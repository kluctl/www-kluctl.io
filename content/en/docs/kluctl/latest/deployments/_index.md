---
description: Deployments and sub-deployments.
github_repo: https://github.com/kluctl/kluctl
lastmod: "2023-12-18T17:40:02+01:00"
linkTitle: Deployments
path_base_for_github_subdir:
    from: .*
    to: main/docs/kluctl/deployments/README.md
title: Deployments
weight: 30
---



A deployment project is collection of deployment items and sub-deployments. Deployment items are usually
[Kustomize](./kustomize.md) deployments, but can also integrate [Helm Charts](./helm.md).

## Basic structure

The following visualization shows the basic structure of a deployment project. The entry point of every deployment
project is the `deployment.yaml` file, which then includes further sub-deployments and kustomize deployments. It also
provides some additional configuration required for multiple kluctl features to work as expected.

As can be seen, sub-deployments can include other sub-deployments, allowing you to structure the deployment project
as you need.

Each level in this structure recursively adds [tags](./tags.md) to each deployed resources, allowing you to control
precisely what is deployed in the future.

```
-- project-dir/
   |-- deployment.yaml
   |-- .gitignore
   |-- kustomize-deployment1/
   |   |-- kustomization.yaml
   |   `-- resource.yaml
   |-- sub-deployment/
   |   |-- deployment.yaml
   |   |-- kustomize-deployment2/
   |   |   |-- resource1.yaml
   |   |   `-- ...
   |   |-- kustomize-deployment3/
   |   |   |-- kustomization.yaml
   |   |   |-- resource1.yaml
   |   |   |-- resource2.yaml
   |   |   |-- patch1.yaml
   |   |   `-- ...
   |   |-- kustomize-with-helm-deployment
   |   |   |-- charts/
   |   |   |   `-- ...
   |   |   |-- kustomization.yaml
   |   |   |-- helm-chart.yaml
   |   |   `-- helm-values.yaml
   |   `-- subsub-deployment/
   |       |-- deployment.yaml
   |       |-- ... kustomize deployments
   |       `-- ... subsubsub deployments
   `-- sub-deployment/
       `-- ...
```

## Order of deployments
Deployments are done in parallel, meaning that there are usually no order guarantees. The only way to somehow control
order, is by placing [barriers](./deployment-yml.md#barriers) between kustomize deployments.
You should however not overuse barriers, as they negatively impact the speed of kluctl.

## Plain Kustomize

It's also possible to use Kluctl on plain Kustomize deployments. Simply run `kluctl deploy` from inside the
folder of your `kustomization.yaml`. If you also don't have a `.kluctl.yaml`, you can also work without targets.

Please note that pruning and deletion is not supported in this mode.