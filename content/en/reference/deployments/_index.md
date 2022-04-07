---
title: "Deployments"
linkTitle: "Deployments"
weight: 2
description: >
    Deployments and sub-deployments.
---

A deployment project is collection of deployment items and sub-deployments. Deployment items are usually
[Kustomize]({{< ref "./kustomize" >}}) deployments, but can also integrate [Helm Charts]({{< ref "./helm" >}}).

## Basic structure

The following visualization shows the basic structure of a deployment project. The entry point of every deployment
project is the `deployment.yml` file, which then includes further sub-deployments and kustomize deployments. It also
provides some additional configuration required for multiple kluctl features to work as expected.

As can be seen, sub-deployments can include other sub-deployments, allowing you to structure the deployment project
as you need.

Each level in this structure recursively adds [tags]({{< ref "./tags" >}}) to each deployed resources, allowing you to control
precisely what is deployed in the future.

Some visualized files/directories have links attached, follow them to get more information.

<pre>
-- project-dir/
   |-- <a href="{{< ref "./deployment-yml" >}}">deployment.yml</a>
   |-- .gitignore
   |-- kustomize-deployment1/
   |   |-- kustomization.yml
   |   `-- resource.yml
   |-- sub-deployment/
   |   |-- deployment.yml
   |   |-- kustomize-deployment2/
   |   |   |-- kustomization.yml
   |   |   |-- resource1.yml
   |   |   `-- ...
   |   |-- kustomize-deployment3/
   |   |   |-- kustomization.yml
   |   |   |-- resource1.yml
   |   |   |-- resource2.yml
   |   |   |-- patch1.yml
   |   |   `-- ...
   |   |-- <a href="{{< ref "./helm" >}}">kustomize-with-helm-deployment/</a>
   |   |   |-- charts/
   |   |   |   `-- ...
   |   |   |-- kustomization.yml
   |   |   |-- helm-chart.yml
   |   |   `-- helm-values.yml
   |   `-- subsub-deployment/
   |       |-- deployment.yml
   |       |-- ... kustomize deployments
   |       `-- ... subsubsub deployments
   `-- sub-deployment/
       `-- ...
</pre>

## Order of deployments
Deployments are done in parallel, meaning that there are usually no order guarantees. The only way to somehow control
order, is by placing [barriers]({{< ref "./deployment-yml#barriers" >}}) between kustomize deployments.
You should however not overuse barriers, as they negatively impact the speed of kluctl.
