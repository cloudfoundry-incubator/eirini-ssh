package authenticators

import (
	cfauth "code.cloudfoundry.org/diego-ssh/authenticators"
	"code.cloudfoundry.org/lager"
	"golang.org/x/crypto/ssh"
)

type kubeBuilder struct {
}

func NewKubeAuth() cfauth.PermissionsBuilder {
	return &kubeBuilder{}
}

func (kb *kubeBuilder) Build(logger lager.Logger, processGuid string, index int, metadata ssh.ConnMetadata) (*ssh.Permissions, error) {
	return nil, nil
}
