package extension

import (
	"fmt"

	. "code.cloudfoundry.org/eirini-ssh/pkg/logger"
	eirinix "code.cloudfoundry.org/eirinix"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/watch"
	typedv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	// "k8s.io/client-go/kubernetes"
	//	"k8s.io/client-go/rest"
)

type CleanupWatcher struct {
	Manager eirinix.Manager
}

func (pw *CleanupWatcher) Handle(manager eirinix.Manager, e watch.Event) {
	LogDebug("Received event: ", fmt.Sprintf("%+v", e))
	if e.Object == nil {
		return
	}

	pod, ok := e.Object.(*corev1.Pod)
	if !ok {
		LogError("Received non-pod object in watcher channel")
		return
	}

	if e.Type == watch.Deleted {
		secretName, err := generateSecretNameForPod(pod)
		if err != nil {
			LogError(err.Error())
			return
		}
		LogInfo("Removing secret " + secretName + " for pod " + pod.GetName())

		config, err := manager.GetKubeConnection()
		if err != nil {
			LogError(err.Error())
			return
		}
		kubeClient, err := typedv1.NewForConfig(config)
		if err != nil {
			LogError(err.Error())
		}

		// https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#DeleteOptions
		err = kubeClient.Secrets(pod.Namespace).Delete(manager.GetContext(), secretName, metav1.DeleteOptions{})
		if err != nil {
			LogError(err.Error())
		}
		LogInfo("Secret '" + secretName + "' removed from namespace " + pod.Namespace)
	} else {
		LogDebug("Ignoring event of type: ", e.Type)
	}
}
