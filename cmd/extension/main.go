package main

import (
	"fmt"
	"os"
	"strconv"

	eirinix "github.com/SUSE/eirinix"
)

func startExtension() {
	var port int32
	ns := os.Getenv("EIRINI_EXTENSION_NAMESPACE")
	if len(ns) == 0 {
		ns = "default"
	}
	host := os.Getenv("EIRINI_EXTENSION_HOST")
	if len(host) == 0 {
		host = "10.0.2.2"
	}
	p := os.Getenv("EIRINI_EXTENSION_PORT")
	if len(p) == 0 {
		port = 3000
	} else {
		po, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
		port = int32(po)
	}

	fmt.Println("Listening on ", host, port)

	filterEiriniApps := true
	x := eirinix.NewManager(
		eirinix.ManagerOptions{
			Namespace:           ns,
			Host:                host,
			Port:                port,
			KubeConfig:          os.Getenv("KUBECONFIG"),
			FilterEiriniApps:    &filterEiriniApps,
			OperatorFingerprint: "eirini-ssh",
		})

	x.AddExtension(&Extension{Namespace: ns})
	fmt.Println(x.Start())
}

func main() {
	startExtension()
}
