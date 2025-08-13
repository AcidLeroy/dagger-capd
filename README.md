# Dagger + CAPD Experiment 

The goal of this experiment is to setup CAPD to run with in Dagger. This is useful when you want to run e2e tests that require cluster api. 

## Running
Requirements: 
- Dagger 
- Docker

```shell
dagger call --socket=unix:///var/run/docker.sock dev-container --docker-file ./ terminal
```

This will stand up a k3s cluster and attempt to get a CAPD cluster up. 

From there, you can inspect the logs of CAPD:
```shell
dagger / $ kubectl logs -n capd-system capd-controller-manager-6cf4bf67d9-pzkfr  | grep "cluster is not reachable"
E0812 20:24:04.098196       1 cluster_accessor.go:262] "Connect failed" err="error creating HTTP client and mapper: cluster is not reachable: Get \"https://172.18.0.2:6443/?timeout=5s\": context deadline exceeded" controller="clustercache" controllerGroup="cluster.x-k8s.io" controllerKind="Cluster" Cluster="test-cluster/test-cluster" namespace="test-cluster" name="test-cluster" reconcileID="e2ffa769-ffeb-4891-a7ca-339264f6d84e"
E0812 20:24:39.105088       1 cluster_accessor.go:262] "Connect failed" err="error creating HTTP client and mapper: cluster is not reachable: Get \"https://172.18.0.2:6443/?timeout=5s\": context deadline exceeded" controller="clustercache" controllerGroup="cluster.x-k8s.io" controllerKind="Cluster" Cluster="test-cluster/test-cluster" namespace="test-cluster" name="test-cluster" reconcileID="1c31a96b-88ac-4155-8455-90baafb0df16"
``` 


## Problem 
Currently, the networking is broken between the k3s and CAPD workload cluster. This is the ongoing issue that needs to be solved. 