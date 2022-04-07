---
title: "Prerequisites"
linkTitle: "Prerequisites"
weight: 1
description: >
    What needs to be done before using kluctl.
---


## Get a Kubernetes cluster

The first step is of course: You need a kubernetes cluster. It doesn't really matter where this cluster is hosted, if
it's a local (e.g. kind) cluster, managed cluster, or a self-hosted cluster, kops or kubespray based, AWS, GCE, Azure, ... and so on. kluctl
is completely independent of how Kubernetes is deployed and where it is hosted.

There is however a minimum Kubernetes version that must be met: 1.20.0. This is due to the heavy use of server-side apply
which was not stable enough in older versions of Kubernetes.

## Prepare your kubeconfig

Your local kubeconfig should be configured to have access to the target Kubernetes cluster via a dedicated context. The context
name should match with the name that you want to use for the cluster from now on. Let's assume the name is `test.example.com`,
then you'd have to ensure that the kubeconfig context `test.example.com` is correctly pointing and authorized for this
cluster.

See [Configure Access to Multiple Clusters](https://kubernetes.io/docs/tasks/access-application-cluster/configure-access-multiple-clusters/) for documentation
on how to manage multiple clusters with a single kubeconfig. Depending on the Kubernets provisioning/deployment tooling
you used, you might also be able to directly export the context into your local kubeconfig. For example,
[kops](https://github.com/kubernetes/kops/blob/master/docs/cli/kops_export.md) is able to export and merge the kubeconfig
for a given cluster.
