
---
title: "Kluctl And Flux"
linkTitle: "Kluctl And Flux"
date: 2022-05-05
author: Alexander Block (@codablock)
---

We're very happy to announce that Kluctl can from now on be used together with [Flux](https://fluxcd.io/). This
will allow you to combine the workflows and features advertised by Kluctl with GitOps style continuous delivery.

## GitOps vs Kluctl
One of the first questions that we usually get when introducing Kluctl to someone is something like:
"Why not GitOps?" or "Why not Flux?". There seems to be a common misunderstanding that arises in many people
when trying to understand Kluctl on first sight, which is to believe that Kluctl is an alternative or competitor 
to GitOps and Flux (or even ArgoCD).

This is not the case. If one wants to compare Kluctl with something else, then it's more appropriate to compare it
to Helm, Kustomize or Helmfile. It should be clear that Kustomize for example is not an alternative/competitor for Flux,
but instead an essential tool and building block to make it work.

Kluctl can be looked at from the same perspective when it comes to Flux. Flux implements Helm and Kustomize support
via different controllers, namely the [kustomize-controller](https://fluxcd.io/docs/components/kustomize/) and the
[helm-controller](https://fluxcd.io/docs/components/helm/). Kluctl can simply do the same to support Flux.

## Introducing the Kluctl Flux Controller
An alpha version of the [Kluctl Flux Controller]({{< ref "docs/flux" >}}) has just been released. It allows to
create [KluctlDeployment]({{< ref "docs/flux/kluctldeployment" >}}) objects which are reconciled in a similar
fashion as [Kustomizations](https://fluxcd.io/docs/components/kustomize/kustomization/).

Each KluctlDeployment specifies a source object (e.g. a [GitRepository](https://fluxcd.io/docs/components/source/gitrepositories/)),
the [target]({{< ref "docs/reference/kluctl-project/targets" >}}) to be deployed and some information on how
to handle kubeconfigs. The controller then regularly reconciles the deployment.

For a simple example, check [this]({{< ref "docs/flux/controller" >}}) documentation.

## Kustomize/Helm vs Kluctl
If you've already read through the [Kluctl documentation]({{< ref "docs" >}}), you've probably already noticed
that Kluctl internally uses Kustomize and Helm extensively.

This might raise the question: Why not use plain Kustomize and/or Helm if Flux is already involved? There are multiple
reasons to prefer Kluctl, which we'll try to elaborate here:

## Kluctl Projects/Deployments
If you prefer the way Kluctl organizes and structures projects and deployments, then using the Flux Kluctl Controller
is obviously the best choice. Kluctl allows you to easily glue together what belongs together. If for example, a redis
database is required to make your application work, you can manage the redis Helm Release and your application in the
same deployment, including the necessary configuration to let them talk to each other.

To see how different a Kluctl deployment is compared to classic Kustomize/Helm + Flux, you can compare the
[flux2-kustomize-helm-example](https://github.com/fluxcd/flux2-kustomize-helm-example) and the
[Kluctl Microservices Demo](https://github.com/kluctl/kluctl-examples/tree/main/microservices-demo/3-templating-and-multi-env)
([here]({{< ref "docs/guides/tutorials/microservices-demo">}}) is tutorial for the demo).

## Native multi-env support
Kluctl allows you to natively create deployment projects that can be deployed multiple times to different
environments/targets. You can for example have one target that is solely meant for `local` (e.g. Kind based) deployments,
one that targets the `test` environment and one for `prod`. You can then use templating to influence deployments in whatever
way you like. For example, you could change the `local` target to set all replicas to 1 and omit resource hungry
support applications.

This is possible in plain Kustomize as well, but requires you to solve it without the concept of
[targets]({{< ref "docs/reference/kluctl-project/targets" >}}) and without templating. In Kustomize, multi-env
deployments must be solved with [overlays](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/#bases-and-overlays),
which does not necessary align with how you prefer your project structure.

## Mix local DevOps and GitOps
The core idea of GitOps is that Git becomes the single source of truth for the desired cluster state. This is something
that is extremely valuable with many advantages compared to other approaches. There are however still situations
where diverging from GitOps is very valuable as well.

For example, when you start a new deployment project, you're usually in a state of frequent changes inside the deployment
project. These frequent changes need frequent deployments and testing until you get to a point where things are stable
enough. If you're forced to adhere to GitOps in that situation, you end up with very noisy Git histories and plenty
of trial-and-error deployment cycles. This is, at least in our opinion and experience, a major productivity killer
without any advantages.

With Kluctl, you can start developing locally and deploying from your local machine, with the guarantee that what you
see is what will also happen later when GitOps is introduced for the same deployment. When you're "ready", you can then
commit and push to Git and create the appropriate `KluctlDeployment` objects and then let GitOps/Flux take over.

You can also use dedicated targets for development purposes and only deploy to them from your local machine, while
other targets are deployed via GitOps/Flux.
