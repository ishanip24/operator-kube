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

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.

	"flag"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/pubsub/pstest"
	"github.com/go-logr/logr"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"

	identityv1 "k8s.io/api/core/v1"

	identityv2 "example.com/m/api/v2"
	identityv3 "example.com/m/api/v3"

	"example.com/m/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme                                = runtime.NewScheme()
	setupLog                              = ctrl.Log.WithName("setup")
	pubsubTopic                           = "pull-topic"
	pubsubSubscription                    = "pull-test-results"
	file_pull_test_results_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
)

type Message struct {
	PublishTime time.Time
	Attributes  map[string]string
	ID          string
	OrderingKey string
	modacks     []Modack
	Modacks     []Modack
	Data        []byte
	Acks        int
	deliveries  int
	acks        int
}

// Modack represents a modack sent to the server.
type Modack struct {
	ReceivedAt  time.Time
	AckID       string
	AckDeadline int32
}

type PubSubInfo struct {
	Client     *pubsub.Client
	SecretKey  string
	TestServer *pstest.Server
}

type PullTestResultEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Info          *gcloud.Common       `protobuf:"bytes,1,opt,name=info,proto3" json:"info,omitempty"`
	DistrictKeyId string               `protobuf:"bytes,2,opt,name=district_key_id,json=districtKeyId,proto3" json:"district_key_id,omitempty"`
	Name          string               `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	CreatedAt     *timestamp.Timestamp `protobuf:"bytes,4,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	UpdatedAt     *timestamp.Timestamp `protobuf:"bytes,5,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	Force         bool                 `protobuf:"varint,6,opt,name=force,proto3" json:"force,omitempty"`
	DryRun        bool                 `protobuf:"varint,7,opt,name=dry_run,json=dryRun,proto3" json:"dry_run,omitempty"`
	Test          bool                 `protobuf:"varint,8,opt,name=test,proto3" json:"test,omitempty"`
	GcsFile       string               `protobuf:"bytes,9,opt,name=gcs_file,json=gcsFile,proto3" json:"gcs_file,omitempty"`
}

func (*PullTestResultEvent) ProtoMessage() {}

func (x *PullTestResultEvent) ProtoReflect() protoreflect.Message {
	mi := &file_pull_test_results_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

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
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	flag.Parse()

	ctrl.SetLogger(logr.Logger{})

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		Port:               9443,
		// LeaderElection:     enableLeaderElection,
		// LeaderElectionID:   "95f0db32.company.org",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.UserIdentityReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("UserIdentity"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "UserIdentity")
		os.Exit(1)
	}
	if err = (&controllers.UserIdentityv2Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "UserIdentityV2")
		os.Exit(1)
	}
	if err = (&controllers.UserIdentityv3Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "UserIdentityV3")
		os.Exit(1)
	}
	if err = (&identityv3.UserIdentityv3{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "UserIdentityV3")
		os.Exit(1)
	}
	if err = (&controllers.UserIdentityv3Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "UserIdentityv3")
		os.Exit(1)
	}
	if err = (&identityv3.UserIdentityv3{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "UserIdentityv3")
		os.Exit(1)
	}
	if err = (&controllers.UserIdentityv2Reconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "UserIdentityv2")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// func main() {
// 	ctx := context.Background()
// 	logger, err := zap.NewProduction()
// 	pubSubTopic := "pull-topic"
// 	pubSubID := pubSubTopic + "_job_runner"
// 	subName := "pull-test-results"
// 	client, _ := pubsub.NewClient(ctx, "khan-academy")
// 	clientSet, clientSetErr := NewKubeClient()
// 	if clientSetErr != nil {
// 		log.Fatal(clientSetErr)
// 	}
// 	pubSubInfo := PubSubInfo{
// 		Client:     client,
// 		SecretKey:  "secretkey",
// 		TestServer: pstest.NewServer(),
// 	}

// 	var forever bool
// 	flag.BoolVar(&forever,
// 		"forever",
// 		false,
// 		"Keep pulling from Pull Test Results request pubsub topic forever",
// 	)

