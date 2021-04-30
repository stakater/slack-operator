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

package v1alpha1

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var channellog = logf.Log.WithName("channel-resource")

func (r *Channel) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

// +kubebuilder:webhook:path=/mutate-slack-stakater-com-v1alpha1-channel,mutating=true,failurePolicy=fail,sideEffects=None,groups=slack.stakater.com,resources=channels,verbs=create;update,versions=v1alpha1,name=mchannel.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Defaulter = &Channel{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Channel) Default() {
	channellog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
// +kubebuilder:webhook:verbs=create;update,path=/validate-slack-stakater-com-v1alpha1-channel,mutating=false,failurePolicy=fail,sideEffects=None,groups=slack.stakater.com,resources=channels,versions=v1alpha1,name=vchannel.kb.io,admissionReviewVersions={v1,v1beta1}

var _ webhook.Validator = &Channel{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Channel) ValidateCreate() error {
	channellog.Info("validate create", "name", r.Name)

	if len(r.Spec.Users) < 1 {
		return fmt.Errorf("Users can not be empty")
	}
	if r.Spec.Topic == "test" {
		return fmt.Errorf("Users can not be empty")
	}

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Channel) ValidateUpdate(old runtime.Object) error {
	channellog.Info("validate update", "name", r.Name)

	oldChannel, ok := old.(*Channel)
	if !ok {
		return fmt.Errorf("Error casting old runtime object to %T from %T", oldChannel, old)
	}

	if len(r.Spec.Users) < 1 {
		return fmt.Errorf("Users can not be empty")
	}

	return ValidateImmutableFields(r, oldChannel)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Channel) ValidateDelete() error {
	channellog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func ValidateImmutableFields(newChannel *Channel, oldChannel *Channel) error {
	if oldChannel.Spec.Private != newChannel.Spec.Private {
		return fmt.Errorf("Field 'isPrivate' is immutable and cannot be changed after Slack Channel has been created")
	}
	return nil
}
