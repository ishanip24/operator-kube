// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v3

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UserIdentityv3) DeepCopyInto(out *UserIdentityv3) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UserIdentityv3.
func (in *UserIdentityv3) DeepCopy() *UserIdentityv3 {
	if in == nil {
		return nil
	}
	out := new(UserIdentityv3)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *UserIdentityv3) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UserIdentityv3List) DeepCopyInto(out *UserIdentityv3List) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]UserIdentityv3, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UserIdentityv3List.
func (in *UserIdentityv3List) DeepCopy() *UserIdentityv3List {
	if in == nil {
		return nil
	}
	out := new(UserIdentityv3List)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *UserIdentityv3List) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *UserIdentityv3Spec) DeepCopyInto(out *UserIdentityv3Spec) {
	*out = *in
	if in.Template != nil {
		in, out := &in.Template, &out.Template
		*out = make([]TemplateObject, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i].Unstructured)
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UserIdentityv3Spec.
func (in *UserIdentityv3Spec) DeepCopy() *UserIdentityv3Spec {
	if in == nil {
		return nil
	}
	out := new(UserIdentityv3Spec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.

// QUESTION: I feel like double pointer is not ideal...
func (in *UserIdentityv3Status) DeepCopyInto(out *UserIdentityv3Status) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		**out = make([]metav1.Condition, len(**in))
		for i := range **in {
			(**in)[i].DeepCopyInto(&(**out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new UserIdentityv3Status.
func (in *UserIdentityv3Status) DeepCopy() *UserIdentityv3Status {
	if in == nil {
		return nil
	}
	out := new(UserIdentityv3Status)
	in.DeepCopyInto(out)
	return out
}