// 	if forever {
// 		err = pullForever(ctx, logger, clientSet, pubSubInfo, pubSubID, subName)
// 	} else {
// 		err = pullMsgsAndCreateKubeJobs(ctx, logger, clientSet, pubSubInfo, pubSubID)
// 	}
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// }

// func NewKubeClient() (*coreClient.Clientset, error) {
// 	config, err := restClient.InClusterConfig()
// 	if err != nil {
// 		kubeconfig := cmdClient.NewDefaultClientConfigLoadingRules().GetDefaultFilename()
// 		config, err = cmdClient.BuildConfigFromFlags("", kubeconfig)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
// 	return coreClient.NewForConfig(config)
// }

// func pullMsgsAndCreateKubeJobs(
// 	timeoutctx context.Context,
// 	logger *zap.Logger,
// 	clientSet *coreClient.Clientset,
// 	pubSubInfo PubSubInfo,
// 	subID string,
// ) error {
// 	logger.Info("Beginning pull", zap.String("PubSubID", subID))
// 	sub := pubSubInfo.Client.Subscription(subID)

// 	// Receive messages for 60 seconds.
// 	timeoutctx, cancel := context.WithTimeout(timeoutctx, 60*time.Second)
// 	defer cancel()

// 	// Create a channel to handle messages to as they come in.
// 	cm := make(chan *pubsub.Message)
// 	defer close(cm)
// 	// Handle individual messages in a goroutine.
// 	go func() {
// 		for msg := range cm {
// 			pullTestResultsRequest := &PullTestResultEvent{}
// 			unmarhsallErr := proto.Unmarshal(msg.Data, pullTestResultsRequest)
// 			if unmarhsallErr != nil {
// 				logger.Error(
// 					"unable to unmarshall protobuf",
// 					zap.Error(unmarhsallErr),
// 					zap.String("data", string(msg.Data)),
// 				)
// 			}
// 			logger.Info("Got message", zap.Any("record", pullTestResultsRequest))

// 			expBackoff := backoff.NewExponentialBackOff()

// 			expBackoff.InitialInterval = 30 * time.Second
// 			expBackoff.Multiplier = 2.0
// 			expBackoff.MaxInterval = 5 * time.Minute
// 			expBackoff.MaxElapsedTime = 0
// 			expBackoff.Reset()

// 			// TODO: error checking
// 			msg.Ack()
// 		}
// 	}()
// 	sub.ReceiveSettings.Synchronous = true
// 	// MaxOutstandingMessages is the maximum number of unprocessed messages the
// 	// client will pull from the server before pausing.
// 	//
// 	// This is only guaranteed when ReceiveSettings.Synchronous is set to true.
// 	// When Synchronous is set to false, the StreamingPull RPC is used which
// 	// can pull a single large batch of messages at once that is greater than
// 	// MaxOustandingMessages before pausing. For more info, see
// 	// https://cloud.google.com/pubsub/docs/pull#streamingpull_dealing_with_large_backlogs_of_small_messages.
// 	sub.ReceiveSettings.MaxOutstandingMessages = 10
// 	// Receive blocks until the context is cancelled or a receive error occurs.
// 	err := sub.Receive(timeoutctx, func(ctx context.Context, msg *pubsub.Message) {
// 		cm <- msg
// 	})
// 	if err != nil {
// 		return errors.Wrap(err, "pubsub receive err")
// 	}
// 	logger.Info("Done pulling for now", zap.String("PubSubID", subID))
// 	return nil
// }

// func pullForever(
// 	ctx context.Context,
// 	logger *zap.Logger,
// 	clientSet *coreClient.Clientset,
// 	pubSubInfo PubSubInfo,
// 	pubSubID string,
// 	subID string,
// ) error {
// 	// Go signal notification works by sending `os.Signal`
// 	// values on a channel. We'll create a channel to
// 	// receive these notifications (we'll also make one to
// 	// notify us when the program can exit).
// 	sigs := make(chan os.Signal, 2)
// 	done := make(chan bool, 1)

