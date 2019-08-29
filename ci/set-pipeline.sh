#!/bin/bash

TARGET=${TARGET:-suse.de}
SECRETS_FILE=${SECRETS_FILE:-../../cloudfoundry/secure/concourse-secrets.yml.gpg}

fly -t "${TARGET}" set-pipeline -p eirini-ssh -c pipeline.yaml --load-vars-from=<(${SECRETS_FILE:+gpg --decrypt --batch ${SECRETS_FILE}})
