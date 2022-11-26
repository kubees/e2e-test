package e2e

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestRedisDeployment(t *testing.T) {
	deploymentFeature := features.New("appsv1/deployment").
		Assess("redis database stst creation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var sts appsv1.StatefulSet
			if err := cfg.Client().Resources().Get(ctx, "redis-master", "databases", &sts); err != nil {
				t.Fatal(err)
			}
			if &sts != nil {
				t.Logf("stateful set found: %s", sts.Name)
			}
			availableReplicas := sts.Status.AvailableReplicas
			if availableReplicas != int32(1) {
				t.Errorf("Expected 1 available replica, got %v", availableReplicas)
			}
			readyReplicas := sts.Status.ReadyReplicas
			if readyReplicas != int32(1) {
				t.Errorf("Expected 1 ready replica, got %v", readyReplicas)
			}
			ports := sts.Spec.Template.Spec.Containers[0].Ports
			if len(ports) < 1 {
				t.Errorf("Expected at least 1 port, got %v", len(ports))
			}
			imagePullPolicy := sts.Spec.Template.Spec.Containers[0].ImagePullPolicy
			if imagePullPolicy != "IfNotPresent" {
				t.Errorf("Image pull policy should be set to IfNotPresent, got %v", imagePullPolicy)
			}
			return context.WithValue(ctx, "redis-sts", &sts)
		}).Feature()

	testenv.Test(t, deploymentFeature)
}
func TestRedisSecret(t *testing.T) {
	secretFeature := features.New("v1/secret").
		Assess("redis secret creation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var secret corev1.Secret
			if err := cfg.Client().Resources().Get(ctx, "redis-secret", "databases", &secret); err != nil {
				t.Fatal(err)
			}
			if &secret != nil {
				t.Logf("secret found: %v", &secret.Name)
			}
			password := secret.Data["password"]
			if string(password) != "redis" {
				t.Errorf("Password expected redis, got: %v", password)
			}
			return context.WithValue(ctx, "redis-secret", &secret)
		}).Feature()

	testenv.Test(t, secretFeature)
}

func TestPlaylistdatabaseService(t *testing.T) {
	secretFeature := features.New("v1/service").
		Assess("redis service creation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var service corev1.Service
			if err := cfg.Client().Resources().Get(ctx, "redis-master", "databases", &service); err != nil {
				t.Fatal(err)
			}
			if &service != nil {
				t.Logf("service found: %v", &service.Name)
			}
			selector := service.Spec.Selector
			if selector["app.kubernetes.io/name"] != "redis" {
				t.Errorf("Expected app.kubernetes.io/name=redis, got: app.kubernetes.io/name=%v", selector)
			}
			targetPort := service.Spec.Ports[0].TargetPort.StrVal
			if targetPort != "redis" {
				t.Errorf("Expected target port redis, got %v", targetPort)
			}
			return context.WithValue(ctx, "redis-service", &service)
		}).Feature()

	testenv.Test(t, secretFeature)
}
