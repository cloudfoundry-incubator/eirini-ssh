# eirini-logging

[Eirinix](https://github.com/SUSE/eirinix) extension for Application Logging in Cloud Foundry

**Note** This is a work in progress to replace fluentd in Eirini

## Concept

fluentd needs access to the container runtime directory on the underlying host in order to read the logs from the running containers.
This is a problem because we need to know what the underlying directory is, so we need to know what the container runtime is (dockerd, cri-o, containerd etc).

This plugin takes a different approach. A [mutating webhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/) listens for new pods on the Eirini namespace and adds a sidecar container to each one upon creation. It also creates an RBAC role that allows the sidecar to read logs from the kube API for this pod. The role is attached to the serviceaccount which will be used to talk to the API (configurable by the operator).

The sidecar reads logs from the kube API and streams them to the [doppler endpoint](https://docs.cloudfoundry.org/loggregator/architecture.html) acting as a loggregator agent [consuming the go-loggregator Golang library](https://github.com/cloudfoundry/go-loggregator). In order for the logs to appear as application logs in when `cf logs` command is used, it sets up the metadata accordingly.

## Build

To compile the code run:

    $> make build

You can also build the Docker image that is used by the injected sidecar container:

    $> DOCKER_ORG="test/" make image

Or define the result image name:

    $> DOCKER_IMAGE="eirini-sidecar" make image

And later on consume it when running the extension

    $> DOCKER_SIDECAR_IMAGE="eirini-sidecar" ./binaries/eirini-logging

the same image can be used to create a pod that acts as the webhook server from withing the kubernetes cluster since it's the same binary that runs as the webhook and as the loggregator agent.

## Run the extension

The extension is listening by default to 10.0.2.2, you can tweak that by setting `EIRINI_EXTENSION_HOST`. This is the IP address of the service in front of the extension pod. You can also set a listening port specifying the environment variable `EIRINI_EXTENSION_PORT`, and the namespace the extension will monitor with `EIRINI_EXTENSION_NAMESPACE`.

If you want to create the extension as a pod inside the cluster you can use the example files here [extension.yaml](contrib/kube/extension.yaml). Make sure you edit this to point to the doppler endpoint and use the correct secret that contains the certificate to talk to doppler.

## Deploy and app and see the logs

If you have [scf](https://github.com/suse/scf) deployed with Eirini you simply have to push an application with `cf push`.

NOTE: For now staging logs do not work because staging happens partly inside Init containers and the sidecar pattern will not work with Init containers because they run in sequence and up to completion.

If you don't have an scf cluster with Eirini but want to debug the extension, you can create a pod that acts as a fake Eirini application. You can find an example of a "fake" Eirini app here: [eirini-fake-app](contrib/kube/eirini_app.yaml). Since you don't have scf runningyou don't have a Doppler endpoint either so there is nowhere to send the logs. You can still debug the extension using `kubectl exec` on the extension pod so it might be useful during development.

## Cleanup

Currently in some cases when pods don't start or stop succesfully, there might be leftovers (Role, Rolebindings etc) in the Eirini namespace. You might have to delete those manually until this is fixed.

## Useful links

- [Eirinix library](https://github.com/SUSE/eirinix)
- [Eirinix sample extension](https://github.com/SUSE/eirinix-sample#eirinix-sample)
- [How to deploy scf with Eirini](https://github.com/SUSE/scf/wiki/Eirini)
- [How to deploy scf on `kind`](https://github.com/SUSE/scf/wiki/scf-on-kind)
