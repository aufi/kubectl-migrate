package framework

import (
	"context"
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ClusterManager manages Kubernetes cluster interactions for testing
type ClusterManager struct {
	Client    kubernetes.Interface
	Dynamic   dynamic.Interface
	Context   string
	Namespace string
}

// NewClusterManager creates a new cluster manager using the current kubeconfig context
func NewClusterManager() (*ClusterManager, error) {
	return NewClusterManagerWithContext("")
}

// NewClusterManagerWithContext creates a new cluster manager with a specific context
func NewClusterManagerWithContext(contextName string) (*ClusterManager, error) {
	// Load kubeconfig
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}

	if contextName != "" {
		configOverrides.CurrentContext = contextName
	}

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}

	// Create dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create dynamic client: %w", err)
	}

	// Get current namespace
	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		namespace = "default"
	}

	return &ClusterManager{
		Client:    clientset,
		Dynamic:   dynamicClient,
		Context:   contextName,
		Namespace: namespace,
	}, nil
}

// WaitForDeployment waits for a deployment to be ready
func (cm *ClusterManager) WaitForDeployment(namespace, name string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return wait.PollImmediate(2*time.Second, timeout, func() (bool, error) {
		deployment, err := cm.Client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		// Check if deployment is ready
		return deployment.Status.ReadyReplicas == *deployment.Spec.Replicas, nil
	})
}

// WaitForPod waits for a pod to be ready
func (cm *ClusterManager) WaitForPod(namespace, name string, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return wait.PollImmediate(2*time.Second, timeout, func() (bool, error) {
		pod, err := cm.Client.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		// Check if pod is ready
		for _, cond := range pod.Status.Conditions {
			if cond.Type == corev1.PodReady && cond.Status == corev1.ConditionTrue {
				return true, nil
			}
		}
		return false, nil
	})
}

// GetPods returns pods matching the given label selector
func (cm *ClusterManager) GetPods(namespace string, labels map[string]string) ([]corev1.Pod, error) {
	ctx := context.Background()

	labelSelector := metav1.LabelSelector{MatchLabels: labels}
	listOptions := metav1.ListOptions{
		LabelSelector: metav1.FormatLabelSelector(&labelSelector),
	}

	podList, err := cm.Client.CoreV1().Pods(namespace).List(ctx, listOptions)
	if err != nil {
		return nil, err
	}

	return podList.Items, nil
}

// GetDeployment returns a deployment by name
func (cm *ClusterManager) GetDeployment(namespace, name string) (*appsv1.Deployment, error) {
	ctx := context.Background()
	return cm.Client.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
}

// GetResource gets a resource using the dynamic client
func (cm *ClusterManager) GetResource(gvr schema.GroupVersionResource, namespace, name string) (*unstructured.Unstructured, error) {
	ctx := context.Background()

	if namespace == "" {
		return cm.Dynamic.Resource(gvr).Get(ctx, name, metav1.GetOptions{})
	}

	return cm.Dynamic.Resource(gvr).Namespace(namespace).Get(ctx, name, metav1.GetOptions{})
}

// ListResources lists resources using the dynamic client
func (cm *ClusterManager) ListResources(gvr schema.GroupVersionResource, namespace string, labelSelector string) (*unstructured.UnstructuredList, error) {
	ctx := context.Background()

	listOptions := metav1.ListOptions{}
	if labelSelector != "" {
		listOptions.LabelSelector = labelSelector
	}

	if namespace == "" {
		return cm.Dynamic.Resource(gvr).List(ctx, listOptions)
	}

	return cm.Dynamic.Resource(gvr).Namespace(namespace).List(ctx, listOptions)
}

// CreateNamespace creates a namespace
func (cm *ClusterManager) CreateNamespace(name string) error {
	ctx := context.Background()

	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	_, err := cm.Client.CoreV1().Namespaces().Create(ctx, ns, metav1.CreateOptions{})
	return err
}

// DeleteNamespace deletes a namespace
func (cm *ClusterManager) DeleteNamespace(name string) error {
	ctx := context.Background()

	deletePolicy := metav1.DeletePropagationForeground
	return cm.Client.CoreV1().Namespaces().Delete(ctx, name, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	})
}

// NamespaceExists checks if a namespace exists
func (cm *ClusterManager) NamespaceExists(name string) (bool, error) {
	ctx := context.Background()

	_, err := cm.Client.CoreV1().Namespaces().Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		if metav1.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

// WaitForNamespaceDeletion waits for a namespace to be deleted
func (cm *ClusterManager) WaitForNamespaceDeletion(name string, timeout time.Duration) error {
	return wait.PollImmediate(2*time.Second, timeout, func() (bool, error) {
		exists, err := cm.NamespaceExists(name)
		if err != nil {
			return false, err
		}
		return !exists, nil
	})
}

// ResourceExists checks if a resource exists
func (cm *ClusterManager) ResourceExists(gvr schema.GroupVersionResource, namespace, name string) (bool, error) {
	_, err := cm.GetResource(gvr, namespace, name)
	if err != nil {
		if metav1.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CountResources counts resources matching the label selector
func (cm *ClusterManager) CountResources(gvr schema.GroupVersionResource, namespace string, labelSelector string) (int, error) {
	list, err := cm.ListResources(gvr, namespace, labelSelector)
	if err != nil {
		return 0, err
	}
	return len(list.Items), nil
}
