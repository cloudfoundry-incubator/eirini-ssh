package authenticators

import (
	"encoding/base64"
	"errors"
	"strconv"
	v1 "k8s.io/api/core/v1"
	"code.cloudfoundry.org/diego-ssh/proxy"
	cfauth "code.cloudfoundry.org/diego-ssh/authenticators"
	"code.cloudfoundry.org/lager"
	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type kubeBuilder struct {
	Kubeconfig string
	Kubeclient client.Client
}

type ConnectionOptions struct {
	PodName, PublicKey, PrivateKey, Fingerprint string
}

func NewKubeAuth(kubeconfig string) cfauth.PermissionsBuilder {
	return &kubeBuilder{}
}

func (kb *kubeBuilder) Build(logger lager.Logger, processGuid string, index int, metadata ssh.ConnMetadata) (*ssh.Permissions, error) {
	conn, err := kb.GetSecretKeys(processGuid + "-" + strconv.Itoa(int) + "-ssh-key-meta")
	if err != nil {
		return nil, errors.Wrap(err, "Could not retreive SSH key for "+processGuid)
	}

	/ Using a typed object.
	pod := &v1.Pod{}
	_ = kb.Kubeclient.Get(context.TODO(), client.ObjectKey{
		Namespace: "eirini",  // FIXME: Hardcoded 
		Name:      conn.PodName,
	}, pod)	
	
	logMessage := fmt.Sprintf("Successful remote access by %s", metadata.RemoteAddr().String())
	
	address:=pod.Status.HostIP  // FIXME: inject containerports if possible ?
	port:=2222 // FIXME: Hardcoded 

	var targetConfig *proxy.TargetConfig
	targetConfig = &proxy.TargetConfig{
		Address:             fmt.Sprintf("%s:%d", address, port),
		TLSAddress:          tlsAddress,
		ServerCertDomainSAN: actual.ActualLRPInstanceKey.InstanceGuid,
		HostFingerprint:     sshRoute.HostFingerprint,
		User:                sshRoute.User,
		Password:            sshRoute.Password,
		PrivateKey:          sshRoute.PrivateKey,
	}

	targetConfigJson, err := json.Marshal(targetConfig)
	if err != nil {
		return nil, err
	}


	logMessageJson, err := json.Marshal(proxy.LogMessage{
		Message: logMessage,
		// Tags:    tags, // FIXME: todo
	})
	if err != nil {
		return nil, err
	}

	return &ssh.Permissions{
		CriticalOptions: map[string]string{
			"proxy-target-config": string(targetConfigJson),
			"log-message":         string(logMessageJson),
		},
	}, nil

}

func (kb *kubeBuilder) GetSecretKeys(name string) (ConnectionOptions, error) {

	// We have to query for the Secret using an unstructured object because the cache for the structured
	// client is not initialized yet at this point in time. See https://github.com/kubernetes-sigs/controller-runtime/issues/180
	secret := &unstructured.Unstructured{}
	secret.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   "",
		Kind:    "Secret",
		Version: "v1",
	})
	kb.Kubeclient.Get(ctx, secretNamespacedName, secret)
	if secret.GetName() == "" {
		return ConnectionOptions{}, errors.New("No secret found")
	}

	data := secret.Object["data"].(map[string]interface{})
	pubKey, err := base64.StdEncoding.DecodeString(data["public_key"].(string))
	if err != nil {
		return ConnectionOptions{}, errors.New("Failed decoding 'public_key' from the secret")
	}
	privKey, err := base64.StdEncoding.DecodeString(data["private_key"].(string))
	if err != nil {
		return ConnectionOptions{}, errors.New("Failed decoding 'private_key' from the secret")
	}
	fingerprint, err := base64.StdEncoding.DecodeString(data["fingerprint"].(string))
	if err != nil {
		return ConnectionOptions{}, errors.New("Failed decoding 'fingerprint' from the secret")
	}
	podName, err := base64.StdEncoding.DecodeString(data["pod_name"].(string))
	if err != nil {
		return ConnectionOptions{}, errors.New("Failed decoding 'pod_name' from the secret")
	}

	return ConnectionOptions{
		PublicKey:   pubKey,
		PrivateKey:  privKey,
		Fingerprint: fingerprint,
		PodName:     podName,
	}, nil
}
