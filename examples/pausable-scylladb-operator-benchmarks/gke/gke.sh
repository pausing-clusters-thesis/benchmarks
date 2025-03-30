#!/usr/bin/env bash

set -euExo pipefail
shopt -s inherit_errexit

gcloud container \
clusters create 'rzetelskik-pausing-benchmarks' \
--zone='europe-west3-a' \
--cluster-version="1.31" \
--machine-type='n2-standard-4' \
--num-nodes='1' \
--disk-type='pd-standard' --disk-size='20' \
--image-type='UBUNTU_CONTAINERD' \
--enable-stackdriver-kubernetes \
--no-enable-autoupgrade \
--no-enable-autorepair

gcloud container \
node-pools create 'scylladb' \
--zone='europe-west3-a' \
--cluster='rzetelskik-pausing-benchmarks' \
--node-version="1.31" \
--machine-type='c2d-standard-2' \
--num-nodes='3' \
--disk-type='pd-ssd' --disk-size='20' \
--local-nvme-ssd-block='count=1' \
--image-type='UBUNTU_CONTAINERD' \
--system-config-from-file='systemconfig.yaml' \
--no-enable-autoupgrade \
--no-enable-autorepair \
--node-labels='scylla.scylladb.com/node-type=scylla' \
--node-taints='scylla-operator.scylladb.com/dedicated=scyllaclusters:NoSchedule'

