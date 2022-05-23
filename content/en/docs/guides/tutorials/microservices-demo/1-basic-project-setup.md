---
title: "1. Basic Project Setup"
linkTitle: "1. Basic Project Setup"
weight: 1
---

## Introduction

This is the first tutorial in a series of tutorials around the [GCP Microservices Demo](https://github.com/GoogleCloudPlatform/microservices-demo)
and the use of kluctl to deploy and manage the demo.

We will start with a simple kluctl project setup (this tutorial) and then advance to a multi-environment and
multi-cluster setup (upcoming tutorial). Afterwards, we will also show how daily business (updates, house keeping, ...)
with such a deployment would look like.

## GCP Microservices Demo
From the README.md of [GCP Microservices Demo](https://github.com/GoogleCloudPlatform/microservices-demo):

> Online Boutique is a cloud-native microservices demo application. Online Boutique consists of a 10-tier microservices application. The application is a web-based e-commerce app where users can browse items, add them to the cart, and purchase them.

This demo application seems to be a good example for a more or less typical application seen on Kubernetes. It has multiple
self-developed microservices while also requiring third-party applications/services (e.g. redis) to be deployed and configured properly.

## Ways to deploy the demo
The simplest and most naive way to deploy the demo is by using `kubectl apply` with the provided release manifests:

```bash
$ kubectl apply -f https://raw.githubusercontent.com/GoogleCloudPlatform/microservices-demo/main/release/kubernetes-manifests.yaml
```

This is also what is shown in the README.md of the microservices demo.

The shortcomings of this approach are however easy to spot, and probably no one would ever follow this approach up to
production. As an example, updates to the application and its dependencies will be hard to maintain. Housekeeping
(deleting orphan resources) will also be hard to achieve. At some point in time, when you start deploying the application
multiple times to different clusters and/or different environments, configuration will also become hard to maintain, as
every target might need different configuration. Long story short...without proper tooling, you'll easily run into
painful limitations.

There are multiple solutions available that each solve parts of the limitations and problems. As an example,
[Helm](https://helm.sh) and [Kustomize](https://kustomize.io/) are well known. Introducing these tools will easily
bring you much further, but you will very likely end up with something complicated/complex around these tools to
make it usable in daily business. In the worst case, you'd start using Bash scripts that orchestrate your deployments.

GitOps oriented solutions like [ArgoCD](https://argoproj.github.io/cd/) and [Flux](https://https://fluxcd.io/) are
able to relieve you from parts of the deployment orchestration tasks, but bring in new complexities that need to be
solved as well.

## Deploying with kluctl
In this tutorial, we'll show how the microservices demo can be deployed and managed with kluctl. We will start with a
simple and naive deployment to a local [kind](https://kind.sigs.k8s.io/) cluster. The next tutorial in this series will
then focus on making the deployment multi-environment and multi-cluster capable.

The goal is to make a deployment as simple as typing:

```shell
$ kluctl deploy -t local
```

## Setting up the kluctl project
The first thing you need is an empty project directory and the [`.kluctl.yml`]({{< ref "docs/reference/kluctl-project" >}}) project configuration:

```shell
$ mkdir -p microservices-demo/1-basic-setup
$ cd microservices-demo/1-basic-setup
```

Inside this new directory, create the file `.kluctl.yml` with the following content:
```yaml
targets:
  - name: local
    context: kind-kind
```

This is a very simple example with only a single target, being a local [kind](https://kind.sigs.k8s.io/) cluster.

You might have noticed that the target configuration refers a kubectl context that is not existing yet. It's time to
create a local kind cluster now. To do so, first ensure that you have [kind installed](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
and then run:

```sh
$ kind create cluster
```

After this, you should have a local cluster setup and your kubeconfig prepared with a new context named `kind-kind`.

## Setting up a minimal deployment project
Inside the kluctl project, you will now have to create a minimal [deployment project]({{< ref "docs/reference/deployments" >}}).
The deployment project starts with the root [`deployment.yml`]({{< ref "docs/reference/deployments/deployment-yml" >}}).

The location of this `deployment.yml` is the same as the `.kluctl.yml`. Create the file with following content:
```yaml
deployments:
  - path: redis

commonLabels:
  examples.kluctl.io/deployment-project: "microservices-demo"
```

This minimal deployment project contains two elements:
1. The list of [deployment items]({{< ref "docs/reference/deployments/deployment-yml#deployments" >}}), which currently
only consists of the upcoming redis deployment. The next chapter will explain this deployment.
2. The [commonLabels]({{< ref "docs/reference/deployments/deployment-yml#commonlabels" >}}), which is a map of common
labels and values. These labels are applied to all deployed resources and are later used by kluctl to identify resources
that belong to this kluctl deployment.

## Setting up the redis deployment
As seen in the previous chapter, the root `deployment.yml` refers to a `redis` deployment item. This deployment item must
be located inside the sub-folder `redis` (hence the `path: redis`). kluctl expects each deployment item to be a
[kustomize](https://kustomize.io/) deployment. Such a kustomize deployment can be as simple as a `kustomization.yml` with
a single `resources` entry or a fully fledged kustomize deployment with overlays, generators, and so on.

For our example, first create the sub-directory `redis`:
```sh
$ mkdir redis
```

Then create the file `redis/kustomization.yml` with the following content:
```yaml
resources:
  - deployment.yml
  - service.yml
```

Then create the file `redis/deployment.yml` with the following content:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: redis-cart
spec:
  selector:
    matchLabels:
      app: redis-cart
  template:
    metadata:
      labels:
        app: redis-cart
    spec:
      containers:
      - name: redis
        image: redis:alpine
        ports:
        - containerPort: 6379
        readinessProbe:
          periodSeconds: 5
          tcpSocket:
            port: 6379
        livenessProbe:
          periodSeconds: 5
          tcpSocket:
            port: 6379
        volumeMounts:
        - mountPath: /data
          name: redis-data
        resources:
          limits:
            memory: 256Mi
            cpu: 125m
          requests:
            cpu: 70m
            memory: 200Mi
      volumes:
      - name: redis-data
        emptyDir: {}
```

And the file `redis/service.yml`:
```yaml
apiVersion: v1
kind: Service
metadata:
  name: redis-cart
spec:
  type: ClusterIP
  selector:
    app: redis-cart
  ports:
  - name: redis
    port: 6379
    targetPort: 6379
```

The above files (`deployment.yml` and `service.yml`) are based on the content of [redis.yaml](https://github.com/GoogleCloudPlatform/microservices-demo/blob/main/kubernetes-manifests/redis.yaml)
from the original GCP Microservices Demo.

As you can see, there is nothing special about the contents of these files so far. It's simple and plain Kubernetes and
YAML resources. The full potential of kluctl will become clear later, when we start to use templating inside these files.
Only with the templating, it will become possible to support multi-environment and multi-cluster deployments.

## Setting up the first microservice
Now it's time to setup the first microservice. It is done the same way as we're already setup the redis deployment.

First, create the sub-directory `cartservice` at the same level as you created the `redis` sub-directory. Then create
the following files.

Another `kustomization.yml`
```yaml
resources:
  - deployment.yml
  - service.yml
```

Another `deployment.yml`, with the content found [here](https://github.com/kluctl/kluctl-examples/blob/main/microservices-demo/1-basic-project/cartservice/deployment.yml)

Another `service.yml`, with the content found [here](https://github.com/kluctl/kluctl-examples/blob/main/microservices-demo/1-basic-project/cartservice/service.yml)

Finally add the new deployment item to the root `deployment.yml`

```yaml
...
deployments:
  ...
  # add this line
  - path: cartservice
...
```

## Setting up all other microservices
The [GCP Microservices Demo](https://github.com/GoogleCloudPlatform/microservices-demo) is composed of multiple other
services, which can be setup the same way as the microservice shown before. You can do this by yourself, or alternatively
switch to the completed example found [here](https://github.com/kluctl/kluctl-examples/blob/main/microservices-demo/1-basic-project).

From now on, we will assume that all microservices have been added (or that you switched to the example project).

## Deploy it!
We now have a minimal kluctl project with two simple kustomize deployments. It's time to deploy it. From inside the
kluct project directory, call:

```sh
$ kluctl deploy -t local
INFO[0000] Rendering templates and Helm charts          
INFO[0000] Building kustomize objects                   
Do you really want to deploy to the context/cluster kind-kind? (y/N) y
INFO[0001] Getting remote objects by commonLabels       
INFO[0001] Getting 24 additional remote objects         
INFO[0001] Running server-side apply for all objects    
INFO[0001] shippingservice: Applying 2 objects          
INFO[0001] paymentservice: Applying 2 objects           
INFO[0001] currencyservice: Applying 2 objects          
INFO[0001] frontend: Applying 3 objects                 
INFO[0001] loadgenerator: Applying 1 objects            
INFO[0001] recommendationservice: Applying 2 objects    
INFO[0001] productcatalogservice: Applying 2 objects    
INFO[0001] adservice: Applying 2 objects                
INFO[0001] cartservice: Applying 2 objects              
INFO[0001] emailservice: Applying 2 objects             
INFO[0001] checkoutservice: Applying 2 objects          
INFO[0001] redis: Applying 2 objects                    

New objects:
  default/Deployment/adservice
  default/Deployment/cartservice
  default/Deployment/checkoutservice
  default/Deployment/currencyservice
  default/Deployment/emailservice
  default/Deployment/frontend
  default/Deployment/loadgenerator
  default/Deployment/paymentservice
  default/Deployment/productcatalogservice
  default/Deployment/recommendationservice
  default/Deployment/redis-cart
  default/Deployment/shippingservice
  default/Service/adservice
  default/Service/cartservice
  default/Service/checkoutservice
  default/Service/currencyservice
  default/Service/emailservice
  default/Service/frontend
  default/Service/frontend-external
  default/Service/paymentservice
  default/Service/productcatalogservice
  default/Service/recommendationservice
  default/Service/redis-cart
  default/Service/shippingservice
```

The `-t local` selects the `local` target which was previously defined in the `.kluctl.yml`. Right now we only have this
one target, but we will add more targets in upcoming tutorials from this series.

Answer with `y` to the question if you really want to deploy. The command will output what is happening and then show
what has been changed on the target.

## Playing around
You have now deployed redis and the cartservice microservice. You can now start to play around with some other kluctl
commands. For example, try to change something inside `cartservice.yml` (e.g. set terminationGracePeriodSeconds to 10)
and then run `kluctl diff -t local`:

```sh
$ kluctl diff -t local
INFO[0000] Rendering templates and Helm charts          
...

Changed objects:
  default/Deployment/cartservice

Diff for object default/Deployment/cartservice
+--------------------------------------------------+---------------------------+
| Path                                             | Diff                      |
+--------------------------------------------------+---------------------------+
| spec.template.spec.terminationGracePeriodSeconds | -5                        |
|                                                  | +10                       |
+--------------------------------------------------+---------------------------+
```

As you can see, kluctl now shows you what will happen. If you'd now perform a `kluctl deploy -t local`, kluctl would
output what has happened (which would be the same as in the diff as long as you don't change anything else).

If you try to remove (or at least comment out) a microservice, e.g. the cartservice and then run
`kluctl diff -t local` again, you will get:

```sh
$ kluctl diff -t local
INFO[0000] Rendering templates and Helm charts          
...

Changed objects:
  default/Deployment/cartservice

Diff for object default/Deployment/cartservice
+--------------------------------------------------+---------------------------+
| Path                                             | Diff                      |
+--------------------------------------------------+---------------------------+
| spec.template.spec.terminationGracePeriodSeconds | -5                        |
|                                                  | +10                       |
+--------------------------------------------------+---------------------------+

Orphan objects:
  default/Service/cartservice
  default/Deployment/cartservice
 ```

As you can see, the resources belonging cartservice are listed as "Orphan objects" now, meaning that these are not found locally anymore.
A `kluctl prune -t local` would then give:

```sh
$ kluctl prune -t local
INFO[0000] Rendering templates and Helm charts          
...
Do you really want to delete 2 objects? (y/N) y

Deleted objects:
  default/Service/cartservice
  default/Deployment/cartservice
```

## How to continue
The result of this tutorial is a naive version of the microservices demo deployment. There are a few things that you
would solve differently in the real world, e.g. use Helm Charts for things like redis instead of proving self-crafted
manifests. The next tutorials in this series will focus on a few improvements and refactorings that will make this
kluctl project more "realistic" and more useful. They will also introduce concepts like multi-environment and multi-cluster
deployments.
