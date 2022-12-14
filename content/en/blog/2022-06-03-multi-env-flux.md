
---
title: "Multiple Environments with Flux and Kluctl"
linkTitle: "Multiple Environments with Flux and Kluctl"
date: 2022-06-03
author: Alexander Block (@codablock)
images:
  - "/images/blog/multi-env-flux.jpg"
---

![multi-env-flux](/images/blog/multi-env-flux.jpg)

Most projects that have server-side components usually need to be deployed multiple times, at least if you don't want to
break things for your users all the time. This means, that there is not only one single "prod" environment, but also something
like "test" or "staging" environments. So usually, there is a minimum of two different environments.

Each of these environments has a defined role to play. The "prod" environment is obviously the environment that is used
in production. A "test" or "staging" environment is usually the one where new versions are deployed so that they can
be tested before being promoted to "prod". A "dev" environment might be the one where even high-risk deployments
are allowed where breaking everything is daily business. You could even have multiple dev environments, e.g. one per
developer.

In the Kubernetes world, it is perfectly viable and good practice to run multiple environments on a single cluster.
Depending on the requirements of your project, you might want to separate prod from non-prod environments, but you
could still have "prod" and "staging" on the same cluster while "test", "dev", "uat" exist on a non-prod cluster.

## Environment Parity

If you want to make this work in a way that does not cause too much pain, you must practice environment parity
as much as possible. The best way to achieve this is to have full automation and everything as code. The same "code"
that was used to deploy "prod" should also be used to deploy "staging", "test" and all other environments. The only
difference in the deployments should be a defined set of configurations.

Such configurations can for example be the target cluster and/or namespace, the resource allocations
(e.g. 1 replica in "dev", 3 replicas in "prod"), external systems (e.g. databases) and ingress configuration
(e.g. DNS and certificates). Configuration can also enable/disable conditional parts of the deployment, for example to
disable advanced monitoring on "dev" environments and enable mocking services as replacements for real external systems.

## Tooling

There are multiple tools available that allow you to implement a multi-env/multi-cluster deployment that is completely
automated and completely "as code". Helm and Kustomize are currently the first tools that will pop up when you try to
look for such tools. As written in my [previous blog post]({{< ref "blog/2022-05-16-rethinking-kubernetes-configuration-management" >}}),
I believe that these tools are the best option for the things that they do very good, but a sub-optimal choice when it
comes to configuration management.

Kluctl is the preferred solution for me right now. Not only because I built it, but also because so far I did not find
a solution that is as easy to learn and use and so flexible at the same time.

## Fully working multi-env example

I suggest to open the [microservices demo]({{< ref "docs/guides/tutorials/microservices-demo" >}})
in another tab and look into it at least briefly (especially the third part). I will from now on pick stuff from this
tutorial as examples in this blog post.

## Targets in Kluctl

Kluctl works with the concept of "targets". A target is a named configuration that acts as the entry point for every
further configuration required for your environment. As an example, look at the targets from
[.kluctl.yaml]({{< ref "docs/kluctl/reference/kluctl-project" >}}) of the microservices demo:

```yaml
targets:
  - name: local
    context: kind-kind
    args:
      env_type: local
  - name: test
    context: kind-kind
    args:
      env_type: real
  - name: prod
    context: kind-kind
    args:
      env_type: real
```

Based on these, the same deployment can be configured differently depending on the target, or actually the `args`
passed via the target. In the microservices demo, `env_type` (name can be chosen by you) is used to include further
configuration inside the `deployment.yaml`:

```yaml
...
vars:
  - file: ./vars/{{ args.env_type }}.yml
...
```

At the same time, Kluctl makes the target configuration itself available to the deployment, making things like this
possible:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: ms-demo-{{ target.name }}
```

You can also conditionally enable/disable parts of your deployment depending on the configuration loaded via the `vars`
block above:

```yaml
deployments:
  - path: adservice
  - path: cartservice
  - path: checkoutservice
  - path: currencyservice
  {% if services.emailservice.enabled %}
  - path: emailservice
  {% endif %}
  - path: frontend
  {% if services.loadgenerator.enabled %}
  - path: loadgenerator
  {% endif %}
  - path: paymentservice
  - path: productcatalogservice
  - path: recommendationservice
  - path: shippingservice
