---
title: "Simple with external repositories"
linkTitle: "simple-with-external-repos"
weight: 2
description: >
    Very simple example with cluster and deployment in external repositories.
---
## Description
This example is very similar to `simple` except that the
target cluster and the deployment is defined externally. You can configure the repositories and the ref in
[.kluctl.yml](https://github.com/kluctl/kluctl-examples/blob/main/simple-with-external-repos/.kluctl.yml).

## Prerequisites
1) A running [kind](https://kind.sigs.k8s.io/) cluster with a context named `kind-kind`.
2) Of course, you need to install kluctl. Please take a look at the
   [installation guide]({{< ref "docs/installation" >}}), in case you need further information.

## How to deploy
```bash
git clone git@github.com:kluctl/kluctl-examples.git
cd kluctl-examples/simple-with-external-repos
kluctl diff --target simple-with-external-repos
kluctl deploy --target simple-with-external-repos
```