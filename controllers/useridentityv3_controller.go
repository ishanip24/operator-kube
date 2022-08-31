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
	"bytes"
	"context"
	"fmt"
	"strconv"
	"text/template"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/source"

	identityv2 "example.com/m/api/v2"
	identityv3 "example.com/m/api/v3"
	"example.com/m/health"
)

// UserIdentityv3Reconciler reconciles a UserIdentityv3 object
type UserIdentityv3Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
	health.HealthCheck
	Log          logr.Logger
	PubsubClient *pubsub.Client
	Recorder     record.EventRecorder
}

type Param struct {
	User               string
	Project            string
	ServiceAccountName string
	NameSpace          string
}

//+kubebuilder:rbac:groups=identity.company.org,resources=useridentityv3s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=identity.company.org,resources=useridentityv3s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=identity.company.org,resources=useridentityv3s/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *UserIdentityv3Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	log := r.Log.WithValues("useridentity", req.NamespacedName)

	// update health checks
	r.HealthCheck.Trigger()

	// get useridentity info from the server
	var userIdentity identityv3.UserIdentityv3
	if err := r.Get(ctx, req.NamespacedName, &userIdentity); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	user := "ishani"
	project := "project"
	serviceAccountName := "default"
	namespace := "roster-sync"

	// template
	t := template.New(`
	apiVersion: v1
	kind: ServiceAccount
	metadata:
	  annotations:
		iam.gke.io/gcp-service-account: {{.User}}@{{.Project}}.iam.gserviceaccount.com
	  name: {{.ServiceAccountName}}
	  namespace: {{.NameSpace}}
`)
	renderTemplate := func(idx int, user string, project string) (*unstructured.Unstructured, error) {
		var buf bytes.Buffer
		// exec json template
		if err := t.ExecuteTemplate(&buf, strconv.Itoa(idx), Param{
			User:               user,
			Project:            project,
			ServiceAccountName: serviceAccountName,
			NameSpace:          namespace,
		}); err != nil {
			return nil, err
		}
		// decode the json file itself and store in u
		decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
		u := unstructured.Unstructured{}
		_, _, err := decoder.Decode(buf.Bytes(), nil, &u)
		return &u, err
	}

	// get each desired TemplateObject
	// QUESTION: Does this assume that the template name is the index of the array in userIdentity.Spec.Template
	// Not sure if I need to keep this indexing logic if not.
	for i := range userIdentity.Spec.Template {
		if err := func() error {
			rendered, err := renderTemplate(i, user, project)
			if err != nil {
				return err
			}

			var existing unstructured.Unstructured
			existing.SetGroupVersionKind(rendered.GroupVersionKind())
			existing.SetNamespace(rendered.GetNamespace())
			existing.SetName(rendered.GetName())

			log.V(10).Info(fmt.Sprintf("expanding %v into %v/%v", rendered.GroupVersionKind(), rendered.GetNamespace(), rendered.GetName()))

			// update object or create obj if DNE based on desired
			_, err = ctrl.CreateOrUpdate(ctx, r.Client, &existing, func() error {
				if err := mergo.Merge(&existing, rendered, mergo.WithOverride); err != nil {
					return err
				}
				labels := existing.GetLabels()
				if labels == nil {
					labels = make(map[string]string)
				}
				// garbage collection
				return ctrl.SetControllerReference(&userIdentity, &existing, r.Scheme)
			})

			return err
		}(); err != nil {
			log.Error(err, fmt.Sprintf("Create resources for user:%s err", user))
			return ctrl.Result{}, nil
		}
	}

	return ctrl.Result{RequeueAfter: 10 * time.Minute}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserIdentityv3Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	// define userevent and run
	ch := make(chan event.GenericEvent)
	subscription := r.PubsubClient.Subscription("pull-test-results")
	userEvent := CreateUserEvents(mgr.GetClient(), subscription, ch)
	go userEvent.Run()

	// return controller
	// perdicate: used by Controllers to filter Events before they are provided to EventHandlers
	return ctrl.NewControllerManagedBy(mgr).
		For(&identityv2.UserIdentityv2{}).
		WithEventFilter(predicate.Or(
			predicate.GenerationChangedPredicate{},
			predicate.NewPredicateFuncs(func(object client.Object) bool {
				obj, ok := object.(Object)
				return ok && IsStatusConditionFalse(*obj.GetConditions(), metav1.ConditionFalse)
			}),
		)).
		Watches(&source.Channel{Source: ch, DestBufferSize: 1024}, &handler.EnqueueRequestForObject{}).
		Complete(r)
}

func (r *UserIdentityv3Reconciler) SetConditionFail(ctx context.Context, err error, userIdentity identityv3.UserIdentityv3, log logr.Logger) error {
	conditions := userIdentity.GetConditions()
	r.Recorder.Event(&userIdentity, corev1.EventTypeWarning, string("Failed"), err.Error())
	if meta.IsStatusConditionPresentAndEqual(*conditions, "Ready", metav1.ConditionFalse) {
		if err := r.Status().Update(ctx, &userIdentity); err != nil {
			log.Error(err, "Set conditions failed")
			r.Recorder.Event(&userIdentity, corev1.EventTypeWarning, string(UpdateFailed), "Failed to update resource status")
			return err
		}
	}
	return nil
}

type Object interface {
	runtime.Object
	metav1.Object
	GetConditions() *[]metav1.ConditionStatus
}

func IsStatusConditionFalse(statuses []metav1.ConditionStatus, condition metav1.ConditionStatus) bool {
	for _, status := range statuses {
		if status == condition {
			return true
		}
	}
	return false
}

// containsString returns true if the string is contained in the slice
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