```

I hope the above snippets give you a feeling about how multi-env deplyoments can be solved via Kluctl. As already
mentioned, I suggest to read through the [microservices demo]({{< ref "docs/guides/tutorials/microservices-demo" >}})
tutorial to get an even better understanding. The first two parts will describe some Kluctl basics while the third
part enters multi-env deployments.

## GitOps and Flux

Kluctl is designed in a way that allows seamless co-existence of CLI based workflows, classical CI/CD and GitOps. This
means, that even if you decide to perform prod deployments only via GitOps, you can still perform exactly the same
deployment to other environments through your CLI. All this without any struggle (you really only need access to the
cluster) and 100% compatible to how the same deployment would be performed via GitOps or CI/CD.

This allows you to adapt your workflow depending on what your current goal is. For example, a developer testing out
risky and bleeding-edge changes in his personal dev environment can deploy from his local machine, avoiding the
painful "modify -> push -> wait -> error -> repeat" cycles seen too often in pure GitOps and CI/CD setups. When the
developer is done with the changes, GitOps can take over on the another (e.g. "test" and later "prod") environment.

Even for "prod", which in the above scenario is GitOps managed, can benefit from the possibility to run Kluctl from
your local machine. Running a "kluctl diff -t prod" before promoting to "prod" can prevent some scary surprises.

Kluctl implements GitOps via the [flux-kluctl-controller]({{< ref "blog/2022-05-11-introducing-kluctl-and-flux" >}}).
It allows to create [`KluctlDeployment`]({{< ref "docs/flux/kluctldeployment" >}}) objects which refer to your Kluctl
project (which relies in Git) and the target to be deployed.

## Installing flux-kluctl-controller

Before being able to create [`KluctlDeployment`]({{< ref "docs/flux/kluctldeployment" >}}) objects, the
flux-kluctl-controller needs to be installed first. Please navigate to the
[installation documentation](https://github.com/kluctl/flux-kluctl-controller/blob/main/docs/install.md)
and follow the instructions found there.

## Microservices Demo and Flux

Deploying the microservices demo via Flux is quite easy. First, we'll need a
[`GitRepository`](https://fluxcd.io/docs/components/source/gitrepositories/) object that refers to the Kluctl project:

```yaml
apiVersion: source.toolkit.fluxcd.io/v1beta1
kind: GitRepository
metadata:
  name: microservices-demo
spec:
  interval: 1m
  url: https://github.com/kluctl/kluctl-examples.git
  ref:
    branch: main
```

This will cause Flux to pull the source whenever it changes. A [`KluctlDeployment`]({{< ref "docs/flux/kluctldeployment" >}})
can then refer to the `GitRepository`:

```yaml
apiVersion: flux.kluctl.io/v1alpha1
kind: KluctlDeployment
metadata:
  name: microservices-demo-test
spec:
  interval: 5m
  path: "./microservices-demo/3-templating-and-multi-env/"
  sourceRef:
    kind: GitRepository
    name: microservices-demo
  timeout: 2m
  target: test
  prune: true
  # kluctl targets specify the expected context name, which does not necessarily match the context name
  # found while it is deployed via the controller. This means we must pass a kubeconfig to kluctl that has the
  # context renamed to the one that it expects.
  renameContexts:
    - oldContext: default
      newContext: kind-kind
```

The above example will cause the controller to deploy the "test" target/environment to the same cluster where the
controller runs on. The same deployment can also be deployed to "prod" with a slightly different `KluctlDeployment`:

```yaml
apiVersion: flux.kluctl.io/v1alpha1
kind: KluctlDeployment
metadata:
  name: microservices-demo-prod
spec:
  interval: 5m
  path: "./microservices-demo/3-templating-and-multi-env/"
  sourceRef:
    kind: GitRepository
    name: microservices-demo
  timeout: 2m
  target: prod
  prune: true
  # kluctl targets specify the expected context name, which does not necessarily match the context name
  # found while it is deployed via the controller. This means we must pass a kubeconfig to kluctl that has the
  # context renamed to the one that it expects.
  renameContexts:
    - oldContext: default
      newContext: kind-kind
```

## Multiple Clusters

To make things easy for now, the above examples stick with a single cluster. The microservices demo project is
deploying all targets to different namespaces already, so this is enough to showcase the Flux support. To make it work
with multiple clusters, simply install the controller on another cluster and create the appropriate `KluctlDeployment`
objects per cluster.

As an alternative, you can have a central Flux (+flux-kluctl-controller) installation that deploys to multiple clusters.
This can be achieved with the help of the [spec.kubeconfig and spec.serviceAccountName]({{< ref "docs/flux/kluctldeployment#kubeconfigs-and-rbac" >}})
field of the `KluctlDeployment` object.

Also, as the examples stem from the [microservices demo]({{< ref "docs/guides/tutorials/microservices-demo" >}}), they
use the `kind-kind` context names. In a more realistic setup, you would use the real cluster/context names here. This
also assumes that all developers will then use the same context names to refer to the same clusters. If this is honored,
you gain 100% compatibility between the GitOps based deployments and CLI based deployments.

## What's next?

The flux-kluctl-controller and Kluctl itself already support dynamic feature/review environments, meaning that the
controller can create new targets/deployments dynamically. The next blog article will go into the details of these
feature/review environments.
