package v2

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// UserIdentityv2Spec defines the desired state of UserIdentityv2
type UserIdentityv2Spec struct {
	// RoleRef is the target ClusterRole reference
	// +kubebuilder:printcolumn
	RoleRef rbacv1.RoleRef `json:"roleRef,omitempty"`
}

// UserIdentityv2Status defines the observed state of UserIdentityv2
type UserIdentityv2Status struct {
	// Conditions is the list of error conditions for this resource
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// UserIdentityv2 is the Schema for the useridentityv2s API
type UserIdentityv2 struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   UserIdentityv2Spec   `json:"spec,omitempty"`
	Status UserIdentityv2Status `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// UserIdentityv2List contains a list of UserIdentityv2
type UserIdentityv2List struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []UserIdentityv2 `json:"items"`
}

func (o *UserIdentityv2) GetConditions() []metav1.Condition {
	return o.Status.Conditions
}

func init() {
	SchemeBuilder.Register(&UserIdentityv2{}, &UserIdentityv2List{})
}
