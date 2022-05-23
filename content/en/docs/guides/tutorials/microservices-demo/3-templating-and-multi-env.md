---
title: "3. Templating and multi-env deployments"
linkTitle: "3. Templating and multi-env deployments"
weight: 3
---

## Introduction
The [second tutorial]({{< ref "docs/guides/tutorials/microservices-demo/2-helm-integration" >}}) in this series demonstrated how
to integrate Helm into your deployment project and how to keep things structured.

The project is however still not flexible enough to be deployed multiple times and/or in different flavors. As an example,
it doesn't make much sense to deploy redis with replication on a local cluster, as there can't be any high availability
with single node. Also, the resource requests currently used are quite demanding for a single node cluster.

## How to start
This tutorial is based on the results of the second tutorial. As an alternative, you can take the `2-helm-integration`
example project found [here](https://github.com/kluctl/kluctl-examples/blob/main/microservices-demo/2-helm-integration)
and use it as the base to be able to continue with this tutorial.

This time, you should start with a fresh kind cluster. If you are sure that you won't loose any critical data by deleting
the existing cluster, simply run:
```sh
$ kind delete cluster
$ kind create cluster
```

If you're unsure or if you want to re-use the existing cluster for some reason, you can also simply delete the old deployment:
```yaml
$ kluctl delete -t local
  INFO[0000] Rendering templates and Helm charts
  INFO[0000] Building kustomize objects
  INFO[0000] Getting remote objects by commonLabels
The following objects will be deleted:
  default/Service/emailservice
  ...
  default/ConfigMap/cart-redis-scripts
  Do you really want to delete 29 objects? (y/N) y

Deleted objects:
  default/ConfigMap/cart-redis-scripts
  ...
  default/StatefulSet/cart-redis-node
```

The reason to start with a fresh deployment is that we will later switch to different namespaces and stop using the
`default` namespace.

## Targets
If we want to allow the deployment to be deployed multiple times, we first need multiple targets in our project. Let's
add 2 [targets]({{< ref "docs/reference/kluctl-project/targets" >}}) called `test` and `prod`. To do so, modify the
content of `.kluctl.yml` to contain:

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

You might notice that all targets point to the kind cluster at the moment. This is of course not how you would do it
in a real project as you'd probably have at least one real production-ready cluster to target your deployments against.

We've also introduced [`args`]({{< ref "docs/reference/kluctl-project/targets#args" >}}) for each target, with each target
having an `env_type` argument configured. This argument will later be used to change details of the deployment, depending
on the value of it. For example, setting it to `local` might change the redis deployment into a single-node/standalone
deployment.

## Dynamic namespaces
One of the most obvious and also useful application of templates is making namespaces dynamic, depending on the target
that you want to deploy. This allows to deploy the same set of deployment/manifests multiple times, even to the same
cluster.

There are a few [predefined variables]({{< ref "docs/reference/templating/predefined-variables" >}}) which are always available
in all deployments. One of these variables is the `target` dictionary which is a copy of the currently processed target.
This means, we can use `{{ target.name }}` to insert the current target name through templating.

There are multiple ways to change the namespaces of involved resources. The most naive way is to go directly into the
manifests and add the `metadata.namespace` field. For example, you could edit `services/adservice/deployment.yml` this
way:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: adservice
  namespace: ms-demo-{{ target.name }}
...
```

This can however easily lead to resources being missed or resources where you are not in control, e.g. rendered
Helm Charts. Another way to set the namespace on multiple resources is by using the
[`namespace` property](https://kubectl.docs.kubernetes.io/guides/config_management/namespaces_names/) of kustomize.
For example, instead of changing the `adservice` deployment directly, you could modify the content of
`services/adservice/kustomization.yml` to:

```yaml
resources:
  - deployment.yml
  - service.yml

namespace: ms-demo-{{ target.name }}
```

This is better than the naive solution, but still limited in a comparable (but not as bad) way. The most powerful and
preferred solution is use [`overrideNamespace`]({{< ref "docs/reference/deployments/deployment-yml#overridenamespace" >}})
in the root `deployment.yml`:

```yaml
...
overrideNamespace: ms-demo-{{ target.name }}
...
```

As an alternative, you could also use `overrideNamespace` separately in `third-party/deployment.yml` and
`services/deployment.yml`. In this case, you're also free to use different prefixes for the namespaces, as long as you
include `{{ target.name }}` in them.

{{< alert >}}
Please note that `overrideNamespace` only takes effect on a kustomize deployment if it does NOT specify `namespace`.
If you followed the `kustomization.yml` example from above, make sure to undo the changes to `kustomization.yml`.
{{< /alert >}}

## Helm Charts and namespaces
The previously described way of making namespaces dynamic in all resources works well for most cases. There are however
situations where this is not enough, mostly when the name of the namespace is used in other places than `metadata.namespace`.

Helm Charts very often do this internally, which makes it necessary to also include the dynamic namespace into the
`helm-chart.yml`'s `namespace` property. You will have to do this for the redis chart as well, so let's modify
`third-party/redis/helm-chart.yml` to:

```yaml
helmChart:
  repo: https://charts.bitnami.com/bitnami
  chartName: redis
  chartVersion: 16.8.2
  releaseName: cart
  namespace: ms-demo-{{ target.name }}
  output: deploy.yml
```

Without this change, redis is going to be deployed successfully but will then fail to start due to wrong internal
references to the default namespace.

## Making commonLabels unique per target
[`commonLabels`]({{< ref "docs/reference/deployments/deployment-yml#commonlabels" >}}) in your root `deployment.yml` has
a very special meaning which is important to understand and work with. The combination of all `commonLabels` MUST be unique
between all supported targets on a cluster, including the ones that don't exist yet and are from other kluctl projects.

This is because kluctl uses these to identify resources belonging to the currently processed deployment/target,
which becomes especially important when deleting or pruning.

To fulfill this requirement, change the root `deployment.yml` to:
```yaml
...
commonLabels:
  examples.kluctl.io/deployment-project: "microservices-demo"
  examples.kluctl.io/deployment-target: "{{ target.name }}"
...
```

`examples.kluctl.io/deployment-project` ensures that we don't get in conflict with any other kluctl project that might
get deployed to the same cluster. `examples.kluctl.io/deployment-target` ensures that the same deployment can be deployed
once per target. The names of the labels are arbitrary, and you can choose whatever you like.

## Creating necessary namespaces
If you'd try to deploy the current state of the project, you'd notice that it will result in many errors where kluctl
says that the desired namespace is not found. This is because kluctl does not create namespaces on its own. It also
does not do this for Helm Charts, even if `helm install` for the same charts would do this. In kluctl you have to
create namespaces by yourself, which ensures that you have full control over them.

This implies that we must create the necessary namespace resource by ourselves. Let's put it into its own kustomize deployment below
the root directory. First, create the `namespaces` directory and place a simple `kustomization.yml` into it:

```yaml
resources:
  - namespace.yml
```

In the same directory, create the manifest `namespace.yml`:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: ms-demo-{{ target.name }}
```

Now add the new kustomize deployment to the root `deployment.yml`:

```yaml
deployments:
  - path: namespaces
  - include: third-party
  - include: services
...
```

## Deploying multiple targets
You're now able to deploy the current deployment multiple times to the same kind cluster. Simply run:
```sh
$ kluctl deploy -t local
$ kluctl deploy -t prod
```

After this, you'll have two namespaces with the same set of microservices and two instances of redis (both replicated with
3 replicas) deployed.

