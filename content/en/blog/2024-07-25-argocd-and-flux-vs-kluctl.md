---
title: "ArgoCD/Flux vs Kluctl"
linkTitle: "ArgoCD/Flux vs Kluctl"
slug: argcd-and-flux-vs-kluctl
date: 2024-07-25
author: Alexander Block (@codablock)
description: |
  A comparison of ArgoCD/Flux vs Kluctl
tags:
- Other
images:
- "/images/blog/2024-03-13-cluster-api-kluctl.jpg"
---

![image](/images/blog/2024-07-25-argocd-and-flux-vs-kluctl.jpg)

Kluctl is very flexible when it comes to deployment strategies. All features implemented by Kluctl can be used via
the [CLI]({{% ref "docs/kluctl/commands/" %}}) or via the [Kluctl Controller]({{% ref "docs/gitops/" %}}).

This makes Kluctl comparable to [ArgoCD](https://argo-cd.readthedocs.io/en/stable/) and [Flux](https://fluxcd.io/), as
these projects also implement the GitOps strategy.

This comparison assumes that you already know Flux and/or ArgoCD to some degree, or at least have heard of them, so it
will not go too deep into comparing these against each other. If you want a deep dive into ArgoCD vs Flux, read
the [Comparison: Flux vs Argo CD](https://earthly.dev/blog/flux-vs-argo-cd/) blog post from Earthly.

This post is meant to be updated over its lifetime when things change in any of the projects. Feel free to use the
comment feature or create an issue or pull request to notify us about things that need to be changed.

## Use of Custom Resources

[Kubernetes Custom Resources](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/) (CRs) are used to extend the Kubernetes API with custom resources and custom behavior. CustomResourceDefinitions define
the new API types while Custom Resources (e.g. a Kustomization, HelmRelease or Application) represent single instances
of these new types. A controller watches for changes of a certain type and acts accordingly, e.g. by applying the desired state to the
cluster to move the actual state closer to the desired state (reconciliation).

### ArgoCD and Flux

Both, ArgoCD and Flux heavily use and depend on CRs and Controllers to implement reconciliation and deployments. To
implement more complex deployments, you usually have to use multiple Custom Resources together. Sometimes you have to even chain them (e.g.
Kustomization -> Kustomization -> HelmRelease). Please note that each Kustomization (the CRs) in this chain also might
imply a `kustomization.yaml` that can also use overlays and components which might lead to even longer chains of CRs. A
HelmRelease could potentially also deploy Kustomizations or other HelmReleases, when dynamic configuration via templating
is desired.

Following and fully understaning these chains can become quite challenging, when the project size and complexity grows.

Another implication of using CRs is that your deployments become 100% dependent on the controllers running
inside Kubernetes, because in-cluster reconciliation is the only way to process the CRs. This means, you can not fully
test or verify your deployments before pushing them to your Git repository. The only way to reliably prevent killing your production
environment is to introduce testing/staging environments, adding even more complexity to your depoyments and processes.

### Kluctl

Kluctl also uses CRs ([KluctlDeployment]({{% ref "docs/gitops/spec/v1beta1/kluctldeployment/" %}}), but only as a
bridge between Kluctl deployment projects and GitOps. The actual project structure is solely defined via the Kluctl
project and deployment YAML files found inside your Git repository. The KluctlDeployment CR does not add anything
special to the deployment project itself, meaning that you never get fully dependent on in-cluster reconciliation.

This means, to deploy your project, you can always revert back to using
the [Kluctl CLI]({{% ref "docs/kluctl/commands/" %}}) even if you leverage GitOps as your main deployment strategy.
This might sound counter intiutive at first when talking about GitOps, but there are actually very good reasons and use
cases why you might consider mixing GitOps with other strategies. Please read the next chapter for more details on this.

## Pull vs. Push

GitOps can be implemented in two different strategies. One is the push strategy, which is actually what has been done
reliably for years in classical continuous delivery, e.g. via Jenkins, Gitlab pipelines or GitHub workflows. The way it
is usually implemented however has some disadvantages, for example the possibility of a growing drift between the state in
Git and the state in the cluster.

The other is the pull strategy, implemented via controllers running inside your target cluster. These controllers
usually pull the Git repository in some defined interval and then (re-)apply all resources continuously. This
reliably ensures that drift is being fixed whenever it occurs.

### ArgoCD and Flux

Both ways implement GitOps as a pull strategy. ArgoCD also natively supports manual synchronization, meaning that you can
disable periodic reconciliation and instead rely on manual syncs via the UI.

Push based GitOps is not possible in ArgoCD or Flux, this is however not considered a downside but a strict design
choice. At the same time, it would be very hard to implement due the strict reliance on Custom Resources for core
features.

### Kluctl

Kluctl allows you to choose between the push and the pull-based strategy. It even allows you to switch back and forth or
mix these in the same project. Please read the [chapter](#Use-of-Custom-Resources) about Custom Resources to understand
why this is possible in Kluctl.

Push-based GitOps is implemented via the [Kluctl CLI]({{% ref "docs/kluctl/commands/" %}}), which you can run from
your local machine or from a continuous delivery pipeline. Pull based GitOps is implemented via
the [Kluctl Controller]({{% ref "docs/gitops/" %}}), which takes
a [KluctlDeployment]({{% ref "docs/gitops/spec/v1beta1/kluctldeployment/" %}}) as input and then performs periodic
reconciliation. In the end, both strategies end up using your Git source to perform exactly the same deployment actions.

There are many use-cases where mixing is useful.

One simple example: Running [diffs]({{% ref "docs/kluctl/commands/diff" %}}) against production, because a diff is
actually implemented as a [server-side](https://kubernetes.io/docs/reference/using-api/server-side-apply/) dry-run apply
with the diff happening on current/real state vs. simulated/dry-run state. This means, you can locally implement a
hot-fix for your production system and verify correctness of the fix by running a diff against production.

Another example is using pull-based GitOps for production and the push-based CLI for development/test environments,
allowing you to perform very fast cycles of modify and test iterations, without the need to commit and push your changes
just to find out if your change applies successfully.

## Bootstrapping/Installation

To start with GitOps, some components must usually be installed into the cluster and optionally on your local machine.

### ArgoCD

Multiple options exist for both solutions. For ArgoCD, you can apply
some [static manifests](https://argo-cd.readthedocs.io/en/stable/getting_started/) and install the `argocd` CLI. After
that, you can either use the UI or the CLI to
add [Applications](https://argo-cd.readthedocs.io/en/stable/user-guide/application-specification/) or even additional
clusters. It's common to have one ArgoCD instance manage multiple clusters.
The [App of Apps](https://argo-cd.readthedocs.io/en/stable/operator-manual/cluster-bootstrapping/) pattern allows you to
hand over most operations to Git, reducing the amount configuration required in the UI or via the CLI.

### Flux

Flux requires a process called [bootstrapping](https://fluxcd.io/flux/get-started/#install-flux-onto-your-cluster) to
pre-create all necessary Kubernetes manifests in your Git repository and perform an initial apply, resulting in the Flux
controllers to start up and then take over reconciliation of that same repository. In Flux, it's common to have one set
of controller on every cluster, while Git repositories might be shared (using different subdirectories).

### Kluctl

Kluctl uses a Kluctl deployment project to install the controller into your target cluster. This project is embedded
into the CLI and can be deployed
via [kluctl controller install]({{% ref "docs/gitops/installation/#using-the-install-sub-command" %}}). As an
alternative, you can set up your own Kluctl deployment project that
just [git includes]({{% ref "docs/gitops/installation/#using-a-kluctl-deployment" %}}) the actual controller
deployment. It's common to have such a bootstrap deployment to setup the controller and all the other things a cluster
requires to function, for example your CNI, cert-manager, ingress/gateway controllers, cloud specific controllers, and
so on.

This bootstrap deployment is then simply deployed via `kluctl deploy --context my-cluster`.

## Reconciliationn and Drift

Kubernetes controllers typically implement a reconciliation loop that reconciles the actual/current state towards the
desired state (defined by a CR). GitOps controllers do the same, with the difference that the CR actually references a
source repository (usually Git) which then contains the desired state (instead of putting it directly into the CR).

### ArgoCD

ArgoCD supports periodic synchronization and manual synchronization. This means, that the reconciliation loop will do
different things depending on how the Application CR is configured. This allows you to adapt different strategies, for
example to perform automatic sync for test/staging and only allows manual sync on prod, or vice versa. Sync windows allow
you to further customize the behavior. Drift detection is performed in all cases and drift is properly shown in the UI
even if syncs are not performed.

ArgoCD seems to try its best to cater many use cases and demands of users with this flexibility.

### Flux

Flux only supports periodic synchronization natively. Anything more advanced/customized requires you to somehow control
it from the outside world, e.g. via CronJob resources that suspend/resume reconciliation at defined times. Flux only
performs drift detection to determine what needs to be re-applied.

Flux tries to adhere to strict GitOps principles and does not accept drift as something that it allows to happen in
daily business.

### Kluctl

By default, Kluctl does only re-apply resources when the involved source code changes. This somewhat mimics what
classical push based continuous delivery is doing. In practice this means that if you change something on the cluster (
e.g. via kubectl or k9s), the change/drift stays intact until something gets changed in the source (usually via a
commit + push or merge). This behaviour can be changed in the KluctlDeployment CR to also perform unconditional periodic
re-apply.

Drift detection is always happening, even if no re-apply happened. The drift is properly reflected in the
KluctlDeployment status and can also be seen in the UI.

The Kluctl project considers drift as something that always might happen in daily business and can even be acceptable or
desired. e.g. when fixing a production issues and relying on GitOps to not apply its state (and thus revert the fix)
before the fix gets through established processes and workflows (review, CI/CD, ...). At the same time, it performs
proper drift detection at all times to allow early detection of undesired drift.

## Pruning

As your deployments change over time, things will get added, removed and renamed. This causes leftovers on the cluster
that must be cleaned up or pruned afterward. If this garbage were left on the cluster, you would be guaranteed to
get into trouble long-term.

### ArgoCD and Flux

Both projects follow the same strategy. They perform some form of bookkeeping (e.g. by storing lists in the CRs
status sub-resource) to remember which resources were applied in the past. This allows them to figure out which
resources got removed from the source code and thus need to be pruned.

Both projects also allow you to fully or partially (on resource level via annotations) disable pruning.

### Kluctl

Kluctl does not perform bookkeeping but instead marks all applied resources
with [discriminator labels]({{% ref "docs/kluctl/kluctl-project/targets/#discriminator" %}}). These allow Kluctl to
efficiently query the cluster for all the resources previously deployed from the given Kluctl deployment project. The
query result is then used to determine what got orphaned and thus needs to be pruned. This type of orphan detection has
the advantage that it also works without the CRs and the Kluctl Controller and thus can also be leveraged by the CLI.

Pruning can be disabled on KluctlDeployment level and on resource level (also via annotations). Orphan resource
detection will however always be performed and reported as part of the drift detection (or manual diff invocations via
the CLI).

## Helm

[Helm](https://helm.sh/) is the de-facto standard package manager in Kubernetes. It uses go-templates to implement
configurability and Helm Repositories or OCI registries to distribute Helm Charts.

Helm Charts are originally installed via the Helm CLI, which maintains the resulting Helm Releases lifecycle via
in-cluster release secrets.

The way Helm Release's lifecycle was initially meant to be managed is not very friendly to how GitOps is meant to
function, which adds a few challenges when combined. This resulted in multiple approaches in the different GitOps
solutions, all with their own drawbacks and issues.

### ArgoCD

ArgoCD internally uses `helm template` to render out all manifests. These manifests are then applied to the cluster the
same way as all other static manifests are applied in ArgoCD. Helm Hooks are simulated/re-implemented in ArgoCD to be as
compatible as possible.

This approach leaves full lifecycle management to ArgoCD and does not use any Helm native feature for this. Most
prominent effect of this is that `helm list` does not list Helm Releases managed by ArgoCD. The advantage on the other
side is that drift detection works out of the box the way you would expect.

Helm values files can be pulled from different source repositories. This means, you could use a publicly released Helm
Chart (e.g. cert-manager) from the public repository and at the same time provide your own Helm values from your private
repository.

Helm values files are completely static in ArgoCD and can not be further configured. ArgoCD currently does not
support [post-renderering](https://github.com/argoproj/argo-cd/issues/3698), which means you can not patch/fix upstream
Helm Charts without forking them.

### Flux

Flux has a dedicated [helm-controller](https://fluxcd.io/flux/components/helm/) to implement Helm support.
A [HelmRelease](https://fluxcd.io/flux/components/helm/helmreleases/) references a `HelmRepository` which in turn
contains the URL to a Helm Repository or OCI registry. The controller then uses native Helm features to install, upgrade
and uninstall Helm Releases.

The advantage is that `helm list` and external tools like [helm-dashboard](https://github.com/komodorio/helm-dashboard)
keep working, because these rely on the Helm Release information stored in the Helm Release secrets.

The disadvantage is that drift detection becomes a lot harder to implement. For a long time, drift detection was
completely missing in Flux. It is implemented now, but required a lot of effort to get to this state.

Helm values can be passed via ConfigMaps and Secrets which are usually also managed by the same Flux deployment,
allowing you to perform some configuration via Flux's substitution feature. Post-rendering is fully supported and allows
you to patch/fix upstream Helm Charts without forking them.

### Kluctl

Kluctl follows a similar approach as ArgoCD. It renders the Helm Charts and then applies the rendered manifests the same
way as any other manifest, so that all Kluctl features (e.g. drift detection and pruning) are working as expected. As a
step in-between, the manifests are passed through Kustomize to allow patching of upstream Charts.

Helm values are passed via a `helm-values.yaml` inside your deployment project. Kluctl templating is applied to this
file the same way as it is applied to all other files in your project, allowing you to have very flexible configuration.

## Kustomize

[Kustomize](https://kustomize.io/) is an alternative to Helm with a completely different approach. It markets itself as
template-free and instead relies on overlays and patches to implement configurability.

### ArgoCD

ArgoCD natively supports Kustomize. You can either point your Application to a Git repository containing a self-contained
Kustomization project, inline some Kustomize directives into the Application or even mix both approaches.

### Flux

Flux also natively supports Kustomize via the
dedicated [kustomize controller](https://fluxcd.io/flux/components/kustomize/).

Flux also
supports [post-build variable substitions](https://fluxcd.io/flux/components/kustomize/kustomizations/#post-build-variable-substitution)
on top of the native Kustomize features set. With this, some form of lightweight templating is possible even though
Kustomize is meant to be template-free. If this feature is used however, the Kustomizations incompatible to native
Kustomize, making it harder to test/verify changes without commiting and pushing.

### Kluctl

Kluctl uses Kustomize as the low-level building block
for [deployment items]({{% ref "docs/kluctl/deployments/deployment-yml/#kustomize-deployments" %}}). Everything,
including Helm Charts, ends up being processed by Kustomize and the resulting manifests are then applied to the cluster.

Before Kustomize is being invoked, [templating]({{% ref "docs/kluctl/templating/" %}}) is performed on all involved
manifests, including the `kustomization.yaml` itself. This allows advanced variable substitution (variables come from
different [sources]({{% ref "docs/kluctl/templating/variable-sources/" %}})) and even conditional
inclusion/exclusion of resources.

This also makes the Kustomizations incompatible to native Kustomize, but still allows to perform all desired actions
(diff, deploy, render, ...) via the Kluctl CLI.

## Dependency management and ordering

Kubernetes manifests are declarative and in most cases the order in which they are applied does not matter, because
constant reconciliation will eventually fix all issues that arise in-between. This, however, has limits and does not
always work. The most prominent examples are CRDs and Namespaces. You can't apply a CR before the corresponding CRD is
applied and you can't apply namespaced resources before the namespace is applied.

Reality will always force you in some way or another to deal with deployment order and dependencies. All projects
discussed in this post have completely different approaches to this problem.

### ArgoCD

ArgoCD implements [Sync Waves](https://argo-cd.readthedocs.io/en/stable/user-guide/sync-waves/), which allows you
control the order in which resources get applied. Resources with a lower wave number get applied first, then the next
highest wave number and so on. For each wave number, ArgoCD waits for healthness of each resource of the current wave.

This gives some good control about the order inside the same Application. Cross-Application sync waves or any other
dependency mechanism between Applications is [currently not supported](https://github.com/argoproj/argo-cd/issues/7437).

### Flux

Flux supports specifying dependencies
in [Kustomization](https://fluxcd.io/flux/components/kustomize/kustomizations/#dependencies)
and [HelmRelease](https://fluxcd.io/flux/components/helm/helmreleases/#dependencies). However, a Kustomization can
currently only depend on another Kustomization and a HelmRelease only on another HelmRelease. Cross-dependencies between
these are currently [not supported](https://github.com/fluxcd/flux2/issues/3364).

Health checks on Kustomizations can be specified to control when a Kustomization is considered ready so that applying
dependent Kustomizations can be deferred until readiness.

### Kluctl

Kluctl by default applies all deployment items in
the [deployment.yaml]({{% ref "docs/kluctl/deployments/deployment-yml/" %}}) in parallel to speed up deployments.
When a [barrier item]({{% ref "docs/kluctl/deployments/deployment-yml/#barriers" %}}) is encountered, Kluctl will
stop and wait for all previously encountered deployment items to fully apply before it continues with further parallel
processing.

This allows you to specify an intuitive and natural ordering. The position in the deployment item list determines the
order and thus allows you to easily specify an intent like "deploy operator X and Y with their CRDs in parallel and only
after this apply the corresponding CRs".

Additional deployment item types
like [waitReadinessObjects]({{% ref "docs/kluctl/deployments/deployment-yml/#waitreadinessobjects" %}}) allow you to
also wait for readiness of individual resources (e.g. an operator implementing a Webhook or performing delayed CRD
installation).

## Dynamic configuration

Very often, the same deployment needs to be deployed to different environments/clusters but with slightly different
configuration. This requires some form of dynamic configuration capabilities offered by the GitOps solution in use.
Performing a copy of all manifests and individually changing the differing resources is the worst option.

The next level is to dynamically create deployments (via GitOps CRs) based on some additional source (e.g. Git files or
branches). This is for example useful to create preview environments.

### ArgoCD

ArgoCD Applications can use Kustomize or Helm to perform configuration. This means, to support multiple environments,
you'd have to point the different Applications to different Kustomize overlays or Helm values.

[ApplicationSets](https://argo-cd.readthedocs.io/en/stable/user-guide/application-set/) can be used to dynamically
create dynamic Applications based on a generator (e.g. clusters or git files). These Applications can also receive some
limited set of variables.

### Flux

Flux also uses Kustomize and Helm to perform configuration. You'd either use a Kustomization CR pointing to a specific
overlay or a HelmRelease CR with dedicated Helm values.

[Post-build variable substitions](https://fluxcd.io/flux/components/kustomize/kustomizations/#post-build-variable-substitution)
can be used to inject variables sourced from ConfigMaps or Secrets. This can be used to further configure environments.

Flux itself does not provide dynamic creation of CRs (like ApplicationSets in ArgoCD). You can however use
the [template-controller]({{% ref "docs/template-controller/" %}}) from the Kluctl project to create dynamic
CRs/environments.

### Kluctl

Even though Kluctl supports the same way of configuration via plain Kustomize and Helm as the other solutions, it is
generally not recommended due to Kluctl offering its own solution. Kluctl allows to
use [templating]({{% ref "docs/kluctl/templating/" %}}) in all involved files inside your deployment project. At the
same time, dynamic [variable sources]({{% ref "docs/kluctl/templating/variable-sources/" %}}) can be used to pull in
all kinds of configuration sources (e.g. plain YAML files, Git repos, ConfigMaps, Secrets, AWSSecretsManager,
Vault, ...).

Dynamic creation of KluctlDeployment CRs can either be implemented via native Kluctl (e.g. by using
a [gitFiles]({{% ref "docs/kluctl/templating/variable-sources/#gitfiles" %}}) source) or via
the [template-controller]({{% ref "docs/template-controller/" %}}).

## Testability

In GitOps, deployments are triggered by commiting and pushing to Git. To avoid pushing broken deployments, some form of
testing and/or verification must be present to avoid breaking production (or any other sensible environment).

Using pull requests with proper reviews is usually a good first step to catch many broken deployments. There are however
many cases where such a process is not enough, because many effects of changes are not obvious just by looking at the
changed manifests/configuration. A prominent example is the renaming of resources, worst case being a namespace, which
actually causes deletion and re-creation of resources.

The next step is to introduce staging/testing/preview environments which are used to test out changes first. This
however requires a more complex project structure and release process. It can also easily multiply infrastructure costs,
making it even harder for smaller companies or teams.

### ArgoCD

In addition to the already described approaches, ArgoCD supports manual syncs to prevent breaking sensible environments.
With manual syncs, you can push/merge to the production branch and then review the effects on the cluster in the ArgoCD
UI. Only if they look good, you would then trigger the sync manually.

The [argocd app diff](https://argo-cd.readthedocs.io/en/stable/user-guide/commands/argocd_app_diff/) CLI command can
also be used to diff a local version of your deployment against the live state in the cluster. This is useful if you
want to verify that a change looks good before commiting and pushing it.

### Flux

Flux does not support something like manual syncs, which means that if you push/merge to Git, it will be applied to the
target environment. This forces you to ensure correctness of deployments before the merge/push happens, e.g. by using
one of the previously mentioned methods (pull request reviews and/or staging/prod environments).

The [flux diff kustomization](https://fluxcd.io/flux/cmd/flux_diff_kustomization/) CLI command is comparable to
the `argocd app diff` command described for ArgoCD. It does not support deep diffing of nested Kustomizations. Helm
diffs are also not supported.

### Kluctl

Kluctl supports [manual deployments]({{% ref "docs/gitops/spec/v1beta1/kluctldeployment/#manual" %}}) which must be
approved via the UI (which shows the diff) before they are deployed. This can be used the same way as manual syncs in
ArgoCD.

The [kluctl diff]({{% ref "docs/kluctl/commands/diff/" %}}) command can be used to locally dry-run apply and diff a
local version of your deployment project. It fully supports nested Kluctl deployment projects.

In addition, [kluctl gitops diff]({{% ref "docs/kluctl/commands/gitops-diff/" %}}) can be used to instruct the
controller to perform a diff, based on the local version of the deployment project.

## UI

Even though GitOps is centered around Git, a UI can be quite useful for additional control and monitoring.

### ArgoCD

ArgoCD is well known for its UI and clearly the forerunner here. Everything can be controlled via the UI, including
cluster and application management. You can see the state of all applications and the corresponding resources. You can
manually sync applications, look into logs of PODs, and much more.

### Flux

Flux does not offer an official UI. There are however some non-official UIs available, for
example [Capacitor](https://fluxcd.io/blog/2024/02/introducing-capacitor/). These are however still very rudimentary and
only allow simple monitoring and actions (e.g. reconcile).

### Kluctl

Kluctl offers an official but still [experimental UI](https://kluctl.io/blog/2023/09/12/introducing-the-kluctl-webui/).
It allows you to monitor and control your KluctlDeployments. You can suspend, resume, reconcile, prune, approve, ...
your deployments. It shows you historical deployment results (with diffs), current drift, validation state, and much
more.

As mentioned, it's still experimental but it already showcases the future potential.
