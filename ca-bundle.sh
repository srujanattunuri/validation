#!/bin/bash
set -o errexit
set -o nounset
set -o pipefail

ROOT=$(cd $(dirname $0)/../../; pwd)

export CA_BUNDLE=$(kubectl get configmap -n kube-system extension-apiserver-authentication -o=jsonpath='{.data.client-ca-file}' | base64 | tr -d '\n')

sed -i "s/caBundle: .*$/caBundle: ${CA_BUNDLE}/g" /Users/srujanaattunuri/github/validation/manifest.yaml