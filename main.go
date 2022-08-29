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

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	"cloud.google.com/go/pubsub"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	identityv1 "k8s.io/api/core/v1"

	identityv2 "example.com/m/api/v2"
	identityv3 "example.com/m/api/v3"
	"example.com/m/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme             = runtime.NewScheme()
	setupLog           = ctrl.Log.WithName("setup")
	pubsubTopic        = "district_create_report_ephemeral_record_event"
	pubsubSubscription = "pull-test-results"
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(identityv1.AddToScheme(scheme))
	utilruntime.Must(identityv2.AddToScheme(scheme))
	utilruntime.Must(identityv3.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "95f0db32.company.org",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.UserIdentityReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "UserIdentity")
		os.Exit(1)
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, pubsubTopic)
	if err != nil {
		log.Fatal(err)
	}
	// Publish "hello world" on pubsubTopic
	topic := client.Topic(pubsubTopic)
	res := topic.Publish(ctx, &pubsub.Message{
		Data: []byte("hello world"),
	})
	// The publish happens asynchronously.
	msgID, err := res.Get(ctx)
	if err != nil {
		log.Fatal(err)
	}
	setupLog.Info("messageID: " + msgID)

	// Use a callback to receive messages via pubsubSubscription.
	sub := client.Subscription(pubsubSubscription) // create subscription
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		fmt.Println(msg.Data)
		msg.Ack() // Acknowledge that we've consumed the message.
	})
	if err != nil {
		log.Println(err)
	}
	// if err = (&controllers.UserIdentityv2Reconciler{
	// 	Client: mgr.GetClient(),
	// 	Scheme: mgr.GetScheme(),
	// }).SetupWithManager(mgr); err != nil {
	// 	setupLog.Error(err, "unable to create controller", "controller", "UserIdentityv2")
	// 	os.Exit(1)
	// }
	// if err = (&controllers.UserIdentityv3Reconciler{
	// 	Client: mgr.GetClient(),
	// 	Scheme: mgr.GetScheme(),
	// }).SetupWithManager(mgr); err != nil {
	// 	setupLog.Error(err, "unable to create controller", "controller", "UserIdentityv3")
	// 	os.Exit(1)
	// }
	// if err = (&identityv3.UserIdentityv3{}).SetupWebhookWithManager(mgr); err != nil {
	// 	setupLog.Error(err, "unable to create webhook", "webhook", "UserIdentityv3")
	// 	os.Exit(1)
	// }
	// //+kubebuilder:scaffold:builder

	// if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
	// 	setupLog.Error(err, "unable to set up health check")
	// 	os.Exit(1)
	// }
	// if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
	// 	setupLog.Error(err, "unable to set up ready check")
	// 	os.Exit(1)
	// }

	// setupLog.Info("starting manager")
	// if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
	// 	setupLog.Error(err, "problem running manager")
	// 	os.Exit(1)
	// }
}
