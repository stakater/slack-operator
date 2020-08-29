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

// Reconcile loop for the Channel resource
func (r *ChannelReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("channel", req.NamespacedName)

	channel := &slackv1alpha1.Channel{}
	err := r.Get(ctx, req.NamespacedName, channel)

	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading channel, requeue
		return ctrl.Result{}, err
	}

	isChannelMarkedToBeDeleted := channel.GetDeletionTimestamp() != nil
	if isChannelMarkedToBeDeleted {
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

	if channel.Status.ID == "" {
		name := channel.Spec.Name
		isPrivate := channel.Spec.Private

		log.Info("Creating new channel", "name", name)

		channelID, err := r.SlackService.CreateChannel(name, isPrivate)

		if err != nil {
			return reconcilerUtil.ManageError(r.Client, channel, err, false)
		}

		channel.Status.ID = *channelID

		err = r.Status().Update(ctx, channel)
		if err != nil {
			log.Error(err, "Failed to update Channel status")
			return ctrl.Result{}, err
		}
		return r.updateSlackChannel(ctx, channel)
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

	return ctrl.Result{}, nil
}

func (r *ChannelReconciler) finalizeChannel(req ctrl.Request, channel *slackv1alpha1.Channel) (ctrl.Result, error) {
	channelID := channel.Status.ID
	log := r.Log.WithValues("channelID", channelID)

	if channel == nil {
		return reconcilerUtil.DoNotRequeue()
	}

	err := r.SlackService.ArchiveChannel(channelID)

	if err != nil && err.Error() != "channel_not_found" {
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