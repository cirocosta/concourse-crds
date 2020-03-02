package v1

import (
	"encoding/json"

	"github.com/concourse/concourse/vars"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var pipelinelog = logf.Log.WithName("pipeline-resource")

func (r *Pipeline) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-concourse-concourse-ci-org-v1-pipeline,mutating=true,failurePolicy=fail,groups=concourse.concourse-ci.org,resources=pipelines,verbs=create;update,versions=v1,name=mpipeline.kb.io

var _ webhook.Defaulter = &Pipeline{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Pipeline) Default() {
	pipelinelog.Info("default", "name", r.Name)

	if r.Spec.Paused == nil {
		r.Spec.Paused = boolRef(true)
	}

	if r.Spec.Vars != nil {
		r.templateConfig()
	}

	// TODO(user): fill in your defaulting logic.
}

func boolRef(v bool) *bool {
	b := new(bool)
	*b = v
	return b
}

func (r *Pipeline) templateConfig() {
	staticVars := new(vars.StaticVariables)

	err := json.Unmarshal(r.Spec.Vars.Raw, staticVars)
	if err != nil {
		pipelinelog.Error(err, "unmarshalling vars")
		return
	}

	template := vars.NewTemplate(r.Spec.Config.Raw)

	res, err := template.Evaluate(staticVars, vars.EvaluateOpts{})
	if err != nil {
		pipelinelog.Error(err, "evaluating template")
		return
	}

	r.Spec.Config = runtime.RawExtension{
		Raw: res,
	}
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-concourse-concourse-ci-org-v1-pipeline,mutating=false,failurePolicy=fail,groups=concourse.concourse-ci.org,resources=pipelines,versions=v1,name=vpipeline.kb.io

var _ webhook.Validator = &Pipeline{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
//
func (r *Pipeline) ValidateCreate() error {
	pipelinelog.Info("validate create", "name", r.Name)

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
//
func (r *Pipeline) ValidateUpdate(old runtime.Object) error {
	pipelinelog.Info("validate update", "name", r.Name)

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
//
func (r *Pipeline) ValidateDelete() error {
	pipelinelog.Info("validate delete", "name", r.Name)

	return nil
}
