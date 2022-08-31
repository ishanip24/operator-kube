/*
Copyright 2020 pc.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
type TemplateObject struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	unstructured.Unstructured `json:",inline"`
}

// UserIdentityv3Spec defines the desired state of UserIdentityv3
type UserIdentityv3Spec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Template is a list of resources to instantiate per repository in Governator
	Template []TemplateObject `json:"template,omitempty"`
}

// UserIdentityv3Status defines the observed state of UserIdentityv3
type UserIdentityv3Status struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions is the list of error conditions for this resource
	Conditions *[]metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true

// UserIdentityv3 is the Schema for the UserIdentityv3s API
type UserIdentityv3 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserIdentityv3Spec   `json:"spec,omitempty"`
	Status UserIdentityv3Status `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// UserIdentityv3List contains a list of UserIdentityv3
type UserIdentityv3List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UserIdentityv3 `json:"items"`
}

func init() {
	SchemeBuilder.Register(&UserIdentityv3{}, &UserIdentityv3List{})
}

func (o *UserIdentityv3) GetConditions() *[]metav1.Condition {
	return o.Status.Conditions
}
