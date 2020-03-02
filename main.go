package main

import (
	"os"

	concoursev1 "github.com/cirocosta/crds/api/v1"
	"github.com/cirocosta/crds/controllers"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	// +kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	_ = clientgoscheme.AddToScheme(scheme)
	_ = concoursev1.AddToScheme(scheme)
	// +kubebuilder:scaffold:scheme
}

func main() {
	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: ":9100",
		Port:               9443,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	err = (&controllers.PipelineReconciler{
		Client: mgr.GetClient(),
		Log:    ctrl.Log.WithName("controllers").WithName("Pipeline"),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr)
	if err != nil {
		setupLog.Error(err,
			"unable to create controller",
			"controller", "Pipeline",
		)
		os.Exit(1)
	}

	// 	if err = (&concoursev1.Pipeline{}).SetupWebhookWithManager(mgr); err != nil {
	// 		setupLog.Error(err, "unable to create webhook", "webhook", "Pipeline")
	// 		os.Exit(1)
	// }

	// +kubebuilder:scaffold:builder

	setupLog.Info("starting manager")

	err = mgr.Start(ctrl.SetupSignalHandler())
	if err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}
