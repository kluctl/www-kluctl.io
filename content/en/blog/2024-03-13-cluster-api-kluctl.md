---
title: "Managing Cluster API with Kluctl"
linkTitle: "Managing Cluster API with Kluctl"
slug: cluster-api-kluctl
date: 2024-03-13
author: Alexander Block (@codablock)
description: |
  A tutorial on how to use Kluctl to manage Cluster API based clusters.
tags:
- Tutorial
images:
- "/images/blog/2024-03-13-cluster-api-kluctl.jpg"
---

![image](/images/blog/2024-03-13-cluster-api-kluctl.jpg)

Kubernetes started as a very promising container orchestrator and in my opinion it was very clear at day one that it would establish itself and take the market. What was not so obvious to me, was that Kubernetes would also morph into some kind of "API Platform".

With the introduction of [Custom Resource Definitions](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/), all kinds of resources could now be managed by Kubernetes. [Controllers and Operators](https://book.kubebuilder.io/cronjob-tutorial/controller-overview.html#whats-in-a-controller) take these Custom Resources and use the reconcile pattern to constantly reconcile the desired state with the real world.

The next step was obvious in hindsight, but still a surprise for me personally: Why not manage Kubernetes Clusters itself from inside Kubernetes Clusters. [Cluster API](https://cluster-api.sigs.k8s.io/) was born.

## Implications of Custom Resources
Having something in the form of a Custom Resource also means that it becomes a regular Kubernetes Resource that can be managed with all available tooling in the Kubernetes ecosystem. It can be managed with plain Kubectl, but also with more advances tools like [Helm](https://helm.sh/), [Flux](https://fluxcd.io/), [ArgoCD](https://argo-cd.readthedocs.io/en/stable/) or [Kluctl](https://kluctl.io).

## So, why Kluctl?
Kluctl is general purpose deployment tool for Kubernetes. It allows you to define Kubernetes deployments of any complexity and manage them via a [unified CLI]({{% ref "docs/kluctl/commands" %}}) and/or an optional [GitOps controller]({{% ref "docs/gitops" %}}). Here a are a few features that make Kluctl interesting for the management of Cluster API based clusters.

1. [Targets]({{% ref "docs/kluctl/kluctl-project/targets" %}}) allow you to manage multiple workload clusters with the same Kluctl deployment project.
2. [Templating]({{% ref "docs/kluctl/templating" %}}) allows you to follow a natural project structure, without the need to use overlays and patching to meet simple requirements.
3. [Deployment projects]({{% ref "docs/kluctl/deployments/deployment-yml" %}}) allow you to reuse parametrised and templated subcomponents without copy-paste.
4. [Variable sources]({{% ref "docs/kluctl/templating/variable-sources" %}}) give you easy to understand and structured configuration for the workload clusters.
5. The [Kluctl diff]({{% ref "docs/kluctl/commands/diff" %}}) command will always tell you if you're good or not when you change things (because it's based on a server-side dry-run).
6. [GitOps]({{% ref "docs/gitops" %}}) is fully supported but also optional. It can even be [mixed]({{% ref "docs/kluctl/commands/gitops-deploy" %}}) with a classical push style CLI.

## Installing Kluctl

For this tutorial, you'll need the Kluctl CLI installed. Please follow the instructions [here]({{% ref "docs/kluctl/installation#installing-the-cli" %}}). There is no need to install the GitOps controller or the Webui, but feel free to try these out as well after the tutorial.

## Let's setup cluster-api
In this tutorial, we'll work completely locally without any cloud resources being involved. This means, we're using [Kind](https://kind.sigs.k8s.io/) and the CAPD (Cluster API Docker) infrastructure provider. In the real world, you'll need to adapt the principles learned here to a proper Cluster API infrastructure provider.

First, lets set up a local Kind cluster. If you don't have Kind installed yet, read through the [installation instructions](https://kind.sigs.k8s.io/#installation-and-usage).

The CAPD provider will need access to the host Docker daemon from inside the Kind cluster. To give access, you'll need to pass through the Docker unix socket. This can be done by using a custom Kind configuration:

```yaml
# contents of kind-config.yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  ipFamily: dual
nodes:
- role: control-plane
  extraMounts:
    - hostPath: /var/run/docker.sock
      containerPath: /var/run/docker.sock
````

Now create the cluster with the above config:

```bash
$ kind create cluster --config kind-config.yaml
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.29.2) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a nice day! 👋
````

The current Kubernetes Context will be set to kind-kind, which is what we'll from now on use to install Cluster API to. Let's do that:

```bash
$ clusterctl init --infrastructure docker
Fetching providers
Installing cert-manager Version="v1.13.2"
Waiting for cert-manager to be available...
Installing Provider="cluster-api" Version="v1.6.1" TargetNamespace="capi-system"
Installing Provider="bootstrap-kubeadm" Version="v1.6.1" TargetNamespace="capi-kubeadm-bootstrap-system"
Installing Provider="control-plane-kubeadm" Version="v1.6.1" TargetNamespace="capi-kubeadm-control-plane-system"
Installing Provider="infrastructure-docker" Version="v1.6.1" TargetNamespace="capd-system"

Your management cluster has been initialized successfully!

You can now create your first workload cluster by running the following:

  clusterctl generate cluster [name] --kubernetes-version [version] | kubectl apply -f -
```

We now have a fully functional Cluster API installation that is able to provision and manage workload clusters in the form of Docker Containers.

## Basic project structure
Let's talk about the basic Kluctl project structure that we'll follow for this tutorial. You can find the full project at https://github.com/kluctl/cluster-api-demo. This repository contains multiple subdirectories with different versions of the project. The first version, as described in this and the next section, is inside `1-initial`.

The root directory will contain 2 files.

The first one is the [.kluctl.yaml]({{% ref "docs/kluctl/kluctl-project" %}}) file, which specifies which [targets]({{% ref "docs/kluctl/kluctl-project/targets" %}}) exists. A target defines where/what to deploy with a Kluctl project and can be anything you want. In a classical application deployment, it would be the target environment. In this case, a target represents a Cluster API workload cluster, deployed to a Cluster API management cluster (our Kind cluster). It serves as the entrypoint to configuration management and will later allow us to load cluster specific configuration.

```yaml
# content of .kluctl.yaml
targets:
  - name: demo-1
    context: kind-kind
  - name: demo-2
    context: kind-kind

discriminator: capi-{{ target.name }}
````

You can also see the first use of templating here in the discriminator. The discriminator is later used to differentiate resources that have been applied to the cluster before. This is useful for cleanup tasks like pruning or deletion.

The second file is the [deployment.yaml]({{% ref "docs/kluctl/deployments/deployment-yml" %}}), which defines the actual deployment project. It includes Kustomize deployments, Helm Charts and other sub-deployment projects.

```yaml
# content of deployment.yaml
deployments:
  - include: clusters

commonAnnotations:
  kluctl.io/force-managed: "true"
````

This will include a sub-deployment in the directory "clusters". Inside this directory, there must be another deployment.yaml. The annotation found in [commonAnnotations]({{% ref "docs/kluctl/deployments/deployment-yml#commonannotations" %}}) will cause Kluctl to [always consider]({{% ref "docs/kluctl/deployments/annotations/all-resources#kluctlioforce-managed" %}}) resources as managed by Kluctl. This is required because Cluster API [claims ownership of resources](https://github.com/kubernetes-sigs/cluster-api/issues/5487#issuecomment-950824947) even though it is not in control of those.

```yaml
# content of clusters/deployment.yaml
deployments:
  - path: {{ target.name }}
````

This will include a [Kustomize]({{% ref "docs/kluctl/deployments/deployment-yml#kustomize-deployments" %}}) deployment from the directory that is resolved via the template `{{ target.name }}`. "target" is a global variable that is always present, and it allows you to access the properties used in the current target, defined in the `.kluctl.yaml` from above. This means, if you later deploy the target "demo-1", Kluctl will load the Kustomize deployment found in the "clusters/demo-1" folder.

## Creating a workload cluster
Now, create the following files in the clusters/demo-1 directory:

```yaml
# contents of clusters/demo-1/kustomization.yaml
resources:
  - namespace.yaml
  - cluster.yaml
  - control-plane.yaml
  - workers.yaml
```

The above file is a regular `kustomization.yaml` that includes the actual resources. Kluctl fully supports [Kustomize]({{% ref "docs/kluctl/deployments/kustomize" %}}) and all its features. You can also omit the `kustomization.yaml` in most cases, causing Kluctl to [auto-generate]({{% ref "docs/kluctl/deployments/deployment-yml#simple-deployments" %}}) the kustomization.yaml. In this case however, this is not recommended as the order is important here: The namespace must be deployed before anything else.

```yaml
# contents clusters/demo-1/namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: cluster-demo-1
```

We create a dedicated namespace for this cluster. We will also create more namespaces later for every other cluster.

```yaml
# contents of clusters/demo-1/cluster.yaml
apiVersion: cluster.x-k8s.io/v1beta1
kind: Cluster
metadata:
  name: "demo-1"
  namespace: "cluster-demo-1"
spec:
  clusterNetwork:
    services:
      cidrBlocks: ["10.128.0.0/12"]
    pods:
      cidrBlocks: ["192.168.0.0/16"]
    serviceDomain: "cluster.local"
  infrastructureRef:
    apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
    kind: DockerCluster
    name: "demo-1"
    namespace: "cluster-demo-1"
  controlPlaneRef:
    kind: KubeadmControlPlane
    apiVersion: controlplane.cluster.x-k8s.io/v1beta1
    name: "demo-1-control-plane"
    namespace: "cluster-demo-1"
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: DockerCluster
metadata:
  name: "demo-1"
  namespace: "cluster-demo-1"
```

The above file describes a [Cluster](https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api/cluster.x-k8s.io/Cluster/v1beta1@v1.6.2) and a [DockerCluster](https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api/infrastructure.cluster.x-k8s.io/DockerCluster/v1beta1@v1.6.2). Please note that we are not using Cluster Topology ([ClusterClass](https://cluster-api.sigs.k8s.io/tasks/experimental-features/cluster-class/)) features. I will later explain why.

```yaml
# contents of clusters/demo-1/control-plane.yaml
kind: KubeadmControlPlane
apiVersion: controlplane.cluster.x-k8s.io/v1beta1
metadata:
  name: "demo-1-control-plane"
  namespace: "cluster-demo-1"
spec:
  replicas: 1
  machineTemplate:
    infrastructureRef:
      kind: DockerMachineTemplate
      apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
      name: "demo-1-control-plane"
      namespace: "cluster-demo-1"
  kubeadmConfigSpec:
    clusterConfiguration:
      controllerManager:
        extraArgs:
          enable-hostpath-provisioner: 'true'
      apiServer:
        certSANs: [localhost, 127.0.0.1, 0.0.0.0]
    initConfiguration:
      nodeRegistration:
        criSocket: /var/run/containerd/containerd.sock
        kubeletExtraArgs:
          cgroup-driver: systemd
          eviction-hard: 'nodefs.available<0%,nodefs.inodesFree<0%,imagefs.available<0%'
    joinConfiguration:
      nodeRegistration:
        criSocket: /var/run/containerd/containerd.sock
        kubeletExtraArgs:
          cgroup-driver: systemd
          eviction-hard: 'nodefs.available<0%,nodefs.inodesFree<0%,imagefs.available<0%'
  version: "1.29.0"
---
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: DockerMachineTemplate
metadata:
  name: "demo-1-control-plane"
  namespace: "cluster-demo-1"
spec:
  template:
    spec:
      extraMounts:
        - containerPath: "/var/run/docker.sock"
          hostPath: "/var/run/docker.sock"
```

The above file describes a [KubeadmControlPlane](https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api/controlplane.cluster.x-k8s.io/KubeadmControlPlane/v1beta1@v1.6.2) and a [DockerMachineTemplate](https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api/infrastructure.cluster.x-k8s.io/DockerMachineTemplate/v1beta1@v1.6.2) for the control plane nodes.

```yaml
# contents of clusters/demo-1/workers.yaml
apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
kind: DockerMachineTemplate
metadata:
  name: "demo-1-md-0"
  namespace: "cluster-demo-1"
spec:
  template:
    spec: {}
---
apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
kind: KubeadmConfigTemplate
metadata:
  name: "demo-1-md-0"
  namespace: "cluster-demo-1"
spec:
  template:
    spec:
      joinConfiguration:
        nodeRegistration:
          kubeletExtraArgs:
            cgroup-driver: systemd
            eviction-hard: 'nodefs.available<0%,nodefs.inodesFree<0%,imagefs.available<0%'
---
apiVersion: cluster.x-k8s.io/v1beta1
kind: MachineDeployment
metadata:
  name: "demo-1-md-0"
spec:
  clusterName: "demo-1"
  replicas: 1
  selector:
    matchLabels:
  template:
    spec:
      clusterName: "demo-1"
      version:  "1.29.0"
      bootstrap:
        configRef:
          name: "demo-1-md-0"
          namespace: "cluster-demo-1"
          apiVersion: bootstrap.cluster.x-k8s.io/v1beta1
          kind: KubeadmConfigTemplate
      infrastructureRef:
        name: "demo-1-md-0"
        namespace: "cluster-demo-1"
        apiVersion: infrastructure.cluster.x-k8s.io/v1beta1
        kind: DockerMachineTemplate
```

The above file describes everything needed to create a pool of nodes. This includes a [DockerMachineTemplate](https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api/infrastructure.cluster.x-k8s.io/DockerMachineTemplate/v1beta1@v1.6.2), a [KubeadmConfigTemplate](https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api/bootstrap.cluster.x-k8s.io/KubeadmConfigTemplate/v1beta1@v1.6.2) and a [MachineDeployment](https://doc.crds.dev/github.com/kubernetes-sigs/cluster-api/cluster.x-k8s.io/MachineDeployment/v1beta1@v1.6.2).

## Deploying the workload cluster

We now have a working Kluctl Deployment Project that can be deployed via the [Kluctl CLI]({{% ref "docs/kluctl/commands" %}}) (we will later also explore GitOps). Execute the following command:

```bash
$ kluctl deploy -t demo-1
```

This will perform a dry-run, show the diff and then after confirmation do the actual deployment. The dry-run will produce a few errors as the underlying server-side dry-run is not perfect in combination with Cluster API, you can ignore these errors and simply confirm.

After a few minutes, the workload cluster should be ready with one control-plane node and one worker node, all running as Docker containers. We now need to get the kubeconfig of this cluster.

```bash
$ kind get kubeconfig --name demo-1 > demo-1.kubeconfig
```

You can now test access to the workload cluster:
```bash
$ kubectl --kubeconfig=demo-1.kubeconfig get node
NAME                         STATUS     ROLES           AGE   VERSION
demo-1-control-plane-bjfvn   NotReady   control-plane   47m   v1.29.0
demo-1-md-0-mtcpn-wnb8v      NotReady   <none>          21m   v1.29.0
```

This will reveal that the cluster is currently not fully functional, simply because a working [CNI](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/) is missing. To install a CNI, run:

```bash
$ kubectl --kubeconfig=./demo-1.kubeconfig \
    apply -f https://raw.githubusercontent.com/projectcalico/calico/v3.26.1/manifests/calico.yaml
```

After a few seconds, re-running the above `get node` command will show that nodes are ready.

## Modifying the workload cluster

You can now try to modify something in the workload cluster manifests.

Lets increase the workers `MachineDeployment` replicas to 2. You can do this by editing `clusters/demo-1/workers.yaml` with your favorite editor, search for the `MashineDeployment` resource and replace `replicas: 1` with `replicas: 2`.

Now, let's [deploy]({{% ref "docs/kluctl/commands/deploy" %}}) this change. We will now start to see the first benefits from Kluctl, specifically the dry-run and diff that happens before we deploy something. You will need to confirm the deployment by pressing `y`.

```bash
$ kluctl deploy -t demo-1
✓ Loading kluctl project-api-demo/1-initial
✓ Initializing k8s client
✓ Rendering templates
✓ Rendering Helm Charts
✓ Building kustomize objects
✓ Postprocessing objects
✓ Writing rendered objects
✓ Getting remote objects by discriminator
✓ Getting namespaces
✓ demo-1: Applied 8 objects.

Changed objects:
  cluster-demo-1/MachineDeployment/demo-1-md-0


Diff for object cluster-demo-1/MachineDeployment/demo-1-md-0
+---------------+----------------------------------------------------------------------------------+
| Path          | Diff                                                                             |
+---------------+----------------------------------------------------------------------------------+
| spec.replicas | -1                                                                               |
|               | +2                                                                               |
+---------------+----------------------------------------------------------------------------------+
✓ The diff succeeded, do you want to proceed? (y/N) y
✓ demo-1: Applied 8 objects.
✓ Writing command result

Changed objects:
  cluster-demo-1/MachineDeployment/demo-1-md-0


Diff for object cluster-demo-1/MachineDeployment/demo-1-md-0
+---------------+----------------------------------------------------------------------------------+
| Path          | Diff                                                                             |
+---------------+----------------------------------------------------------------------------------+
| spec.replicas | -1                                                                               |
|               | +2                                                                               |
+---------------+----------------------------------------------------------------------------------+
```

If you check the Cluster API management cluster, you will see that a new node will appear now.

```bash
$ kubectl --kubeconfig=demo-1.kubeconfig get node
demo-1-control-plane-bjfvn   Ready      control-plane   12h   v1.29.0
demo-1-md-0-mtcpn-n2jdt      NotReady   <none>          20s   v1.29.0
demo-1-md-0-mtcpn-wnb8v      Ready      <none>          12h   v1.29.0
```

## Add and remove node pools
You can also try more types of modifications. It gets especially interesting when you start to add or remove resources, for example if you add another node pool by copying `workers.yaml` to `workers-2.yaml` (don't forget to also update `kustomization.yaml`) and replace all occurrences of `md-0` with `md-1`. When you deploy this, Kluctl will show you that new resources will be created and actually create these after confirmation.

If you tried this, also try to delete `workers-2.yaml` again and then see what `kluctl deploy -t demo-1` will do. It will inform you about the orphaned resources, which you then can [prune]({{% ref "docs/kluctl/commands/prune" %}}) via `kluctl prune -t demo-1`. Pruning can also be combined with deploying via `kluctl deploy -t demo-1 --prune`. We won't get into more detail at this point, because this will get more clear and powerful when we combine this with templating in the next section.

## Introducing templating
So far, we've only used very static manifests. To introduce new clusters, or even node pools, we'd have to do a lot of copy-paste while replacing names everywhere. This is of course not considered best practice and we should seek for a better way. Cluster API has an experimental feature called [cluster classes](https://cluster-api.sigs.k8s.io/tasks/experimental-features/cluster-class/) which tries to solve this problem. We'll however not use these in this tutorial and instead rely on Kluctl's templating functionality to solve the same requirements. A later section will also explain why templating is a viable alternative to ClusterClass.

The following changes to the project structure and files can also be found in the same [repository](https://github.com/kluctl/cluster-api-demo) already mentioned before, but inside the `2-templating` directory.

## Preparing some templated deployments
We will now introduce two reusable and templated Kustomize deployments for the cluster iteself and its workers. The cluster deployment is meant to be included once for per cluster. The workers deployment can be included multiple times, depending on how many different worker node pools you need.

Let's start by moving `kustomization.yaml`, `namespace.yaml`, `cluster.yaml` and `control-plane.yaml` into `templates/cluster/`. Also remove `workers.yaml` from the resources list in `kustomization.yaml`. This will be the cluster deployment.

Now, replace all occurrences of `demo-1` with `{{ cluster.name }}` in all the manifests found in the `templates/cluster` directory. Also, in the `KubeadmControlPlane` inside `control-plane.yaml`, replace `replicas: 1` with `{{ cluster.replicas }}`. This introduces some simple [Jinja2 templating]({{% ref "docs/kluctl/templating" %}}) to inject the cluster name. The global `cluster` variable seen here will be introduced later.

Next, move the `workers.yaml` manifest into `templates/workers`. This time, there is no need for a `kustomization.yaml` as we don't care about deployment order here (there is no namespace involved here), which means we can allow Kluctl to [auto-generate]({{% ref "docs/kluctl/deployments/deployment-yml#simple-deployments" %}}) it. Then, replace all occurences of `demo-1` with `{{ cluster.name }}` and all occurrences of `md-0` with `{{ workers.name }}`. Finally, find `replicas: 1` (or whatever you set it to before) and replace it with `replicas: {{ workers.replicas }}`.

Please note that this tutorial keeps the amount of configuration possible in these deployments to a minimum. You can maybe imagine that a lot can be achieved via templating here. For example, AWS or Azure instance types could be configured via `{{ workers.instanceType }}`.

Also, a real world example might consider putting the cluster/worker templates in seprate git repositories and including them via [git]({{% ref "docs/kluctl/deployments/deployment-yml#git-includes" %}}) or [oci]({{% ref "docs/kluctl/deployments/deployment-yml#oci-includes" %}}) includes. Both will allow you to implement versioning and other best practices for the templates.

## Using the templated deployments
The previously prepared templated deployments can now be included as often as you want, with different configuration.

For this to work, we must however change the `clusters/demo-1` Kustomize deployment to become an [included sub-deployment]({{% ref "docs/kluctl/deployments/deployment-yml#includes" %}}). Replace `path` with `include` inside `clusters/deployment.yaml`:

```yaml
# content of clusters/deployment.yaml
deployments:
  - include: {{ target.name }}
```

Now, create a `deployment.yaml` inside `clusers/demo-1`:

```yaml
# content of clusters/demo-1/deployment.yaml
vars:
  - values:
      cluster:
        name: demo-1
        replicas: 1

deployments:
  - path: ../../templates/cluster
  - barrier: true
  - path: ../../templates/workers
    vars:
      - values:
          workers:
            name: md-0
            replicas: 1
  - path: ../../templates/workers
    vars:
      - values:
          workers:
            name: md-1
            replicas: 2
```

The above sub-deployment defines some global configuration (e.g. `cluster.name`) and includes the two previously prepared Kustomize deployments. The cluster level configuration is loaded on sub-deployment level so that all items in `deployments` have access to the configuration found there. The worker specific configuration is specified in-line as part of the deployment item itself. This way, each workers item can have its own configuration (e.g. own name and replicas), which is also demonstrated here by introducing a new node pool.

You'll also find a [barrier]({{% ref "docs/kluctl/deployments/deployment-yml#barriers" %}}) in the list of deployment items. This barrier ensures that Kluctl does not continue deploying worker resources before the cluster resources are applied already.

## Deploying the refactored workload cluster

Simply re-run the deploy command:

```bash
$ kluctl deploy -t demo-1
✓ Loading kluctl project
✓ Initializing k8s client
✓ Rendering templates
✓ Rendering Helm Charts
✓ Building kustomize objects
✓ Postprocessing objects
✓ Writing rendered objects
✓ Getting remote objects by discriminator
✓ Getting namespaces
✓ ../../templates/workers: Applied 3 objects.
✓ ../../templates/cluster: Applied 5 objects.

Changed objects:
  Namespace/cluster-demo-1
  cluster-demo-1/KubeadmConfigTemplate/demo-1-md-0
  cluster-demo-1/Cluster/demo-1
  cluster-demo-1/MachineDeployment/demo-1-md-0
  cluster-demo-1/KubeadmControlPlane/demo-1-control-plane
  cluster-demo-1/DockerCluster/demo-1
  cluster-demo-1/DockerMachineTemplate/demo-1-control-plane
  cluster-demo-1/DockerMachineTemplate/demo-1-md-0

Diff for object Namespace/cluster-demo-1
+-------------------------------------------------------+------------------------------------------+
| Path                                                  | Diff                                     |
+-------------------------------------------------------+------------------------------------------+
| metadata.annotations["kluctl.io/deployment-item-dir"] | -1-initial/clusters/demo-1               |
|                                                       | +2-templating/templates/cluster          |
+-------------------------------------------------------+------------------------------------------+
| metadata.labels["kluctl.io/tag-0"]                    | -clusters                                |
|                                                       | +demo-1                                  |
+-------------------------------------------------------+------------------------------------------+
...
```

You'll see a lot of changes in regard to [tags]({{% ref "docs/kluctl/deployments/tags" %}}) and the `kluctl.io/deployment-item-dir` annotation. These are happening due to the movement of manifests and can be ignored for this tutorial. Simply confirm and let it deploy it.

You should also see that the new workers are being created. You could now try to experiment a little bit by adding more workers or removing old ones. Kluctl will always support you by showing what is new and what got orphaned, allowing you to prune these either via `kluctl prune -t demo-1` or via `kluctl deploy -t demo-1 --prune`.

## Adding more clusters

Adding more clusters is hopefully self-explanatory at this point. It's basically just copying the `demo-1` directory, changing the cluster name in `deployment.yaml` and adding a new target in `.kluctl.yaml`.

## Introducing GitOps

If you prefer to manage your workload clusters via GitOps, the same Kluctl project can be re-used via a simple [KluctlDeployment]({{% ref "docs/gitops/spec/v1beta1/kluctldeployment" %}}) pointing to your Git repository. We won't go into more detail about GitOps here, but feel free to read the documentation and try it on your own. Moving to GitOps doesn't mean that you have to do a full buy-in, as you'll always be able to mix non-GitOp related workflows with GitOps workflows. For example, a `kluctl diff` / `kluctl gitops diff` can always be used even if the same deployment is already managed via GitOps.

## Kluctl vs. ClusterClass

You might ask why one would use Kluctl instead of simply relying on [ClusterClass](https://cluster-api.sigs.k8s.io/tasks/experimental-features/cluster-class/), which is a cluster-api native way of achieving reusability. There are multiple reasons why I believe that Kluctl is a good alternative to ClusterClass, let's go through a few of them.

#### Generic solution

Kluctl is a very generic solution for templated deployments. This means, you can implement a lot of different ways and scenarios that meet different needs. If you already use Kluctl somewhere else, or consider using it somewhere else, you'll easily get used to managing Cluster API via Kluctl. With ClusterClass, you have to learn a new and very Cluster API specific way of templating.

I also believe that it's very likely that you will end up using at least some additional tool on top of the Cluster API manifests, simply because plain `kubectl apply -f ...` is not the best way to do it. Classically, this would be Kustomize or Helm. If GitOps is desired, you might also end up using Flux or ArgoCD. So, if this additional layer of tooling is already required, why not give Kluctl a try and while at it, completely avoid uses of ClusterClass with it.

#### Not limited to Cluster API resources

With ClusterClass, you can only glue together Cluster API related resources. A cluster might however need much more, for example an instance of [Cluster Autoscaler](https://github.com/kubernetes/autoscaler/tree/master/cluster-autoscaler). With ClusterClass, the only option you have is to use a [ClusterResourceSet](https://cluster-api.sigs.k8s.io/tasks/experimental-features/cluster-resource-set) that deploys plain manifests to the workload cluster. These CRSs are however not templated, which will limit you quite a bit in what can be achieved. Also, you must use plain manifests and can't use Helm Charts, which means that the burden of keeping manifests up-to-date is on you. Also, CRSs only allow to deploy additional resource to the workload cluster, but not into the management cluster itself.

With Kluctl, you can use whatever resources you want for the cluster and/or worker templates. Adding Cluster Autoscaler becomes as easy as adding a [Helm Chart]({{% ref "docs/kluctl/deployments/helm" %}}) with proper Helm values (which can also use the `cluster` configuration via templating).

#### Migrations/Modifications to cluster templates

[Changing a ClusterClass](https://cluster-api.sigs.k8s.io/tasks/experimental-features/cluster-class/change-clusterclass) is a risky thing and in my opinion it is crucial to have proper dry-run and diff capabilites. With ClusterClass, this is [supported](https://cluster-api.sigs.k8s.io/clusterctl/commands/alpha-topology-plan#clusterctl-alpha-topology-plan) to some degree but hard to use and [not 100% reliable](https://cluster-api.sigs.k8s.io/clusterctl/commands/alpha-topology-plan#limitations-server-side-apply). With Kluctl, testing changes becomes as easy as changing something and then running `kluctl diff -t demo-1`.

## Wrapping it up

If you want to try out the results of this tutorial without copy-pasing all the YAML, simply clone https://github.com/kluctl/cluster-api-demo and follow the instructions in the README.md.

For a more generic explanation of what Kluctl can do, watch [this live demo](https://www.youtube.com/watch?v=fJgLOyEHmN8) at the [Rawkode Academy](https://www.youtube.com/@RawkodeAcademy) YouTube channel. The documentation at {{% ref "docs" %}} is also worthwhile to read.

You can also join the projects #kluctl channel in the [CNCF Slack](https://cloud-native.slack.com/) and get in contact with existing users and maintainers.