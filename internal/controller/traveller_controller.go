package controller

import (
	"context"
	travelv1 "github.com/example-inc/lab8-operator/api/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// +kubebuilder:rbac:groups=traveller.example.com,resources=travellers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=traveller.example.com,resources=travellers/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=traveller.example.com,resources=travellers/finalizers,verbs=update

// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Traveller object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile



type subReconciler interface {
        Reconcile(*travelv1.Traveller) error
}

// TravellerReconciler reconciles a Traveller object
type TravellerReconciler struct {
        client.Client
        Scheme *runtime.Scheme
}


func (r *TravellerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithValues("Traveller", req.NamespacedName)

	var err error
	instance := &travelv1.Traveller{}
	if stop := r.reconcileTraveller(req, instance, err); stop {
		return reconcile.Result{}, err
	}

	deplReconciler := deploymentReconciler{client: r.Client, Scheme: r.Scheme, log: &logger}
	svcReconciler := serviceReconciler{client: r.Client, Scheme: r.Scheme, log: &logger}

	for _, sr := range []subReconciler{deplReconciler, svcReconciler} {
		if err := sr.Reconcile(instance); err != nil {
			return reconcile.Result{}, err
		}
	}

	logger.Info("Skip reconcile: Deployment and service already exists")
	return reconcile.Result{}, nil
}

func (r *TravellerReconciler) reconcileTraveller(req ctrl.Request, instance *travelv1.Traveller, err error) bool {
	err = r.Get(context.TODO(), req.NamespacedName, instance)

	if err != nil {
		if errors.IsNotFound(err) {
			err = nil
		}
		return true
	}

	return false
}

// SetupWithManager sets up the controller with the Manager.
func (r *TravellerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&travelv1.Traveller{}).
		Complete(r)
}
