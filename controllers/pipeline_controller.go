package controllers

import (
	"context"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/cirocosta/crds/concourse"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	concoursev1 "github.com/cirocosta/crds/api/v1"
)

type Concourse interface {
	SetPipeline(
		ctx context.Context,
		team, name string,
		config []byte,
	) (err error)

	PausePipeline(
		ctx context.Context,
		team, name string,
	) (err error)

	UnpausePipeline(
		ctx context.Context,
		team, name string,
	) (err error)

	DestroyPipeline(
		ctx context.Context,
		team, name string,
	) (err error)
}

const finalizer = "finalizer.concourse-ci.org"

// PipelineReconciler reconciles a Pipeline object
//
type PipelineReconciler struct {
	client.Client
	Log       logr.Logger
	Scheme    *runtime.Scheme
	concourse Concourse
}

// +kubebuilder:rbac:groups=concourse.concourse-ci.org,resources=pipelines,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=concourse.concourse-ci.org,resources=pipelines/status,verbs=get;update;patch

func (r *PipelineReconciler) Reconcile(req ctrl.Request) (result ctrl.Result, err error) {
	var (
		ctx = context.Background()
		log = r.Log.WithValues("pipeline", req.NamespacedName)
	)

	log.Info("start")
	defer func() {
		if err != nil {
			log.Error(err, "failed")
			return
		}

		log.Info("finish")
	}()

	var pipeline concoursev1.Pipeline
	err = r.Get(ctx, req.NamespacedName, &pipeline)
	if err != nil {
		err = client.IgnoreNotFound(err)
		if err != nil {
			err = fmt.Errorf("fetch pipeline: %w", err)
			return
		}
		return
	}

	if !pipeline.ObjectMeta.DeletionTimestamp.IsZero() {
		err = r.destroy(ctx, &pipeline)
		if err != nil {
			err = fmt.Errorf("destroy: %w", err)
			return
		}

		return
	}

	err = r.ensureFinalizerSet(ctx, &pipeline)
	if err != nil {
		err = fmt.Errorf("ensure finalizer set: %w", err)
		return
	}

	err = r.sync(ctx, &pipeline)
	if err != nil {
		err = fmt.Errorf("sync: %w", err)
		return
	}

	return
}

func (r *PipelineReconciler) destroy(ctx context.Context, pipeline *concoursev1.Pipeline) (err error) {
	team, name := pipeline.Spec.Team, pipeline.GetName()

	err = r.concourse.DestroyPipeline(ctx, team, name)
	if err != nil {
		err = fmt.Errorf("destroy pipeline: %w", err)
		return
	}

	err = r.removeFinalizer(ctx, pipeline)
	if err != nil {
		err = fmt.Errorf("remove finalizer: %w", err)
		return
	}

	return
}

func (r *PipelineReconciler) removeFinalizer(ctx context.Context, pipeline *concoursev1.Pipeline) (err error) {
	if !contains(pipeline.ObjectMeta.Finalizers, finalizer) {
		return
	}

	pipeline.ObjectMeta.Finalizers = without(
		pipeline.ObjectMeta.Finalizers,
		finalizer,
	)

	err = r.Update(ctx, pipeline)
	if err != nil {
		err = fmt.Errorf("update: %w", err)
		return
	}

	return
}

func (r *PipelineReconciler) ensureFinalizerSet(ctx context.Context, pipeline *concoursev1.Pipeline) (err error) {
	if contains(pipeline.ObjectMeta.Finalizers, finalizer) {
		return
	}

	pipeline.ObjectMeta.Finalizers = append(
		pipeline.ObjectMeta.Finalizers,
		finalizer,
	)

	err = r.Update(ctx, pipeline)
	if err != nil {
		err = fmt.Errorf("update: %w", err)
		return
	}

	return
}

func (r *PipelineReconciler) sync(ctx context.Context, pipeline *concoursev1.Pipeline) (err error) {
	team, name, config := pipeline.Spec.Team,
		pipeline.GetName(),
		pipeline.Spec.Config.Raw

	err = r.concourse.SetPipeline(ctx, team, name, config)
	if err != nil {
		err = fmt.Errorf("create pipeline: %w", err)
		return
	}

	if pipeline.Spec.Paused != nil {
		if *pipeline.Spec.Paused {
			err = r.concourse.PausePipeline(ctx, team, name)
			if err != nil {
				err = fmt.Errorf("pause pipeline: %w", err)
				return
			}
		} else {
			err = r.concourse.UnpausePipeline(ctx, team, name)
			if err != nil {
				err = fmt.Errorf("unpause pipeline: %w", err)
				return
			}
		}
	}

	pipeline.Status = concoursev1.PipelineStatus{
		LastSetTime: &metav1.Time{Time: time.Now()},
	}

	err = r.Status().Update(ctx, pipeline)
	if err != nil {
		err = fmt.Errorf("updating status: %w", err)
		return
	}

	return
}

func (r *PipelineReconciler) SetupWithManager(mgr ctrl.Manager) (err error) {
	var cfg struct {
		Url      string `env:"CONCOURSE_URL"      envDefault:"http://localhost:8080"`
		User     string `env:"CONCOURSE_USERNAME" envDefault:"test"`
		Password string `env:"CONCOURSE_PASSWORD" envDefault:"test"`
	}

	err = env.Parse(&cfg)
	if err != nil {
		err = fmt.Errorf("env parse: %w", err)
		return
	}

	client, err := concourse.New(
		context.Background(),
		cfg.Url, cfg.User, cfg.Password,
	)
	if err != nil {
		err = fmt.Errorf("concourse new: %w", err)
		return
	}

	r.concourse = client

	err = ctrl.NewControllerManagedBy(mgr).
		For(&concoursev1.Pipeline{}).
		Complete(r)

	return
}

func without(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}

		result = append(result, item)
	}

	return
}

func contains(strs []string, str string) bool {
	for _, s := range strs {
		if s != str {
			continue
		}

		return true
	}

	return false
}
