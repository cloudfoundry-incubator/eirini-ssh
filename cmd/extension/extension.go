package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	inmemorygenerator "code.cloudfoundry.org/cf-operator/pkg/credsgen/in_memory_generator"
	eirinix "github.com/SUSE/eirinix"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission/types"
)

type Extension struct{ Namespace string }

func getVolume(name, path string) (v1.Volume, v1.VolumeMount) {
	mount := v1.VolumeMount{
		Name:      name,
		MountPath: path,
	}

	vol := v1.Volume{
		Name: name,
	}

	return vol, mount
}

func extractInstanceID(s string) string {
	var res string
	el := strings.Split(s, "-")
	if len(el) != 0 {
		res = el[len(el)-1]
		if _, err := strconv.Atoi(res); err == nil {
			return res
		}
	}

	return "0"
}

func (ext *Extension) Handle(ctx context.Context, eiriniManager eirinix.Manager, pod *v1.Pod, req types.Request) types.Response {

	if pod == nil {
		return admission.ErrorResponse(http.StatusBadRequest, errors.New("No pod could be decoded from the request"))
	}

	config, err := eiriniManager.GetKubeConnection()
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, errors.Wrap(err, "Failed getting the Kube connection"))
	}

	podCopy := pod.DeepCopy()

	// Mount the serviceaccount token in the container
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, errors.Wrap(err, "Failed to create a kube client"))
	}
	guid, ok := pod.GetLabels()["guid"]
	if !ok {
		return admission.ErrorResponse(http.StatusBadRequest, errors.New("Couldn't get Eirini APP UID"))
	}

	index := extractInstanceID(pod.Name)
	// TODO:
	// - Create or append to the existing secret a new SSH key for this app
	// - Create a volume and a mount the secrete we created (and only that) as
	//   an environment variable inside the application pod
	// - Cleanup any non-existing application keys from from the secret.
	//   (NOTE: This is not HA! A better approach is to have a Watcher watching for pod deletions on the eirini namespace and remove the
	//   relevant key when a pod is deleted)

	generator := inmemorygenerator.NewInMemoryGenerator(eiriniManager.GetLogger())
	key, err := generator.GenerateSSHKey(podCopy.Name)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, errors.Wrap(err, "Failed to generate SSH key for the application"))
	}
	secretName := guid + "-" + index + "-ssh-key-meta"
	fmt.Println("Creating", secretName)
	newSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: podCopy.Namespace,
		},
		StringData: map[string]string{
			"public_key":  string(key.PublicKey),
			"private_key": string(key.PrivateKey),
			"fingerprint": key.Fingerprint,
			"pod_name":    pod.Name,
		},
	}
	_, err = kubeClient.CoreV1().Secrets(podCopy.Namespace).Create(newSecret)
	if err != nil {
		return admission.ErrorResponse(http.StatusBadRequest, errors.Wrap(err, "Failed to create a kube secret for the application SSH key"))
	}

	for i, c := range podCopy.Spec.Containers {
		c.Env = append(c.Env, v1.EnvVar{
			Name: "EIRINI_SSH_KEY",
			ValueFrom: &v1.EnvVarSource{
				SecretKeyRef: &v1.SecretKeySelector{
					LocalObjectReference: v1.LocalObjectReference{Name: secretName},
					Key:                  "public_key",
				},
			},
		})
		podCopy.Spec.Containers[i] = c
	}

	return admission.PatchResponse(pod, podCopy)
}
