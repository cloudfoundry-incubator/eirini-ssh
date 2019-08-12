FROM bitnami/kubectl as kubectl

FROM opensuse:leap
COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /bin/kubectl
ADD binaries/eirini-logging /bin/eirini-logging
ENTRYPOINT ["/bin/eirini-logging"]
