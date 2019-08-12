package authenticators

import (
	"encoding/json"
	"fmt"

	"code.cloudfoundry.org/bbs"
	"code.cloudfoundry.org/bbs/models"
	"code.cloudfoundry.org/diego-ssh/proxy"
	"code.cloudfoundry.org/diego-ssh/routes"
	"code.cloudfoundry.org/lager"
	"golang.org/x/crypto/ssh"
)

type kubeBuilder struct {
}

func NewKubeAuth() PermissionsBuilder {
	return &kubeBuilder{}
}

func (pb *permissionsBuilder) Build(logger lager.Logger, processGuid string, index int, metadata ssh.ConnMetadata) (*ssh.Permissions, error) {
	actual, err := pb.bbsClient.ActualLRPGroupByProcessGuidAndIndex(logger, processGuid, index)
	return nil, err
}
