ARG BASE_IMAGE=opensuse/leap
FROM golang:1.12 as build
ARG user="SUSE CFCIBot"
ARG email=ci-ci-bot@suse.de
ADD . /eirini-ssh
WORKDIR /eirini-ssh
RUN git config --global user.name ${user}
RUN git config --global user.email ${email}
RUN GO111MODULE=on go mod vendor
RUN bin/build-extension
RUN bin/build-ssh-proxy

FROM bitnami/kubectl as kubectl

FROM $BASE_IMAGE
COPY --from=kubectl /opt/bitnami/kubectl/bin/kubectl /bin/kubectl
COPY --from=build /eirini-ssh/binaries/ssh-extension /bin/eirini-ssh-extension
COPY --from=build /eirini-ssh/binaries/ssh-proxy /bin/eirini-ssh-proxy