// 	// Tickers use a similar mechanism to timers: a
// 	// channel that is sent values. Here we'll use the
// 	// `range` builtin on the channel to iterate over
// 	// the values as they arrive every 60s.
// 	ticker := time.NewTicker(time.Second * 60)
// 	var pullErr error

// 	go func() {
// 		for t := range ticker.C {
// 			logger.Info("Tick happened", zap.Time("tick time", t))
// 			err := pullMsgsAndCreateKubeJobs(
// 				ctx,
// 				logger,
// 				clientSet,
// 				pubSubInfo,
// 				pubSubID,
// 			)
// 			if err != nil {
// 				pullErr = err
// 				logger.Error("unable to pull from pubsub topic", zap.Error(err))
// 				ticker.Stop()
// 			}
// 		}
// 	}()

// 	// `signal.Notify` registers the given channel to
// 	// receive notifications of the specified signals.
// 	signal.Notify(sigs,
// 		os.Interrupt,    // interrupt is syscall.SIGINT, Ctrl+C
// 		syscall.SIGQUIT, // Ctrl-\
// 		syscall.SIGHUP,  // "terminal is disconnected"
// 		syscall.SIGTERM, // "the normal way to politely ask a program to terminate"
// 	)
// 	// This goroutine executes a blocking receive for
// 	// signals. When it gets one it'll print it out
// 	// and then notify the program that it can finish.
// 	go func() {
// 		sig := <-sigs
// 		logger.Info("received os Signal", zap.Any("signal", sig))

// 		// Tickers can be stopped like timers.
// 		ticker.Stop()
// 		logger.Info("Ticker stopped")

// 		done <- true
// 	}()
// 	// The program will wait here until it gets the
// 	// expected signal (as indicated by the goroutine
// 	// above sending a value on `done`) and then exit.
// 	logger.Info("awaiting signals")
// 	<-done
// 	logger.Info("exiting")
// 	return pullErr
// }

// Old Main Code
// var metricsAddr string
// var enableLeaderElection bool
// var probeAddr string
// flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
// flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
// flag.BoolVar(&enableLeaderElection, "leader-elect", false,
// 	"Enable leader election for controller manager. "+
// 		"Enabling this will ensure there is only one active controller manager.")
// opts := zap.Options{
// 	Development: true,
// }
// opts.BindFlags(flag.CommandLine)
// flag.Parse()

// ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

// mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
// 	Scheme:                 scheme,
// 	MetricsBindAddress:     metricsAddr,
// 	Port:                   9443,
// 	HealthProbeBindAddress: probeAddr,
// 	LeaderElection:         enableLeaderElection,
// 	LeaderElectionID:       "95f0db32.company.org",
// })
// if err != nil {
// 	setupLog.Error(err, "unable to start manager")
// 	os.Exit(1)
// }

// if err = (&controllers.UserIdentityReconciler{
// 	Client: mgr.GetClient(),
// 	Scheme: mgr.GetScheme(),
// }).SetupWithManager(mgr); err != nil {
// 	setupLog.Error(err, "unable to create controller", "controller", "UserIdentity")
// 	os.Exit(1)
// }

// ctx := context.Background()
// client, err := pubsub.NewClient(ctx, pubsubTopic)
// if err != nil {
// 	log.Fatal(err)
// }
// // Publish "hello world" on pubsubTopic
// topic := client.Topic(pubsubTopic)
// res := topic.Publish(ctx, &pubsub.Message{
// 	Data: []byte("hello world"),
// })
// // The publish happens asynchronously.
// msgID, err := res.Get(ctx)
// if err != nil {
// 	log.Fatal(err)
// }
// setupLog.Info("messageID: " + msgID)

// // Use a callback to receive messages via pubsubSubscription.
// sub := client.Subscription(pubsubSubscription) // create subscription
// err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
// 	fmt.Println(msg.Data)
// 	msg.Ack() // Acknowledge that we've consumed the message.
// })
// if err != nil {
// 	log.Println(err)
// }
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
