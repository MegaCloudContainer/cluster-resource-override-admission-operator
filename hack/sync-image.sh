#!/bin/bash

set -e

TARGET_VERSION=${1:-v4.16}

repos=(
     release.daocloud.io
     release-ci.daocloud.io 
)

for repo in ${repos[@]}; do
    skopeo copy docker://ghcr.io/megacloudcontainer/cro-operator:"${TARGET_VERSION}" docker://"${repo}"/openshift/cro-operator:v4.16  --all
    skopeo copy docker://ghcr.io/megacloudcontainer/cro:"${TARGET_VERSION}" docker://"${repo}"/openshift/cro:v4.16  --all
done
