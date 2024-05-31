---
title: "Deploying via GitOps"
linkTitle: "Deploying via GitOps"
weight: 40
---

This recipe will try to give best practices on how to leverage the kluctl controller to implement Kluctl GitOps.
Before exploring Kluctl GitOps, it is suggested to first learn how Kluctl works without GitOps being involved.

You should also try to understand [how to deploy to multiple targets/environments]({{% ref "docs/recipes/multi-env" %}})
first to get a basic understanding of how the same deployment project can be deployed multiple times.

The source shown in this recipe can also be found on GitHub in the [kluctl-examples repository](https://github.com/kluctl/kluctl-examples/tree/main/gitops-deployment)

## GitOps is optional

Kluctl follows a command-line-first approach, which means that all features implemented into Kluctl will always be
added in a way so that you can keep using the CLI. This means, that Kluctl does not depend on the controller to
implement all its features.

Letting the controller take over is optional and can even be done in a way so that you can
mix CLI based (push-based GitOps) approaches and controller based approaches (pull-based GitOps).

## GitOps is just an interface

Kluctl considers GitOps as just another interface for your deployments. This means that everything that can be
performed and configured via the CLI can also be configured through the Kluctl CRDs
([`KluctlDeployment`]({{% ref "docs/gitops/spec/v1beta1/kluctldeployment" %}})).

Consider a deployment project that you usually deploy via these commands:

```shell
$ git clone https://github.com/kluctl/kluctl-examples.git
$ cd simple
$ kluctl deploy -t simple -a environment=test
```

The above lines perform a deployment in the "push" style, meaning that you (or your CI) pushes the deployment to the
target cluster. That same deployment project can also be deployed in "pull" style, which involves the
[kluctl-controller]({{% ref "docs/gitops" %}}) running on the target cluster that "pulls" the deployment into the cluster.

If you have the controller already [installed]({{% ref "docs/gitops/installation" %}}), you can apply the following
[`KluctlDeployment`]({{% ref "docs/gitops/spec/v1beta1/kluctldeployment" %}}) to your target cluster:

```yaml
# file example-deployment.yaml
apiVersion: gitops.kluctl.io/v1beta1
kind: KluctlDeployment
metadata:
  name: example-deployment
  namespace: kluctl-system
spec:
  interval: 5m
  source:
    git:
      url: https://github.com/kluctl/kluctl-examples.git
      path: simple
  target: simple
  args:
    environment: test
  context: default
```

The above manifest can be applied via plain `kubectl apply -f example-deployment.yaml` or via a Kluctl deployment
project. [Later sections](#managing-gitops-deployments) will go into more detail about some possible options.

After the `KluctlDeployment` got applied, the controller will periodically (5m interval) clone the repository and check
if the result of the rendering process differs since the last deployment. If it differs, the controller will deploy
the deployment project with the given options (which are equal to options of the CLI example from above).

## The reconciliation loop

After a `KluctlDeployment` is applied to the cluster, the kluctl-controller will immediately pick up that deployment
and start to periodically reconcile the deployment. Reconciliation basically performs the following steps:

1. Clone the referenced source (don't worry, this fast due to internal caching)
2. Render the deployment with all the provided options (target, args, ...)
3. Check if the rendered result has changed since the last performed deployment
4. If it has not changed, sleep for the duration specified via `interval` and then repeat the reconciliation loop
5. If it has changed, perform a deployment and record the deployment result in the cluster (this can then be used via the Kluctl Webui)
6. Sleep for the duration specified via `interval` and then repeat the reconciliation loop

If you already know GitOps from other solutions (e.g. Flux), you might notice that Kluctl does not deploy
on every reconciliation iteration but instead only when the source changes. This deviation from other GitOps solutions
is intended as it enabled more flexible intervention and processes (e.g, mixing GitOps with push-based processes).

To mitigate drift between the source and the cluster state, drift detection is performed on every reconciliation
iteration. If necessary, the drift can be viewed and fixed via the [Kluctl Webui]({{% ref "docs/webui" %}}) or via
the [GitOps commands](#gitops-commands).

You can also override this behavior to match the behavior of other GitOps solutions by using
[deployInterval]({{% ref "docs/gitops/spec/v1beta1/kluctldeployment/#deployinterval" %}}), which will cause the
reconciliation loop to periodically perform a deployment even if the source does not change.

## Starting with Kluctl GitOps

To start using Kluctl GitOps, [install]({{% ref "docs/gitops/installation" %}}) it into your cluster first.

Optionally, if you want to use the [Kluctl Webui]({{% ref "docs/webui" %}}) to monitor and control your GitOps
deployments, either run it [locally]({{% ref "docs/webui/running-locally" %}}) or [install]({{% ref "docs/webui/installation" %}})
it into the cluster.

## Managing GitOps deployments

`KluctlDeployment` resources need to be applied and managed the same way as any other Kubernetes resource. You might
easily end up managing dozens or even hundreds of `KluctlDeployment`s per cluster. The recommended way to do this is
to introduce a dedicated GitOps deployment project which is only responsible for the management of other deployments.

Other options exist as well, it's for example also possible to include the `KluctlDeployment` resource into the
deployment itself, so when you perform the initial deployment, you will automatically let GitOps take over. The following
sections will go into more detail.

## Dedicated GitOps deployment project

In this setup, you'll have one dedicated directory ([a simple deployment item]({{% ref "docs/kluctl/deployments/deployment-yml/#simple-deployments" %}}))
for each cluster. These deployment items will contain one or more `KluctlDeployment` resources.

The deployment works by using a simple templated entry in `deployments` which uses the argument `cluster_name` so that
a different directory is loaded for each cluster.

An `clusters/all` deployment item is loaded as well for each cluster. The `clusters/all` deployment item is meant to
add common deployments that are needed on all clusters. One of these deployments is the GitOps deployment itself, so
that it is also managed via GitOps.

The `namespaces` deployment item is used to create the `kluctl-gitops` namespace which we then use to deploy the
`KluctlDeployment` resources into. It's generally best practice to use a dedicated namespace for GitOps. 

### Project structure

Consider the following project structure:

```
gitops-deployment
├── namespaces
│   └── kluctl-gitops.yaml
├── clusters/
│   ├── test.example.com/
│   │   ├── app1.yaml
│   │   └── app2.yaml
│   ├── prod.example.com/
│   │   ├── app1.yaml
│   │   └── app2.yaml
│   ├── all/
│   │   └── gitops.yaml
│   └── deployment.yaml
├── .kluctl.yaml
└── deployment.yaml
```

And the following YAML files and manifests:

```yaml
# .kluctl.yaml
args:
  # This allows us to deploy the GitOps deployment to different clusters. It is used to include dedicated deployment
  # items for the selected cluster.
  - name: cluster_name

targets:
  - name: gitops

# Without a discriminator, pruning won't work. Make sure the rendered result is unique on the target cluster
discriminator: gitops-{{ args.cluster_name | slugify }}
```

```yaml
# deployment.yaml
deployments:
  - path: namespaces
  - barrier: true
  - include: clusters
```

```yaml
# clusters/deployment.yaml
deployments:
  # Include things that are required on all clusters (e.g., the KluctlDeployment for the GitOps deployment itself)
  - path: all
  # We use simple templating to change a dedicated deployment item per cluster
  - path: {{ args.cluster_name }}
```

```yaml
# namespaces/kluctl-gitops.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: kluctl-gitops
```

```yaml
# clusters/test.example.com/app1.yaml
# and clusters/prod.example.com/app1.yaml
# but with adjusted specs (e.g., environment names differ)
apiVersion: gitops.kluctl.io/v1beta1
kind: KluctlDeployment
metadata:
  name: app1
  namespace: kluctl-gitops
spec:
  interval: 5m
  source:
    git:
      url: https://github.com/kluctl/kluctl-examples.git
      path: simple
  target: simple
  args:
    environment: test
  context: default
  # Let it automatically clean up orphan resources and delete all resources when the KluctlDeployment itself gets
  # deleted. You might consider setting these to false for prod and instead do manual pruning and deletion when the
  # need arises.
  prune: true
  delete: true
```

```yaml
# clusters/test.example.com/app2.yaml
# and clusters/prod.example.com/app2.yaml
# but with adjusted specs (e.g., environment names differ)
apiVersion: gitops.kluctl.io/v1beta1
kind: KluctlDeployment
metadata:
  name: app2
  namespace: kluctl-gitops
spec:
  interval: 5m
  source:
    git:
      url: https://github.com/kluctl/kluctl-examples.git
      path: simple-helm
  target: simple-helm
  args:
    environment: test
  context: default
  # Let it automatically clean up orphan resources and delete all resources when the KluctlDeployment itself gets
  # deleted. You might consider setting these to false for prod and instead do manual pruning and deletion when the
  # need arises.
  prune: true
  delete: true
```

```yaml
# clusters/all/gitops.yaml
apiVersion: gitops.kluctl.io/v1beta1
kind: KluctlDeployment
metadata:
  name: gitops
  namespace: kluctl-gitops
spec:
  interval: 5m
  source:
    git:
      url: https://github.com/kluctl/kluctl-examples.git
      path: gitops-deployment # You could also use a dedicated repository without a sub-directory
  target: gitops
  args:
    # this passes the cluster_name initially passed via `kluctl deploy -a cluster_name=xxx.example.com` into the KluctlDeployment
    cluster_name: {{ args.cluster_name }}
  context: default
  # let it automatically clean up orphan KluctlDeployment resources
  prune: true
  delete: true
```

### Managing the GitOps deployment project

Please ensure that you have committed and pushed all required files before you bootstrap the GitOps deployment.
Otherwise, you'll end up deploying different states from your local version while the controller will apply the
Git version.

To bootstrap the GitOps deployment project, simply perform a `kluctl deploy`:

```shell
$ cd gitops-deployment
$ kluctl deploy -a cluster_name=test.example.com
```

This will deploy the GitOps deployment to the current context cluster. After this deployment, the `kluctl-controller` will
immediately start reconciling all deployed `KluctlDeployment` resources, including the one for the GitOps deployment
itself.

This means, to change any of the deployments, perform the changes in Git via your already established processes
(e.g., pull-requests or direct pushes to the main branch).

## GitOps commands

Each individual `KluctlDeployment` can be controlled and inspected via the
[Kluctl CLI]({{% ref "docs/kluctl/commands" %}}) (check the `kluctl gitops xxx` sub-commands). Each command takes the
`KluctlDeployment` name and its namespace as arguments.

In addition, if `--name` and `--namespace` are omitted, the CLI will try to auto-detect the `KluctlDeployment` if your
current directory is inside a Kluctl deployment project. It does so by using the URL of the Git `origin` remote and the
subdirectory inside the Git repository to find one or more `KluctlDeployment` that refers to this project.

### Suspend and resume

The CLI can suspend and resume individual `KluctlDeployment`s. This is useful if you need to perform work that would
otherwise be hard to perform with constant reconciliation being active. This includes refactorings, migrations and other
more complex tasks. While suspended, manual reconciliation via the [CLI](#manual-reconciliation) and the Webui is still
possible.

To suspend the `app1` deployment, run the following CLI command:

```shell
$ kluctl gitops suspend --namespace kluctl-gitops --name app1
```

While suspended, you can perform whatever actions you need without the `kluctl-controller` intervening. Then, to resume
the deployment, run:

```shell
$ kluctl gitops resume --namespace kluctl-gitops --name app1
```

### Manual reconciliation

You can trigger different manual requests via the CLI. Please note that these requests are executed by the
controller even though the usage of the CLI feels like things are executed locally.

Every manual request command is able to override many of the spec fields found in the `KluctlDeployment`. The CLI
tries its best to mimic the interface already found in the non-GitOps based commands (e.g. `kluctl deploy`).

As an example, with `kluctl gitops deploy --namespace=xxx --name=yyy` you can pass deployment arguments
via `-a my_arg=my_value` the same way as you can already do with `kluctl deploy`.

{{% alert context="warning" %}}
WARNING: Currently, the Kluctl GitOps sub-commands do not ask for confirmation before actually performing the requested action.
Always run a `kluctl gitops diff ...` before running any potentially disruptive commands. This behavior might change
in the future.
{{% /alert %}}

The CLI will also try to detect if the Git repository in which you're currently in is related to the Git repository used
in the referenced `KluctlDeployment`. In that case, the CLI will upload the local source code to the controller for a
one-time override. This means, that the `kluctl-controller` will actually work with your local version of the project.
This is mostly useful when you want to verify that changes are valid before actually pushing/merging your changes.

The following invocation will request a single reconciliation iteration. This means, it will do the same as described
in [The reconciliation loop](#the-reconciliation-loop).

```shell
$ kluctl gitops reconcile --namespace kluctl-gitops --name app1
```

The following invocation will perform a diff and print the result. This is especially useful if your local version of
the source code contains modifications which you'd like to verify.

```shell
$ kluctl gitops diff --namespace kluctl-gitops --name app1
```

The following invocation will cause a manual prune (delete orphan objects).

```shell
$ kluctl gitops prune --namespace kluctl-gitops --name app1
```

## Viewing controller logs

The following CLI command can be used to view controller logs related to a given `KluctlDeployment`:

```shell
$ kluctl gitops logs --namespace kluctl-gitops --name app1 -f
```

## Using the Webui

In addition to the Kluctl GitOps commands, the [Kluctl Webui]({{% ref "docs/webui" %}}) can be used to monitor and
control the `KluctlDeployment`s.

The Webui is still very experimental, meaning that many features are still missing. But generally, performing manual
requests, viewing state, diffs and logs should already work good enough as of now.

## Mixing

Kluctl allows you to mix pull-based GitOps with push-based CLI workflows. You can use GitOps for some
targets/environments (e.g. prod) and revert to using push-based CLI workflows in other targets/environments (e.g. dev
environments). This is useful if you want the security and stability of GitOps on prod while still having the flexibility
and speed of development on non-prod environments.

You can also use GitOps for a target/environment to perform the actuall deployments while using `kluctl diff` in the
push fashion to test/verify changes before actually pushing/merging the main branch.
