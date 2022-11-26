package e2e

import (
	"context"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"
)

func TestVideosMicroserviceDeployment(t *testing.T) {
	deploymentFeature := features.New("appsv1/deployment").
		Assess("videos microservice creation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var dep appsv1.Deployment
			if err := cfg.Client().Resources().Get(ctx, "videos-api", "microservices", &dep); err != nil {
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
			if len(ports) != 2 {
				t.Errorf("Expected 2 ports, got %v", len(ports))
			}
			imagePullPolicy := dep.Spec.Template.Spec.Containers[0].ImagePullPolicy
			if imagePullPolicy != "Always" {
				t.Errorf("Image pull policy should be set to always")
			}
			return context.WithValue(ctx, "videos-api-deployment", &dep)
		}).Feature()

	testenv.Test(t, deploymentFeature)
}
func TestVideosMicroserviceSecret(t *testing.T) {
	secretFeature := features.New("v1/secret").
		Assess("videos secret creation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var secret corev1.Secret
			if err := cfg.Client().Resources().Get(ctx, "videos-redis-secret", "microservices", &secret); err != nil {
				t.Fatal(err)
			}
			if &secret != nil {
				t.Logf("secret found: %v", &secret.Name)
			}
			password := secret.Data["PASSWORD"]
			if string(password) != "redis" {
				t.Errorf("Password expected redis, got: %v", password)
			}
			return context.WithValue(ctx, "videos-api-secret", &secret)
		}).Feature()

	testenv.Test(t, secretFeature)
}

func TestVideosMicroserviceConfigmap(t *testing.T) {
	secretFeature := features.New("v1/configmap").
		Assess("videos configmap creation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var configmap corev1.ConfigMap
			if err := cfg.Client().Resources().Get(ctx, "videos-env", "microservices", &configmap); err != nil {
				t.Fatal(err)
			}
			if &configmap != nil {
				t.Logf("configmap found: %v", &configmap.Name)
			}
			redisHost := configmap.Data["REDIS_HOST"]
			if string(redisHost) != "redis-master.databases" {
				t.Errorf("Password expected redis-master.databases, got: %v", redisHost)
			}
			redisPort := configmap.Data["REDIS_PORT"]
			if string(redisPort) != "6379" {
				t.Errorf("Password expected 6379, got: %v", redisPort)
			}
			return context.WithValue(ctx, "videos-api-configmap", &configmap)
		}).Feature()

	testenv.Test(t, secretFeature)
}

func TestVideosMicroserviceService(t *testing.T) {
	secretFeature := features.New("v1/service").
		Assess("videos service creation", func(ctx context.Context, t *testing.T, cfg *envconf.Config) context.Context {
			var service corev1.Service
			if err := cfg.Client().Resources().Get(ctx, "videos-api", "microservices", &service); err != nil {
				t.Fatal(err)
			}
			if &service != nil {
				t.Logf("service found: %v", &service.Name)
			}
			selector := service.Spec.Selector
			if selector["app"] != "videos-api" {
				t.Errorf("Expected app=videos-api, got: app=%v", selector)
			}
			targetPort := service.Spec.Ports[0].TargetPort.IntVal
			if targetPort != 10010 {
				t.Errorf("Expected target port 10010, got %v", targetPort)
			}
			return context.WithValue(ctx, "videos-api-service", &service)
		}).Feature()

	testenv.Test(t, secretFeature)
}
