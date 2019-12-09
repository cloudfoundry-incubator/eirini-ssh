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
	serviceName := os.Getenv("OPERATOR_SERVICE_NAME")
	webhookNamespace := os.Getenv("OPERATOR_WEBHOOK_NAMESPACE")
	register := os.Getenv("EIRINI_EXTENSION_REGISTER_ONLY")
	start := os.Getenv("EIRINI_EXTENSION_START_ONLY")
	startOnly := start == "true"
	registerOnly := register == "true"

	fmt.Println("Listening on ", host, port)

	RegisterWebhooks := true
	if startOnly {
		fmt.Println("start-only supplied, the extension will start without registering")
		RegisterWebhooks = false
	}

	filterEiriniApps := true
	x := eirinix.NewManager(
		eirinix.ManagerOptions{
			Namespace:           ns,
			Host:                host,
			Port:                port,
			KubeConfig:          os.Getenv("KUBECONFIG"),
			FilterEiriniApps:    &filterEiriniApps,
			OperatorFingerprint: "eirini-ssh",
			ServiceName:         serviceName,
			WebhookNamespace:    webhookNamespace,
			RegisterWebHook:     &RegisterWebhooks,
		})

	x.AddExtension(&Extension{Namespace: ns})
	x.AddWatcher(&CleanupWatcher{})

	if registerOnly {
		fmt.Println("Registering the extension")
		err := x.RegisterExtensions()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		return
	}

	var v chan error
	go func() {
		fmt.Println("Starting watcher")
		err := x.Watch()
		if err != nil {
			v <- err
			fmt.Println(err.Error())
			return
		}
	}()
	go func() {
		fmt.Println("Starting extension")
		err := x.Start()
		if err != nil {
			v <- err
			fmt.Println(err.Error())
			return
		}
	}()

	for err := range v {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {
	startExtension()
}
