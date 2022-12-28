
---
title: "Introducing the Template Controller and building GitOps Preview Environments"
linkTitle: "Introducing the Template Controller and building GitOps Preview Environments"
slug: template-controller
date: 2022-12-28
author: Alexander Block (@codablock)
images:
- "/images/blog/images/blog/template-controller.jpg"
---

![image](/images/blog/template-controller.jpg)

This blog post serves two purposes. The first one is to announce and present the 
[Template Controller](https://kluctl.io/docs/template-controller/) ([Source](https://github.com/kluctl/template-controller)).
The second purpose is to demonstrate it by setting up a simple GitOps based Kubernetes deployment with dynamic preview
environments.

## The Template Controller
The template-controller is a Kubernetes controller that is able to create arbitrary objects based on dynamic templates
and arbitrary input objects. It is inspired by ArgoCD's [ApplicationSet](https://argo-cd.readthedocs.io/en/stable/operator-manual/applicationset/),
which is able to create dynamic ArgoCD [Applications](https://argo-cd.readthedocs.io/en/stable/operator-manual/declarative-setup/#applications)
from a list of generators (e.g. Git) and an `Application` template.

The Template Controller uses a different approach, making it more flexible and independent of the GitOps system being
used. It uses arbitrary Kubernetes objects as inputs and allows to create templated objects of any kind 
(e.g. a Flux Helm Release or a [KluctlDeployment](https://kluctl.io/docs/flux/spec/v1alpha1/kluctldeployment/)).
This makes the controller very extensible, as any type of input can be implemented with the help of additional
controllers which are not necessarily part of the project.

When specifying the input objects, you'd also specify which part of the object to use as input. This is done by
specifying a [JSON Path](https://goessner.net/articles/JsonPath/) that select the subfield of the object to use, e.g.
`status.result` for a [GitProjector](https://kluctl.io/docs/template-controller/spec/v1alpha1/gitprojector/) or
`data` for a ConfigMap.

The Template Controller implements this functionality through the [ObjectTemplate](https://kluctl.io/docs/template-controller/spec/v1alpha1/objecttemplate/)
CRD. As the name implies, it also uses a templating engine, which is identical to the one used in
[Kluctl](https://kluctl.io/docs/kluctl/reference/templating/), with the `ObjectTemplate's` input matrix available as
global variables.

## Preparation
To try the examples provided in this blog post, you'll need to have a running cluster ready. You could for example use
a local [kind](https://kind.sigs.k8s.io/) cluster:

```shell
$ kind create cluster
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.25.3) 🖼 
 ✓ Preparing nodes 📦  
 ✓ Writing configuration 📜 
 ✓ Starting control-plane 🕹️ 
 ✓ Installing CNI 🔌 
 ✓ Installing StorageClass 💾 
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Thanks for using kind! 😊
```

Then, you will need the Template Controller installed into this cluster:

```shell
$ helm repo add kluctl http://kluctl.github.io/charts
"kluctl" has been added to your repositories
$ helm install -n kluctl-system --create-namespace template-controller kluctl/template-controller
NAME: template-controller
LAST DEPLOYED: Wed Dec 28 17:24:47 2022
NAMESPACE: default
STATUS: deployed
REVISION: 1
```

You will also need the [Flux Kluctl Controller](https://github.com/kluctl/flux-kluctl-controller) installed for the
examples. If you decide to use plain Flux deployments, you will need to [install Flux](https://fluxcd.io/flux/installation/) instead.

```shell
$ # we assume that you have the Helm repository installed already
$ helm install -n kluctl-system --create-namespace flux-kluctl-controller kluctl/flux-kluctl-controller 
NAME: flux-kluctl-controller
LAST DEPLOYED: Wed Dec 28 17:27:53 2022
NAMESPACE: default
STATUS: deployed
REVISION: 1
```

## A few words on security
The Template Controller can access and create any kind of Kubernetes resource. This makes it very powerful but also
very dangerous. The Template Controller needs to run with a quite privileged service account but at the same time uses user
impersonation to downgrade permissions to less privileged service accounts while processing `ObjectTemplates`.

By default, the controller will use the `default` service account of the namespace that the `ObjectTemplate` is deployed
to. This will by default limit permissions to basically nothing, unless you explicitly bind roles to the `default`
service account. A better way is to create a dedicated service account instead and bind limited roles with the required
permissions to that dedicated service account. This service account can then be specified in the `ObjectTemplate` via
`spec.serviceAccountName`.

The following RBAC rules can be used when going through the examples in this blog post, it is however suggested to not
blindly reuse them in your real deployments. You should carefully asses which permissions are really needed and limit
the roles appropriately. Also pay attention when using the `cluster-admin` role or any other `ClusterRole`, as it easily
allows privilege escalation and at least allows to deploy into other namespaces than the `ObjectTemplate's` namespace.

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: template-controller
  namespace: default
---
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: template-controller
  namespace: default
rules:
  - apiGroups: ["*"]
    resources: ["*"]
    verbs: ["*"]
---
kind: RoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: template-controller
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: Role
  name: template-controller
subjects:
  - kind: ServiceAccount
    name: template-controller
    namespace: default
```

## Dynamic Preview Environments
The initial use case and the main reason the Template Controller was created was to implement preview environments
in a GitOps setup. Preview environments are dynamic environments that are spinned up and down on demand, for example
when a pull request introduces changes that need to be tested before these are merged into the main branch.

In a GitOps-based setup, one would need to create the relevant custom resources per preview environment, for example
a [Flux Kustomization](https://fluxcd.io/flux/components/kustomize/kustomization/),
[Flux HelmRelease](https://fluxcd.io/flux/components/helm/helmreleases/) or a
or [KluctlDeployment](https://kluctl.io/docs/flux/spec/v1alpha1/kluctldeployment/). The underlying GitOps controller
would then take over and perform the actual deployment.

In the following examples we will concentrate on using `KluctlDeployments`. Changing it to use `Kustomizations` or
`HelmReleases` should be self-explanatory, as the Template Controller treats all resource kinds equally.

There are multiple options on how to define the desired pre-conditions and configuration of a preview environment, which will
be described in the next chapters.

## Linking Preview Environments to Branches
You can, for example, link Git branches to preview environments, so that for each new branch a preview environment is
created, with configuration being read from a yaml file inside the branch. This can be achieved by using a
[GitProjector](https://kluctl.io/docs/template-controller/spec/v1alpha1/gitprojector/), which will periodically
clone the configured Git repository, scan for matching branches and files and project the result into the `GitProjector`
status. The status is then available as matrix input inside the `ObjectTemplate`.

Example `GitProjector`:

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: GitProjector
metadata:
  name: preview
  namespace: default
spec:
  interval: 1m
  # You'll need to fork this repository and update the url to actually test what happens when you create new branches
  url: https://github.com/kluctl/kluctl-examples.git
  # In case you use a private repository:
  # secretRef:
  #   name: git-credentials
  ref:
    # let's only take preview- branches into account
    branch: preview-.*
```

The following `ObjectTemplate` can then use the `GitProjetor's` `status.result` field to create one `KluctlDeployment`
per branch.

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: ObjectTemplate
metadata:
  name: preview
  namespace: default
spec:
  # this is the service account described at the top
  serviceAccountName: template-controller
  prune: true
  matrix:
    - name: git
      object:
        ref:
          apiVersion: templates.kluctl.io/v1alpha1
          kind: GitProjector
          name: preview
        jsonPath: status.result
        # status.result is a list of matching branches. Without expandLists being true, the matrix input would treat
        # that list as a single input, but we actually want the list items being one input per item
        expandLists: true
  templates:
    - object:
        apiVersion: flux.kluctl.io/v1alpha1
        kind: KluctlDeployment
        metadata:
          # We can use the branch name as basis for the object name
          name: "{{ matrix.git.ref.branch | slugify }}"
          namespace: default
        spec:
          interval: 1m
          deployInterval: "never"
          timeout: 5m
          source:
            url: https://github.com/kluctl/kluctl-examples.git
            path: simple-helm
            ref:
              # ensure we deploy from the correct branch
              branch: "{{ matrix.git.ref.branch }}"
            # secretRef:
            #  name: git-credentials
          target: "simple-helm"
          context: default
          args:
            # Look into the simple-helm example to figure out what the environment arg does. It basically sets the
            # namespace to use, but could theoretically do much more.
            environment: "{{ matrix.git.ref.branch | slugify }}"
          prune: true
```

Please note the use of templating variables and filters inside the actual template under `spec.templates`. Each template
will be rendered once per matrix input, which in this case means once per branch. The templates can use the current
matrix input in the form of a variable, accessible via `matrix.<name>`, e.g. `matrix.git` in this case. Please read
the documentation of [GitProjector](https://kluctl.io/docs/template-controller/spec/v1alpha1/gitprojector/) to figure
out what is available in `matrix.git`, which is basically just a copy of the individual `status.result` list items.

## One Preview Environment per pull requests
Another option is to use pull requests instead of the underlying Git branches to create preview environments. This might
be useful if you want to report the status of your preview environment to the pull request, e.g. by updating the
commit status when the deployment turns green or red. One might also want to post complex status comments, for example
the result of the deployment in the form of structured and beautiful diff.

To achieve this, you can use the [ListGithubPullRequests](https://kluctl.io/docs/template-controller/spec/v1alpha1/listgithubpullrequests/)
or the [ListGitlabMergeRequests](https://kluctl.io/docs/template-controller/spec/v1alpha1/listgitlabmergerequests/)
custom resource, which are also provided by the Template Controller.

Consider the following example:

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: ListGithubPullRequests
metadata:
  name: preview
  namespace: default
spec:
  interval: 1m
  owner: kluctl
  repo: kluctl-examples
  state: open
  base: main
  head: kluctl:preview-.*
  # in case you forked the repo into a private repo
  #tokenRef:
  #  secretName: git-credentials
  #  key: github-token
```

The above example will regularly (1m interval) query the GitHub API for pull requests inside the kluctl-examples
repository. It will filter for pull requests which are in the state "open" and are targeted against the "main" branch.
The result of the query is then stored in
the `status.pullRequests` field of the custom resource. The content of the pullRequests field basically matches what
GitHub would return via the [Pulls API](https://docs.github.com/en/rest/pulls/pulls) (with some fields omitted to
reduce the size).

You can now use the result in an `ObjectTemplate`:

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: ObjectTemplate
metadata:
  name: pr-preview
  namespace: default
spec:
  # this is the service account described at the top
  serviceAccountName: template-controller
  prune: true
  matrix:
    - name: pr
      object:
        ref:
          apiVersion: templates.kluctl.io/v1alpha1
          kind: ListGithubPullRequests
          name: preview
        jsonPath: status.pullRequests
        # status.result is a list of matching branches. Without expandLists being true, the matrix input would treat
        # that list as a single input, but we actually want the list items being one input per item
        expandLists: true
  templates:
    - object:
        apiVersion: flux.kluctl.io/v1alpha1
        kind: KluctlDeployment
        metadata:
          # We can use the head branch name as basis for the object name
          name: pr-{{ matrix.pr.head.ref | slugify }}
          namespace: default
        spec:
          interval: 1m
          deployInterval: "never"
          timeout: 5m
          source:
            url: https://github.com/kluctl/kluctl-examples.git
            path: simple-helm
            ref:
              # ensure we deploy from the correct branch
              branch: "{{ matrix.pr.head.ref }}"
            # secretRef:
            #  name: git-credentials
          target: "simple-helm"
          context: default
          args:
            # Look into the simple-helm example to figure out what the environment arg does. It basically sets the
            # namespace to use, but could theoretically do much more.
            environment: "pr-{{ matrix.pr.head.ref | slugify }}"
          prune: true
```

## Using Flux instead of Kluctl
The above examples can easily be changed to work with Flux resources, e.g. `Kustomization` and `HelmRelease`. Simply
replace the template `KluctlDeployment` with the appropriate resources and use the same template variables where needed.

## Shutting down preview environments
All the above examples will result in pruning the whole deployment when the branch gets deleted. This works because
after deletion of the branch, the `GitProjector` will not contain that branch in the `status.result` after the next
reconciliation. This will lead to the `ObjectTemplate` not having the matrix input anymore, causing a prune of the
`KluctlDeployment` object, which in turn causes deletion of all resources deployed by the Kluctl controller.

For the pull request based examples, it will also prune the deployments when the pull request gets closed or merged,
simply because the `ListGithubPullRequests` filters for open pull requests and thus will not list the pull
requests anymore after they have been closed/merged.

## Another use case: Transformation of objects
The following is a very simple example of an `ObjectTemplate` that uses a `Secret` as input and simply transforms it
into another secret.

This can turn out to be quite useful if you need to re-use the same secret value multiple times
but in different forms. For example, if you have a secret that stores database credentials, you might also need the same
username and password inside a JDBC url. Typically, you'd have to store the secret twice in both required forms.

This is however not an option if the secret is generated after the deployment, e.g. by the Zalando Postgres operator
or by the AWS RDS operator. Using the [Mittwald Secret Generator](https://github.com/mittwald/kubernetes-secret-generator)
is also something that I can recommend, as it removes the need to pre-generate secrets. With the `ObjectTemplate` based
transformation, the generated secrets can then be used in many new ways.

Consider the following input `Secret`:

```yaml
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  name: input-secret
  namespace: default
stringData:
  password: my-user
  username: secret-password
```

And the following `ObjectTemplate`:

```yaml
apiVersion: templates.kluctl.io/v1alpha1
kind: ObjectTemplate
metadata:
  name: transformer-template
  namespace: default
spec:
  serviceAccountName: template-controller
  prune: true
  matrix:
    - name: secret
      object:
        ref:
          apiVersion: v1
          kind: Secret
          name: input-secret
  templates:
  - object:
      apiVersion: v1
      kind: Secret
      metadata:
        name: "transformed-secret"
      stringData:
        jdbc_url: "jdbc:postgresql://db-host:5432/db?user={{ matrix.secret.data.username | b64decode }}&password={{ matrix.secret.data.password | b64decode }}"
        # sometimes the key names inside a secret are not what another component requires, so we can simply use different names if we want
        username_with_different_key: "{{ matrix.secret.data.username | b64decode }}"
        password_with_different_key: "{{ matrix.secret.data.password | b64decode }}"
```

I hope the example above is self-explanatory. It simply transforms one secret into another, which is then in a form
that can be consumed by a Java application for example that can only work with JDBC urls.

## What's next?
I believe that there are much more use cases for the Template Controller, and I'm absolutely convinced that the community
will invent completely new ones and maybe share them by posting examples and ideas. Due to the flexible nature of the
matrix inputs and template definition, a lot is possible. Think big! 😸

The Template Controller currently comes with a few additional custom resources 
(e.g. `GitProjector` and `ListGithubPullRequests`) that are meant to be used as matrix inputs. I can imagine that
other custom resources might be candidates to be included in the controller as well, and I'm also open for ideas.