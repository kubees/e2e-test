package e2e

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestFrontendDeployment(t *testing.T) {
	deploymentFeature := features.New("appsv1/deployment").
		Assess("frontend deployment creation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var dep appsv1.Deployment
			if err := cfg.Client().Resources().Get(ctx, "videos-web", "frontend", &dep); err != nil {
				t.Fatal(err)
			}
			if &dep != nil {
				t.Logf("deployment found: %s", dep.Name)
			}
			availableReplicas := dep.Status.AvailableReplicas
			if availableReplicas != int32(1) {
				t.Errorf("Expected 1 available replica, got %v", availableReplicas)
			}
			readyReplicas := dep.Status.ReadyReplicas
			if readyReplicas != int32(1) {
				t.Errorf("Expected 1 ready replica, got %v", readyReplicas)
			}
			ports := dep.Spec.Template.Spec.Containers[0].Ports
			if len(ports) != 1 {
				t.Errorf("Expected 1 port, got %v", len(ports))
			}
			imagePullPolicy := dep.Spec.Template.Spec.Containers[0].ImagePullPolicy
			if imagePullPolicy != "Always" {
				t.Errorf("Image pull policy should be set to always")
			}
			return context.WithValue(ctx, "videos-web-deployment", &dep)
		}).Feature()

	testenv.Test(t, deploymentFeature)
}

func TestFrontendService(t *testing.T) {
	secretFeature := features.New("v1/service").
		Assess("frontend service creation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var service corev1.Service
			if err := cfg.Client().Resources().Get(ctx, "videos-web", "frontend", &service); err != nil {
				t.Fatal(err)
			}
			if &service != nil {
				t.Logf("service found: %v", &service.Name)
			}
			selector := service.Spec.Selector
			if selector["app"] != "videos-web" {
				t.Errorf("Expected app=videos-web, got: app=%v", selector)
			}
			targetPort := service.Spec.Ports[0].TargetPort.IntVal
			if targetPort != 80 {
				t.Errorf("Expected target port 80, got %v", targetPort)
			}
			return context.WithValue(ctx, "videos-web-service", &service)
		}).Feature()

	testenv.Test(t, secretFeature)
}
