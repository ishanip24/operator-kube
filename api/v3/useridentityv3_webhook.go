/*
Copyright 2022 pc.

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

package v3

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var useridentityv3log = logf.Log.WithName("useridentityv3-resource")

func (r *UserIdentityv3) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// TODO(user): EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-identity-company-org-v3-useridentityv3,mutating=true,failurePolicy=fail,sideEffects=None,groups=identity.company.org,resources=useridentityv3s,verbs=create;update,versions=v3,name=museridentityv3.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &UserIdentityv3{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *UserIdentityv3) Default() {
	useridentityv3log.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-identity-company-org-v3-useridentityv3,mutating=false,failurePolicy=fail,sideEffects=None,groups=identity.company.org,resources=useridentityv3s,verbs=create;update,versions=v3,name=vuseridentityv3.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &UserIdentityv3{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *UserIdentityv3) ValidateCreate() error {
	useridentityv3log.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *UserIdentityv3) ValidateUpdate(old runtime.Object) error {
	useridentityv3log.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *UserIdentityv3) ValidateDelete() error {
	useridentityv3log.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
