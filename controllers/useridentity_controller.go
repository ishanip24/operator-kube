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

package controllers

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

	// "k8s.io/kubernetes/pkg/apis/core"
	// "k8s.io/kubernetes/pkg/apis/rbac"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	identityv1 "example.com/m/api/v1"
	// 	identityv1 "k8s.io/apiserver/pkg/storage/value/encrypt/identity"
)

// UserIdentityReconciler reconciles a UserIdentity object
type UserIdentityReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=identity.company.org,resources=useridentities,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=identity.company.org,resources=useridentities/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=identity.company.org,resources=useridentities/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the UserIdentity object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *UserIdentityReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	var userIdentity identityv1.UserIdentity
	if err := r.Get(context.Background(), req.NamespacedName, &userIdentity); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	user := "jenny"      // pretend we get the name
	project := "project" // pretend we get the project name

	var serviceAccount corev1.ServiceAccount
	serviceAccount.Name = "default"
	annotations := make(map[string]string, 1)
	annotations["iam.gke.io/gcp-service-account"] = fmt.Sprintf("%s@%s.iam.gserviceaccount.com", user, project)
	serviceAccount.Annotations = annotations
	_, err := ctrl.CreateOrUpdate(context.Background(), r.Client, &serviceAccount, func() error {
		return ctrl.SetControllerReference(&userIdentity, &serviceAccount, r.Scheme)
	})
	if err != nil {
		return ctrl.Result{}, nil
	}

	var clusterRoleBinding rbacv1.ClusterRoleBinding
	clusterRoleBinding.Name = req.Name
	clusterRoleBinding.Namespace = req.Namespace
	_, err = ctrl.CreateOrUpdate(context.Background(), r.Client, &clusterRoleBinding, func() error {
		clusterRoleBinding.RoleRef = userIdentity.Spec.RoleRef

		clusterRoleBinding.Subjects = []rbacv1.Subject{
			{
				Kind: "ServiceAccount",
				Name: "default",
			},
		}
		return ctrl.SetControllerReference(&userIdentity, &clusterRoleBinding, r.Scheme)
	})

	if err != nil {
		return ctrl.Result{}, nil
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserIdentityReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&identityv1.UserIdentity{}).
		Complete(r)
}
