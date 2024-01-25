---
title: "Simple Helm"
linkTitle: "simple-helm"
weight: 3
description: >
    Very simple example of a helm-based deployment.
---
## Description
This example is very similar to `simple` but it deploys a Helm-based nginx to
give a first impression how kluctl and Helm work together.

## Prerequisites
1) A running [kind](https://kind.sigs.k8s.io/) cluster with a context named `kind-kind`.
2) Of course, you need to install kluctl. Please take a look at the
   [installation guide]({{< ref "docs/kluctl/installation" >}}), if you need further information.
3) You also need to install [Helm](https://helm.sh/). Please take a look at the 
[Helm installation guide](https://helm.sh/docs/intro/install/) for further information. 

## How to deploy
```bash
git clone git@github.com:kluctl/kluctl-examples.git
cd kluctl-examples/simple-helm
kluctl helm-pull
kluctl diff --target simple-helm
kluctl deploy --target simple-helm
```