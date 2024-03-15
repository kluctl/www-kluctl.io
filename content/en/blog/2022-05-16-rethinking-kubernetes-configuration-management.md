
---
title: "Rethinking Kubernetes Configuration Management"
linkTitle: "Rethinking Kubernetes Configuration Management"
date: 2022-05-16
author: Alexander Block (@codablock)
---
One of the big advantages of Kubernetes is its declarative nature when it comes to deployed resources.
You define what state you expect, apply it to Kubernetes and then let it figure out how to get to that state.

At first, this sounds like a big win, which it actually is. It however doesn’t take long for beginners to realize
that simply writing yamls and applying them via “kubectl apply” is not enough.

The beginner will realize pretty early that configuration between different manifests must be synchronized, e.g.
a port defined in a Service must match the port used when accessing the service from another Deployment. Hardcoding
these configuration values will work till some degree, but get messy when your deployment projects get larger.

It will get completely unmanageable when the same resources need to get deployed multiple times, e.g. for multiple
environments (local, test, uat, prod, …) and/or clusters. Copying and modifying all resources is probably the worst
solution one could come up with at this stage.

A poor-man's solution that works a little bit better is to introduce glue bash code, that performs substitutions
and decides which manifests to actually apply. As with most bash scripting, this tends to start with “just a few lines”
and ends up with unmaintainable spaghetti code.

Spoiler: This blog post is meant to introduce and promote Kluctl, which is a new tool and approach to configuration
management in Kubernetes. Skip the next sections about existing tools if you don’t want to read about stuff you probably
already know well enough.

## Multiple tools to the rescue
Luckily, there are multiple solutions available that try to solve the described situation from above. Each of
these tools has a different approach with its pros and cons. Many of the pros/cons are a matter of personal
preference, e.g. text based templating is preferred by some and avoided by others.

Based on the last sentence, I’ll try to list a few tools that I know of and what my personal
opinion/preference/experience with these is.

