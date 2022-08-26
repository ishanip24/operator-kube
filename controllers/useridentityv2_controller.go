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
	"time"

	identityv2 "example.com/m/api/v2"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// v1 "k8s.io/api/core/v1"
	// rbacv1 "k8s.io/api/rbac/v1"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	// "sigs.k8s.io/controller-runtime/pkg/event"
	// "sigs.k8s.io/controller-runtime/pkg/handler"
	// "sigs.k8s.io/controller-runtime/pkg/source"
)

// UserIdentityv2Reconciler reconciles a UserIdentityv2 object
type UserIdentityv2Reconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

const (
	Ready        metav1.ConditionStatus = "Ready"
	UpToDate     metav1.StatusReason    = "UpToDate"
	UpdateFailed metav1.StatusReason    = "UpdateFailed"
)

//+kubebuilder:rbac:groups=identity.company.org,resources=useridentityv2s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=identity.company.org,resources=useridentityv2s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=identity.company.org,resources=useridentityv2s/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the UserIdentityv2 object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *UserIdentityv2Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("useridentity", req.NamespacedName)

	var userIdentity identityv2.UserIdentityv2
	if err := r.Get(ctx, req.NamespacedName, &userIdentity); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	user := "ishani"  // pretend we get the name
	project := "proj" // pretend we get the project name

	log.V(10).Info(fmt.Sprintf("Create Resources for User:%s, Project:%s", user, project))

	var serviceAccount corev1.ServiceAccount
	serviceAccount.Name = "default"
	annotations := make(map[string]string, 1)
	annotations["iam.gke.io/gcp-service-account"] = fmt.Sprintf("%s@%s.iam.gserviceaccount.com", user, project)
	serviceAccount.Annotations = annotations
	_, err := ctrl.CreateOrUpdate(ctx, r.Client, &serviceAccount, func() error {
		return ctrl.SetControllerReference(&userIdentity, &serviceAccount, r.Scheme)
	})
	if err != nil {
		log.Error(err, fmt.Sprintf("Error create ServiceAccount for user: %s, project: %s", user, project))
		_ = r.SetConditionFail(err, userIdentity, log)
		return ctrl.Result{}, nil
	}

	// Create Service Account
	log.V(10).Info(fmt.Sprintf("Create ServiceAccount for User:%s, Project:%s finished", user, project))

	return ctrl.Result{RequeueAfter: 10 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserIdentityv2Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&identityv2.UserIdentityv2{}).
		Complete(r)
}

// create a func to wrap the logic
func (r *UserIdentityv2Reconciler) SetConditionFail(err error, userIdentity identityv2.UserIdentityv2, log logr.Logger) error {
	conditions := userIdentity.GetConditions()
	condition := metav1.Condition{
		Type:    "Ready",
		Status:  metav1.ConditionFalse,
		Reason:  "Update failed",
		Message: err.Error(),
	}
	if conditions.DeepCopyInto(condition) {
		if err := r.Status().Update(context.Background(), &userIdentity); err != nil {
			log.Error(err, "Set conditions failed")
			// r.Recorder.Event(userIdentity, corev1.EventTypeWarning, string(UpdateFailed), "Failed to update resource status")
			return err
		}
	}
	return nil
}
