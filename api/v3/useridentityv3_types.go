package v3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

type TemplateObject struct {
	// +kubebuilder:pruning:PreserveUnknownFields
	unstructured.Unstructured `json:",inline"`
}

// UserIdentityv3Spec defines the desired state of UserIdentityv3
type UserIdentityv3Spec struct {
	// Template is a list of resources to instantiate per repository in Governator
	Template []TemplateObject `json:"template,omitempty"`
}

// UserIdentityv3Status defines the observed state of UserIdentityv3
type UserIdentityv3Status struct {
	// Conditions is the list of error conditions for this resource
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// UserIdentityv3 is the Schema for the useridentityv3s API
type UserIdentityv3 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserIdentityv3Spec   `json:"spec,omitempty"`
	Status UserIdentityv3Status `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// UserIdentityv3List contains a list of UserIdentityv3
type UserIdentityv3List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UserIdentityv3 `json:"items"`
}

func (o *UserIdentityv3) GetConditions() []metav1.Condition {
	return o.Status.Conditions
}

func init() {
	SchemeBuilder.Register(&UserIdentityv3{}, &UserIdentityv3List{})
}
