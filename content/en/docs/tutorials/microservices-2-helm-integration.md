---
title: "Microservices Demo 2 - Helm Integration"
linkTitle: "Microservices Demo 2 - Helm Integration"
weight: 2
---

## Introduction
The [first tutorial]({{< ref "docs/tutorials/microservices-1-basic-project-setup" >}}) in this series demonstrated how
to setup a simple kluctl project that is able to deploy the [GCP Microservices Demo](https://github.com/GoogleCloudPlatform/microservices-demo)
to a local kind cluster.

This initial kluctl project was however quite naive and too simple to be any way realistic. For example, the project
structure is too flat and will likely result in chaos when the project grows. Also, the project used self-crafted
manifests while it might have been better to reuse feature rich Helm Charts. We will fix both these issues in this
tutorial.

## How to proceed
This tutorial is based on the results of the first tutorial. As an alternative, you can take the `1-basic-project`
example project found [here](https://github.com/kluctl/kluctl-examples/blob/main/microservices-demo/1-basic-project)
and use it the base to be able to continue with this tutorial.

You can also deploy the base project and then incrementally perform deployment after each step in this tutorial. This way
you will also gain some experience and feeling for to use kluctl.

## A simple refactoring
Let's start with a simple refactoring. Having all deployment items on the root level will easily get unmaintainable.

kluctl allows you to structure your project in all kinds of fashions by leveraging sub-deployments. The
[deployment items]({{< ref "docs/reference/deployments/deployment-yml#deployments" >}}) found in deployment projects
allows specifying [includes]({{< ref "docs/reference/deployments/deployment-yml#includes" >}}) which point to sub-directory
with another `deployment.yml`.

Let's split the deployment into third-party applications (currently only redis) and the project specific microservices.
To do this, create the sub-directories `third-party` and `microservices`. Then move the `redis` directory into `third-party`
and all microservice sub-directories into `microservices`:

```sh
$ mkdir third-party
$ mkdir microservices
$ mv redis third-party/
$ mv adservice cartservice checkoutservice currencyservice emailservice \
    frontend loadgenerator paymentservice \
    productcatalogservice recommendationservice shippingservice microservices/
```

Now change the `deployments` list inside the root `deployment.yml` to:
```yaml
deployments:
  - include: third-party
  - include: services
```

Add a `deployment.yml` with the following content into the `third-party` sub-directory:
```yaml
deployments:
  - path: redis
```

And finally a `deployment.yml` with the following content into the `microservices` sub-directory:
```yaml
deployments:
  - path: adservice
  - path: cartservice
  - path: checkoutservice
  - path: currencyservice
  - path: emailservice
  - path: frontend
  - path: loadgenerator
  - path: paymentservice
  - path: productcatalogservice
  - path: recommendationservice
  - path: shippingservice
```

To get an overview of these changes, look into [this commit](https://github.com/kluctl/kluctl-examples/commit/1388dd025e5471c3d6727fe01e626cee6091e4fb)
inside the example project belonging to this tutorial.

If you deploy the new state of the project, you'll notice that only labels will change. These labels are automatically
added to all resources and represent the [tags]({{< ref "docs/reference/deployments/tags" >}}) of the corresponding
deployment items.

## Some notes on project structure
The refactoring from above is meant as an example that demonstrates how sub-deployments can be used to structure your
project. Such sub-deployments can also include deeper sub-deployments, allowing you to structure your project in any
way and complexity that fits your needs.

## Introducing the first Helm Chart
There are many examples where self-crafting of Kubernetes manifests is not the best solution, simply because there is
already a large ecosystem of pre-created Kubernetes packages in the form of [Helm Charts](https://hub.helm.sh).

The redis deployment found in the microservices demo is a good example for this, especially as many available Helm Charts
offer quite some functionality, for example high availability.

kluctl allows the [integration]({{< ref "docs/reference/deployments/helm" >}}) of Helm Charts, which we will do now to
replace the self-crafted redis deployment with the [Bitname Redis Chart](https://artifacthub.io/packages/helm/bitnami/redis).

First, create the file [`third-party/redis/helm-chart.yml`]({{< ref "docs/reference/deployments/helm#helm-chartyml" >}}) with the following content:
```yaml
helmChart:
  repo: https://charts.bitnami.com/bitnami
  chartName: redis
  chartVersion: 16.8.0
  releaseName: cart
  namespace: default
  output: deploy.yml
```

Most of the above configuration can directly be mapped to Helm invocations (pull, install, ...). The `output`
value has a special meaning and must be reflected inside the `kustomization.yml` resources list. The reason is that
kluctl solves the Helm integration by running [helm template](https://helm.sh/docs/helm/helm_template/) and writing
the result to the file configured via `output`. After this, kluctl expects that kustomize takes over, which requires
that the generated file is references in `kustomization.yml`.

To do so, simply replace the content of `third-party/redis/kustomization.yml` with:
```yaml
resources:
  - deploy.yml
```

We now need some configuration for the redis chart, which is provides via `[`third-party/redis/helm-values.yml`]({{< ref "docs/reference/deployments/helm#helm-valuesyml" >}}):
```yaml
architecture: replication

auth:
  enable: false

sentinel:
  enabled: true
  quorum: 2

replica:
  replicaCount: 3
  persistence:
    enabled: true

master:
  persistence:
    enabled: true
```

The above configuration will configure redis to run in replication mode with sentinel and 3 replicas, giving us some
high availability (at least in theory, as we'd still need a HA Kubernetes cluster and proper affinity configuration).

The Redis Chart will also deploy a `Service` resource, but with a different name as the self-crafted version. This means
we have to fix the service name in `microservices/cartservice/deployment.yml` (look for the environment variable REDIS_ADDR)
to point to `cart-redis:6379` instead of `redis-cart:6379`.

You can now remove the old redis related manifests (`third-party/redis/deployment.yml` and `third-party/redis/service.yml`).

All the above changes can be found in [this commit](https://github.com/kluctl/kluctl-examples/commit/e90dd85a2947402d08172295bb3ac22d27e72123)
from the example project.

## Pulling Helm Charts
We have now added a Helm Chart to our deployment, but to make it deployable it must be pre-pulled first. kluctl 
requires Helm Charts to be pre-pulled for multiple reasons. The most important reasons are performance and reproducibility.
Performance would significantly suffer if Helm Chart would have to be pulled on-demand at deployment time. Also,
Helm Charts have no functionality to ensure that a chart that you pulled yesterday is equivalent to the chart pulled today,
even if the version is unchanged.

To pre-pull the redis Helm Chart, simply call:
```sh
$ kluctl helm-pull
INFO[0000] Pulling for third-party/redis/helm-chart.yml
```

This will pre-pull the chart into the sub-directory `third-party/redis/charts`. This directory is meant to be added
to version control, so that it is always available when deploying.

If you ever change the chart version in `helm-chart.yml`, don't forget to re-run the above command and commit the
resulting changes.

## Deploying the current state
It's time to deploy the current state again:

```sh
$ kluctl deploy -t local
INFO[0000] Rendering templates and Helm charts          
...          

New objects:
  default/ConfigMap/cart-redis-configuration
  default/ConfigMap/cart-redis-health
  default/ConfigMap/cart-redis-scripts
  default/Service/cart-redis
  default/Service/cart-redis-headless
  default/ServiceAccount/cart-redis
  default/StatefulSet/cart-redis-node

Changed objects:
  default/Deployment/cartservice

Diff for object default/Deployment/cartservice
+-------------------------------------------------------+------------------------------+
| Path                                                  | Diff                         |
+-------------------------------------------------------+------------------------------+
| spec.template.spec.containers[0].env.REDIS_ADDR.value | -redis-cart:6379             |
|                                                       | +cart-redis:6379             |
+-------------------------------------------------------+------------------------------+

Orphan objects:
  default/Deployment/redis-cart
  default/Service/redis-cart
```

As you can see, the changes that we did to the kluctl project are reflected in the output of the deploy call, meaning
that we can perfectly see what happened. We can see a few new resources which are all redis related, the change of the
service name and the old redis resources being marked as orphan. Let's get rid of the orphan resources:

```sh
$ kluctl prune -t local
INFO[0000] Rendering templates and Helm charts          
INFO[0000] Building kustomize objects                   
INFO[0000] Getting remote objects by commonLabels       
The following objects will be deleted:
  default/Service/redis-cart
  default/Deployment/redis-cart
Do you really want to delete 2 objects? (y/N) y

Deleted objects:
  default/Service/redis-cart
  default/Deployment/redis-cart
```

You have just performed your first house-keeping, which you'll probably do quite often from now on in your daily
DevOps business.

## More house-keeping
When time passes, new versions of the Helm Charts that you integrated are going to be released. You might have to keep
your deployments up-to-date in such cases. The most naive way is to simply increase the chart version inside `helm-chart.yml`
and then simply re-call `kluctl helm-pull`.

As the number of used charts can easily grow to a number where it becomes hard to keep everything up-to-date, kluctl
offers a command to support you in this:

```sh
$ kluctl helm-update
INFO[0005] Chart third-party/redis/helm-chart.yml has new version 16.8.2 available. Old version is 16.8.0. 
```

As you can see, it will display charts with new versions. You can also use the same command to actually update the
`helm-chart.yml` files and ultimately commit these to git:
```sh
$ kluctl helm-update --upgrade --commit
INFO[0005] Chart third-party/redis/helm-chart.yml has new version 16.8.2 available. Old version is 16.8.0. 
INFO[0005] Pulling for third-party/redis/helm-chart.yml 
INFO[0010] Committing: Updated helm chart third-party/redis from 16.8.0 to 16.8.2
```

## How to continue
After this tutorial, you have hopefully learned how to better structure your projects and how to integrate third-party
Helm Charts into your project, including some basic house-keeping tasks.

The next tutorials in this series will show you how to use this kluctl project as a base to implement a multi-environment
and multi-cluster deployment.
