/*
Copyright 2024 jongsoo.

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

package controller

import (
	"context"
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strconv"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	jsgroupv1beta1 "js.domain/fastapi-operator/api/v1beta1"
)

// JsKindReconciler reconciles a JsKind object
type JsKindReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=js-group.js.domain,resources=jskinds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=js-group.js.domain,resources=jskinds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=js-group.js.domain,resources=jskinds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the JsKind object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *JsKindReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	// TODO(user): your logic here

	var app jsgroupv1beta1.JsKind

	if err := r.Get(ctx, req.NamespacedName, &app); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		l.Error(err, "JsKind을 가져올 수 없습니다!")
		return ctrl.Result{}, err
	}

	/*
	   finalizer는 쿠버네티스에게 오브젝트 삭제에 대한 통제가 필요하다는 것을 알려주기 때문에 필수!
	   finalizer가 없으면, 쿠버네티스 가비지 컬렉터를 삭제할 수 없고 클러스터에 쓸모없는 리소스를 가질 위험이 있다.
	*/
	if !controllerutil.ContainsFinalizer(&app, finalizer) {
		l.Info("Finalizer 추가")
		controllerutil.AddFinalizer(&app, finalizer)
		return ctrl.Result{}, r.Update(ctx, &app)
	}

	if !app.DeletionTimestamp.IsZero() {
		l.Info("JsKind 삭제")
		return r.reconcileDelete(ctx, &app)
	}

	l.Info("JsKind 생성")

	l.Info("Deployment 생성")
	err := r.createOrUpdateDeployment(ctx, &app)
	if err != nil {
		return ctrl.Result{}, err
	}
	l.Info("Service 생성")
	err = r.createOrUpdateService(ctx, &app)
	if err != nil {
		return ctrl.Result{}, err
	}
	l.Info("Ingress 생성")
	err = r.createOrUpdateIngress(ctx, &app)
	if err != nil {
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JsKindReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&jsgroupv1beta1.JsKind{}).
		Complete(r)
}

func (r *JsKindReconciler) createOrUpdateDeployment(ctx context.Context, app *jsgroupv1beta1.JsKind) error {
	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.ObjectMeta.Name,
			Namespace: app.ObjectMeta.Namespace,
			Labels:    map[string]string{"label": app.ObjectMeta.Name, "app": app.ObjectMeta.Name},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &app.Spec.Size,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"label": app.ObjectMeta.Name},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"label": app.ObjectMeta.Name, "app": app.ObjectMeta.Name},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  app.ObjectMeta.Name + "-container",
							Image: app.Spec.Image,
							Ports: []v1.ContainerPort{
								{
									ContainerPort: app.Spec.Port,
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, &deployment, func() error {
		return nil
	})
	if err != nil {
		return fmt.Errorf("Deployment를 가져올 수 없습니다: %v", err)
	}
	return nil
}

func (r *JsKindReconciler) createOrUpdateService(ctx context.Context, app *jsgroupv1beta1.JsKind) error {
	service := v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.ObjectMeta.Name,
			Namespace: app.ObjectMeta.Namespace,
			Labels:    map[string]string{"app": app.ObjectMeta.Name},
		},
		Spec: v1.ServiceSpec{
			Type:                  v1.ServiceTypeNodePort,
			ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeLocal,
			Selector:              map[string]string{"app": app.ObjectMeta.Name},
			Ports: []v1.ServicePort{
				{
					Name:       "http",
					Port:       app.Spec.Port,
					Protocol:   v1.ProtocolTCP,
					TargetPort: intstr.Parse(strconv.Itoa(int(app.Spec.Port))),
				},
			},
		},
		Status: v1.ServiceStatus{},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, &service, func() error {
		return nil
	})

	if err != nil {
		return fmt.Errorf("Service를 생성할 수 없습니다!: %v", err)
	}
	return nil
}
func (r *JsKindReconciler) createOrUpdateIngress(ctx context.Context, app *jsgroupv1beta1.JsKind) error {
	var ingressClassName = "nginx"
	var pathType = networkingv1.PathTypePrefix

	ingress := networkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      app.ObjectMeta.Name,
			Namespace: app.ObjectMeta.Namespace,
			Labels:    map[string]string{"app": app.ObjectMeta.Name},
		},
		Spec: networkingv1.IngressSpec{
			IngressClassName: &ingressClassName,
			Rules: []networkingv1.IngressRule{
				{
					IngressRuleValue: networkingv1.IngressRuleValue{
						HTTP: &networkingv1.HTTPIngressRuleValue{
							Paths: []networkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: &pathType,
									Backend: networkingv1.IngressBackend{
										Service: &networkingv1.IngressServiceBackend{
											Name: app.ObjectMeta.Name,
											Port: networkingv1.ServiceBackendPort{
												Number: app.Spec.Port,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	_, err := controllerutil.CreateOrUpdate(ctx, r.Client, &ingress, func() error {
		return nil
	})

	if err != nil {
		return fmt.Errorf("Ingress를 생성할 수 없습니다!: %v", err)
	}
	return nil
}

const finalizer = "my-finalizer"

func (r *JsKindReconciler) reconcileDelete(ctx context.Context, app *jsgroupv1beta1.JsKind) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	l.Info("removing application")

	fmt.Printf("app.ObjectMeta.Finalizers: %+v", app.ObjectMeta.Finalizers)

	controllerutil.RemoveFinalizer(app, finalizer)
	err := r.Update(ctx, app)
	if err != nil {
		return ctrl.Result{}, fmt.Errorf("Error removing finalizer %v", err)
	}
	return ctrl.Result{}, nil
}
