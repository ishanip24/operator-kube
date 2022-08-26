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

	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	identityv3 "example.com/m/api/v3"
	"example.com/m/health"
)

// UserIdentityv3Reconciler reconciles a UserIdentityv3 object
type UserIdentityv3Reconciler struct {
	client.Client
	Scheme *runtime.Scheme
	health.HealthCheck
	Log logr.Logger
}

type Param struct {
	User               string
	Project            string
	ServiceAccountName string
}

//+kubebuilder:rbac:groups=identity.company.org,resources=useridentityv3s,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=identity.company.org,resources=useridentityv3s/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=identity.company.org,resources=useridentityv3s/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the UserIdentityv3 object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *UserIdentityv3Reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()
	log := r.Log.WithValues("useridentity", req.NamespacedName)

	// your logic here

	// update execution time
	r.HealthCheck.Trigger()

	var userIdentity identityv3.UserIdentityv3
	if err := r.Get(ctx, req.NamespacedName, &userIdentity); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	user := "ishani"
	project := "project"

	t := template.New("")

	renderTemplate := func(idx int, user string, project string) (*unstructured.Unstructured, error) {
		var buf bytes.Buffer
		if err := t.ExecuteTemplate(&buf, strconv.Itoa(idx), Param{
			User:               user,
			Project:            project,
			ServiceAccountName: "default",
		}); err != nil {
			return nil, err
		}
		decoder := yaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
		u := unstructured.Unstructured{}
		_, _, err := decoder.Decode(buf.Bytes(), nil, &u)
		return &u, err
	}

	// Apply resources
	for i := range userIdentity.Spec.Foo {
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

			_, err = ctrl.CreateOrUpdate(ctx, r.Client, &existing, func() error {
				if err := mergo.Merge(&existing, rendered, mergo.WithOverride); err != nil {
					return err
				}
				labels := existing.GetLabels()
				if labels == nil {
					labels = make(map[string]string)
				}

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
	return ctrl.NewControllerManagedBy(mgr).
		For(&identityv3.UserIdentityv3{}).
		Complete(r)
}