## Kustomize
[Kustomize](https://kustomize.io/) advertises itself as a “template-free” solution for configuration management.
Instead of templating, it relies on ["bases and overlays"](https://kubernetes.io/docs/tasks/manage-kubernetes-objects/kustomization/#bases-and-overlays)
for composition and configuration of different environments.

Overlays include bases while performing some basic modifications to these (e.g. add a name prefix). More complex
modifications to bases require the use of strategic or json6902 patches.

This approach usually starts clean and is easy to understand and maintain in the beginning. It can however easily
get complicated and messy due to the unnatural project structure. It can become hard to follow the manifests and how/where modifications/patches got introduced.

The biggest pro point for Kustomize is its ease of use and the flat learning curve. There is nearly no need to
learn anything new, as everything is plain YAML and looks very familiar if you already know how to write and
read Kubernetes manifests. Kustomize also offers some very convenient tools, like secret and configmap generators.
Kustomize is also integrated into kubectl and thus does not need any extra installation.

An interesting and revealing observation I made is that even though Kustomize does not offer templating on its own,
other tools that leverage Kustomize (e.g. Flux or ArgoCD) offer features like substitutions to fulfill a recurring
need of their users, loosening the “template-free” nature of Kustomize and signaling a clear demand for templating
with Kustomize.

## Helm
The first thing you see on [helm.sh](https://helm.sh) is: “The package manager for Kubernetes”. This already reveals what
it’s very good for: Package management. If you want to pick a ready-to-use, pre-packaged and well-configurable
application from a huge library of “Charts”, simply go to [hub.helm.sh](https://hub.helm.sh) and look for what you need.

Installation is as simple as two command line calls, “helm repo add …” and “helm install …”. Other commands
then allow you to manage the installation, e.g. upgrade or delete it.

What it’s not good at is composition and configuration management for your infrastructure and applications. Helm
Charts require quite some boilerplate templating if you want to get to the flexibility that is offered by many of
the publicly available charts. At the same time, you usually don’t need that flexibility and end up with templates
that only have very few configuration options.

If you still decide to use Helm for composition and configuration management, you might end up with all the overhead
of Helm without any of the advantages of it. For example, you would have to provide and maintain a Helm repository
and actually perform releases into it, even if such release management doesn’t make sense for you.

Composition of Helm Charts requires the use of “Umbrella Charts”. These are charts that don’t by itself define
resources but instead refer to multiple dependent charts. Configuration of these sub-charts is passed down from
the umbrella chart. The umbrella chart, as all other charts, will also require you to perform release management
and maintenance of a repository.

In my personal opinion, umbrella charts do not feel “natural”. It feels like a solution that only exists because
it is technically possible to do it this way.

Helm and Kustomize can not be used together without some outer glue, e.g. scripting or special tooling. Many solutions
that glue together Helm and Kustomize will however lose some crucial Helm features in-between, e.g. Helm Hooks and
release management.

Using both tools in combination can still be very powerful. For example, it allows you to customize existing public
Helm Charts in ways not offered by the chart maintainers. For example, if a chart does not offer configuration to
add nodeSelectors, you could use a Kustomize patch to add the desired modifications to the resources in question.

(Kustomize has some experimental Helm integration, but the Kustomize developers strongly discourage
use of it in production).

## Flux CD
Flux allows you to implement GitOps style continuous delivery of your infrastructure and applications. It
internally heavily relies on Kustomize and Helm to provide the actual resources.

Flux can also be used for composition and configuration management. It is for example possible to have a single
Kustomization deployment that actually glues together multiple other Kustomizations and HelmReleases.

Flux also “enhances” Kustomize and Helm by allowing different kinds of preprocessing operations, e.g. substitutions
for Kustomize and patches for Helm Charts. These enhancements are what make configuration management much easier
with Flux.

The biggest advantage of Flux (and GitOps in general) is however also one of its biggest disadvantages, at least
in my humble opinion. You become completely dependent on the Flux infrastructure, even for testing of simple
changes to your deployments. I assume every DevOps engineer knows the pain of recurring
“modify -> push -> wait -> error -> repeat” cycles. Such trial and error sessions can become quite
time-consuming and frustrating.

It would be much easier to test out your changes locally (e.g. via kind or minikube) and push to Git only when
you feel confident enough that your changes are ready. This is not easily possible to do in a reliable way with
Flux based configuration management. Even a dedicated environment/playground stays painful as it does not remove
the described cycle.

Plain Kustomize and Helm would allow you to run the tools from your local machine, targeted against a local cluster
or remote playground environment, but this would require you to replicate some of the features offered by Flux so
that your testing stays compatible with what is later deployed via Flux.

## Argo CD
Argo CD is another GitOps implementation that is similar to Flux. The most visible difference is that it also
offers some very good UI to visualize your deployments and the state of these. I do not have enough experience
with ArgoCD to go into detail, but I assume that many of the things I wrote about Flux also apply for ArgoCD.

## Rethinking Kubernetes Configuration Management
I propose to rethink how the available tooling is used in daily business. Basically, I suggest that the previously
described tools keep doing the parts that they are best at and introduce new tooling to solve the configuration
management problem.

Basically, we need a tool that acts as “glue” between Kustomize, Helm and GitOps. Something that allows to perform
declarative composition of self-written and third-party Kustomize deployments, third-party Helm Charts and a powerful
but easy to learn mechanism surrounding it to synchronize configuration between all components.

It would get even better if a unified CLI would manage all deployments the exact same way, no matter how large or
complex. A CLI that does not interfere with GitOps style continues delivery and allows friendly co-existence.

## Introducing Kluctl
My proposed solution for all this is Kluctl. It’s “the missing glue” (hence the name, glue with k and without e).
Kluctl was born out of the need to create a unified solution that removes the need for all other glue around
Kubernetes deployments. Development of it happened based on large real-life production environments inside a
corporate environment with practicability being one of the highest priorities.

It provides a powerful but easy to learn project structure and templating system that feels natural when used
daily. The project structure follows a simple inclusion hierarchy and gives you full freedom on how to organize
your project.

The templating engine allows you to load configuration from multiple sources (usually arbitrary structured yaml
files) and use the configuration in all components of your project.

The project is then deployed to “targets”, which can represent an environment and/or cluster. Multiple targets
can exist on the same cluster or on different clusters, it’s all up to you. The “target” acts as the entry point
for your configuration, which then implicitly defines what a target means.

The Kluctl CLI allows to diff, deploy, prune and delete targets, all via a unified CLI. If you remember how to
do a “deploy”, you will easily figure out how to diff or prune a target. If you remember how to do it for project
X, you will also know how to do it for project Y. It’s always the same from the CLI perspective, no matter how
simple or complex a deployment actually is.

Kluctl will perform a dry-run and show a diff before a deployment is actually performed. Only if you agree with
what would be changed, the changes are actually applied. After the deployment, Kluctl will show another diff that
visualizes what has been changed.

This means, you always know what will happen and what has happened. It allows you to trust your deployments, even
if you somehow lost track about the state of a target environment.

In addition, Kluctl works according to the premise: Everything can, nothing must. That's why it works equally well
for small and large deployment projects compared to the overhead some other tools impose.

Flux integration is possible via the flux-kluctl-controller and allows friendly co-existence of classical DevOps
and GitOps. This means you can deploy to a playground environment from your local machine and let GitOps handle all
other deployments, adhering to the GitOps principles.

## Learning Kluctl
If you want to learn about Kluctl, go to [kluctl.io](https://kluctl.io) and read the documentation. Best is actually
to start with the Microservices Demo Tutorial, as it tries to introduce you to the Kluctl concepts step by step.

After you’ve finished the tutorial, you should be able to understand how projects are structured, how Kustomize and
Helm is integrated and how targets are used to implement multi-env and multi-cluster deployments.

Based on that, you should be able to imagine how Kluctl could play a role in your own projects and maybe you decide
to give it a try.

More tutorials, documentation and blog posts will follow soon.
