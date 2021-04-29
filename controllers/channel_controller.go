/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	finalizerUtil "github.com/stakater/operator-utils/util/finalizer"
	reconcilerUtil "github.com/stakater/operator-utils/util/reconciler"
	slackv1alpha1 "github.com/stakater/slack-operator/api/v1alpha1"
	slack "github.com/stakater/slack-operator/pkg/slack"
)

var (
	channelFinalizer string = "slack.stakater.com/channel"
)

// ChannelReconciler reconciles a Channel object
type ChannelReconciler struct {
	client.Client
	Log          logr.Logger
	Scheme       *runtime.Scheme
	SlackService slack.Service
}

// +kubebuilder:rbac:groups=slack.stakater.com,resources=channels,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=slack.stakater.com,resources=channels/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=secrets,verbs=get;list

// Reconcile loop for the Channel resource
func (r *ChannelReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("channel", req.NamespacedName)

	channel := &slackv1alpha1.Channel{}
	err := r.Get(ctx, req.NamespacedName, channel)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcilerUtil.DoNotRequeue()
		}
		// Error reading channel, requeue
		return reconcilerUtil.RequeueWithError(err)
	}

	// Channel is marked for deletion
	if channel.GetDeletionTimestamp() != nil {
		log.Info("Deletion timestamp found for channel " + req.Name)
		if finalizerUtil.HasFinalizer(channel, channelFinalizer) {
			return r.finalizeChannel(req, channel)
		}
		// Finalizer doesn't exist so clean up is already done
		return reconcilerUtil.DoNotRequeue()
	}

	// Add finalizer if it doesn't exist
	if !finalizerUtil.HasFinalizer(channel, channelFinalizer) {
		log.Info("Adding finalizer for channel " + req.Name)

		finalizerUtil.AddFinalizer(channel, channelFinalizer)

		err := r.Client.Update(ctx, channel)
		if err != nil {
			return reconcilerUtil.ManageError(r.Client, channel, err, true)
		}
	}

	// Check for validity of slack channel custom resource
	err = r.SlackService.IsValidChannel(channel)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, channel, err, true)
	}

	if channel.Status.ID == "" {
		name := channel.Spec.Name
		isPrivate := channel.Spec.Private

		log.Info("Creating new channel", "name", name)

		channelID, err := r.SlackService.CreateChannel(name, isPrivate)
		if err != nil {
			if err.Error() == "name_taken" {
				// Check if the channel already exists and then just reconstruct the status accordingly
				existingChannel, err := r.SlackService.GetChannelByName(name)
				if err != nil {
					return reconcilerUtil.ManageError(r.Client, channel, err, false)
				}

				if existingChannel != nil && existingChannel.GroupConversation.IsArchived {
					err = r.SlackService.UnArchiveChannel(existingChannel)
					if err != nil {
						return reconcilerUtil.ManageError(r.Client, channel, err, false)
					}
				}
				channelID = &existingChannel.ID
			} else {
				return reconcilerUtil.ManageError(r.Client, channel, err, false)
			}
		}

		channel.Status.ID = *channelID

		err = r.Status().Update(ctx, channel)
		if err != nil {
			log.Error(err, "Failed to update Channel status")
			return reconcilerUtil.ManageError(r.Client, channel, err, true)
		}
		return r.updateSlackChannel(ctx, channel)
	}

	existingChannel, err := r.SlackService.GetChannel(channel.Status.ID)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, channel, err, true)
	}

	existingChannelCR := r.SlackService.GetChannelCRFromChannel(existingChannel)

	err = slackv1alpha1.ValidateImmutableFields(existingChannelCR, channel)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, channel, err, true)
	}

	updated, err := r.SlackService.IsChannelUpdated(channel)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, channel, err, true)
	}

	if !updated {
		log.Info("Skipping update. No changes found")
		return reconcilerUtil.DoNotRequeue()
	}

	return r.updateSlackChannel(ctx, channel)
}

func (r *ChannelReconciler) updateSlackChannel(ctx context.Context, channel *slackv1alpha1.Channel) (ctrl.Result, error) {
	channelID := channel.Status.ID
	log := r.Log.WithValues("channelID", channelID)

	log.Info("Updating channel details")

	name := channel.Spec.Name
	users := channel.Spec.Users
	topic := channel.Spec.Topic
	description := channel.Spec.Description

	_, err := r.SlackService.RenameChannel(channelID, name)
	if err != nil {
		log.Error(err, "Error renaming channel")
		return reconcilerUtil.ManageError(r.Client, channel, err, false)
	}

	_, err = r.SlackService.SetTopic(channelID, topic)
	if err != nil {
		log.Error(err, "Error setting channel topic")
		return reconcilerUtil.ManageError(r.Client, channel, err, false)
	}

	_, err = r.SlackService.SetDescription(channelID, description)
	if err != nil {
		log.Error(err, "Error setting channel description")
		return reconcilerUtil.ManageError(r.Client, channel, err, false)
	}

	err = r.SlackService.InviteUsers(channelID, users)
	if err != nil {
		log.Error(err, "Error inviting users to channel")
		return reconcilerUtil.ManageError(r.Client, channel, err, false)
	}

	err = r.SlackService.RemoveUsers(channelID, users)
	if err != nil {
		log.Error(err, "Error removing users from the channel")
		return reconcilerUtil.ManageError(r.Client, channel, err, false)
	}

	return reconcilerUtil.ManageSuccess(r.Client, channel)
}

func (r *ChannelReconciler) finalizeChannel(req ctrl.Request, channel *slackv1alpha1.Channel) (ctrl.Result, error) {
	if channel == nil {
		return reconcilerUtil.DoNotRequeue()
	}

	channelID := channel.Status.ID
	log := r.Log.WithValues("channelID", channelID)

	err := r.SlackService.ArchiveChannel(channelID)

	if err != nil && err.Error() != "channel_not_found" && err.Error() != "already_archived" {
		return reconcilerUtil.ManageError(r.Client, channel, err, false)
	}

	finalizerUtil.DeleteFinalizer(channel, channelFinalizer)
	log.V(1).Info("Finalizer removed for channel")

	err = r.Client.Update(context.Background(), channel)
	if err != nil {
		return reconcilerUtil.ManageError(r.Client, channel, err, false)
	}

	return reconcilerUtil.DoNotRequeue()
}

// SetupWithManager - Controller-Manager binding configuration
func (r *ChannelReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&slackv1alpha1.Channel{}).
		Complete(r)
}
