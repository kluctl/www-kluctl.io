---
title: "Simple"
linkTitle: "simple"
weight: 1
description: >
    Very simple example with cluster and deployment in a single repository.
---
## Description
This example is a very simple one that shows how to define a target cluster, context, create a
namespace and deploy a nginx. You can configure the name of the namespace by changing the arg `environment` in
[.kluctl.yml](https://github.com/kluctl/kluctl-examples/blob/main/simple/.kluctl.yml).

## Prerequisites
1) A running [kind](https://kind.sigs.k8s.io/) cluster with a context named `kind-kind`.
2) Of course, you need to install kluctl. Please take a look at the 
[installation guide]({{< ref "docs/latest/kluctl/installation" >}}), in case you need further information.

## How to deploy
```bash
git clone git@github.com:kluctl/kluctl-examples.git
cd kluctl-examples/simple
kluctl diff --target simple
kluctl deploy --target simple
```