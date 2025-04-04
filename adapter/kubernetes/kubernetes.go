package kubernetes

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/bayu-aditya/ideagate/backend/core/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type IKubernetesAdapter interface {
	CreateDeployment(ctx context.Context, deployment appsv1.Deployment) error
	ListPods(ctx context.Context, appName string) ([]corev1.Pod, error)
	SetReplicas(ctx context.Context, deploymentName string, replica int32) error
	RestartDeployment(ctx context.Context, deploymentName string) error
}

func New() (IKubernetesAdapter, error) {
	kubeConfig, err := getKubeConfig()
	if err != nil {
		return nil, err
	}

	// creates the clientSet
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		return nil, err
	}

	return &adapter{
		namespace: "staging",
		client:    clientSet,
	}, nil
}

func getKubeConfig() (*rest.Config, error) {
	// Using in-cluster config
	kubeConfig, err := rest.InClusterConfig()
	if kubeConfig != nil && err == nil {
		return kubeConfig, nil
	}

	// Using out-of-cluster config
	kubeConfigPaths := os.Getenv("KUBECONFIG")
	if kubeConfigPaths == "" {
		kubeConfigPaths = fmt.Sprintf("%s/.kube/config", homedir.HomeDir())
	}

	kubeConfig, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{Precedence: strings.Split(kubeConfigPaths, ":")},
		&clientcmd.ConfigOverrides{},
	).ClientConfig()
	if err != nil {
		return nil, err
	}

	return kubeConfig, nil
}

type adapter struct {
	client    *kubernetes.Clientset
	namespace string
}

func (a *adapter) CreateDeployment(ctx context.Context, deployment appsv1.Deployment) error {
	deploymentClient := a.client.AppsV1().Deployments(a.namespace)

	if _, err := deploymentClient.Create(ctx, &deployment, metav1.CreateOptions{}); err != nil {
		return err
	}

	return nil
}

func (a *adapter) ListPods(ctx context.Context, appName string) ([]corev1.Pod, error) {
	podClient := a.client.CoreV1().Pods(a.namespace)

	options := metav1.ListOptions{
		//LabelSelector: fmt.Sprintf("appName=%s", appName),
	}
	podList, err := podClient.List(ctx, options)
	if err != nil {
		return nil, err
	}

	return podList.Items, nil
}

func (a *adapter) SetReplicas(ctx context.Context, deploymentName string, replica int32) error {
	deploymentClient := a.client.AppsV1().Deployments(a.namespace)

	deployment, err := deploymentClient.Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return fmt.Errorf("deployment %s not found", deploymentName)
		}
		return err
	}

	deployment.Spec.Replicas = utils.ToPtr(replica)

	_, err = deploymentClient.Update(ctx, deployment, metav1.UpdateOptions{})
	return err
}

func (a *adapter) RestartDeployment(ctx context.Context, deploymentName string) error {
	deploymentClient := a.client.AppsV1().Deployments(a.namespace)

	deployment, err := deploymentClient.Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return fmt.Errorf("deployment %s not found", deploymentName)
		}
	}

	deployment.Spec.Template.Annotations = map[string]string{
		"kubectl.kubernetes.io/restartedAt": metav1.Now().Format(time.RFC3339),
	}

	if _, err = deploymentClient.Update(ctx, deployment, metav1.UpdateOptions{}); err != nil {
		return err
	}

	return nil
}
