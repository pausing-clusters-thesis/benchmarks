#!/usr/bin/env bash

set -euExo pipefail
shopt -s inherit_errexit

gcloud container \
clusters create 'rzetelskik-proxy-csi-driver-benchmarks' \
--zone='europe-west3-a' \
--cluster-version="1.31" \
--machine-type='n2-standard-4' \
--num-nodes='1' \
--disk-type='pd-ssd' --disk-size='20' \
--local-nvme-ssd-block='count=1' \
--image-type='UBUNTU_CONTAINERD' \
--enable-stackdriver-kubernetes \
--no-enable-autoupgrade \
--no-enable-autorepair