## All changes together
For a complete overview of the necessary changes to get to this point, look into [this commit](https://github.com/kluctl/kluctl-examples/commit/511fc0e06790152bfcedf9caa6b402029fea0ffb).

## Make the local target more lightweight
Having the microservices demo deployed twice might easily lead to you local cluster being completely overloaded. The
solution would obviously be to not deploy the prod target to your local cluster and instead use a real cluster.

However, for the sake of this tutorial, we'll instead try to introduce a few differences between targets so that they
fit better onto the local cluster.

To do so, let's introduce [variables files]({{< ref "docs/reference/templating/variable-sources" >}}) that contain
different sets of configuration for different environment types. These variables files are simply yaml files with
arbitrary content, which is then available in future templating contexts.

First, create the sub-directory `vars` in the root project directory. The name of this directory is arbitrary and up to
you, it must however match what is later used in the `deployment.yml`.

Inside this directory, create the file `local.yml` with the following content:
```yaml
redis:
  architecture: standalone
  # the standalone architecture exposes redis via a different service then the replication architecture (which uses sentinel)
  svcName: cart-redis-master
```

And the file `real.yml` with the following content:
```yaml
redis:
  architecture: replication
  # the standalone architecture exposes redis via a different service then the replication architecture (which uses sentinel)
  svcName: cart-redis
```

To load these variables files into the templating context, modify the root `deployment.yml` and add the following to the top:
```yaml
vars:
  - file: ./vars/{{ args.env_type }}.yml
...
```

As you can see, we can even use templating inside the `deployment.yml`. Generally, templating can be used everywhere,
with a few limitations outlined in the documentation.

The above changes will now load a different variables file, depending on which `env_type` was specified in the currently
processed target. This allows us to customize all kinds of configurations via templating. You're completely free in how
you use this feature, including loading multiple variables files where each one can use the variables loaded by the
previous variables file.

To use the newly introduces variables, first modify the content of `third-party/redis/helm-values.yml` to:
```yaml
architecture: {{ redis.architecture }}

auth:
  enabled: false

{% if redis.architecture == "replication" %}
sentinel:
  enabled: true
  quorum: 2

replica:
  replicaCount: 3
  persistence:
    enabled: true
{% endif %}

master:
  persistence:
    enabled: true
```

The templating engine used by kluctl is currently [Jinja2](https://jinja.palletsprojects.com). We suggest reading through
the documentation of Jinja2 to understand what is possible. In the example above, we use simple
[variable expressions](https://jinja.palletsprojects.com/en/3.1.x/templates/#variables) and
[if/else](https://jinja.palletsprojects.com/en/3.1.x/templates/#if) statements.

We will also have to replace the occurrence of `cart-redis:6379` with `{{ redis.svcName }}:6379` inside
`services/cartservice/deployment.yml`.

For an overview of the above changes, look into [this commit](https://github.com/kluctl/kluctl-examples/commit/9b2f5a516f75493f420f5997a8559f624c5b1b21).

## Deploying the current state
You can now try to deploy the `local` and `test` targets. You'll notice that the `local` deployment will result in quite
a few changes (seen in the diff) and the `test` target not having any changes at all. You might also want to do a prune
for the `local` target to get rid of the old redis deployment.

## Disable a few services for local
Some services are not needed locally or might not even be able to run properly. Let's assume this applies to the services
`loadgenerator` and `emailservice`. We can conditionally remove these from the deployment with simple boolean variables
in `vars/local.yml` and `vars/real.yml` and if/else statements in `services/deployment.yml`.

Add the following variables to `vars/local.yml`:
```yaml
...
services:
  emailservice:
    enabled: false
  loadgenerator:
    enabled: false
```

And the following variables to `vars/real.yml`:
```yaml
...
services:
  emailservice:
    enabled: true
  loadgenerator:
    enabled: true
```

Now change the content of `services/deployment.yml` to:
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

A deployment to `test` should not change anything now. Deploying to `local` however should reveal multiple orphan resources,
which you can then prune.

For an overview of the above changes, look into [this commit](https://github.com/kluctl/kluctl-examples/commit/dd91592f5f7848444d1fe05b4e1295805604e3c2).

## How to continue
After this tutorial, you should have a basic understanding how templating in kluctl works and how a multi-environment
deployment can be implemented.

We however only deployed to a single cluster so far and are unable to properly manage the image versions of our microservices
at the moment. In the next tutorial of this series, we'll learn how to deploy to multiple clusters and split third-party
image management and (self developed) application image management.
