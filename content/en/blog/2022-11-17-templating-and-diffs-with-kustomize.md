
---
title: "Templating in Kustomize with Kluctl style deploy and diff"
linkTitle: "Templating in Kustomize with Kluctl style deploy and diff"
slug: templating-and-diffs-with-kustomize
date: 2022-11-17
author: Alexander Block (@codablock)
images:
  - "/images/blog/cover/Templating in Kustomize.png"
---

![image](/images/blog/templating-and-diffs-with-kustomize.png)

[Kustomize](https://kustomize.io/) is currently one of the most used tools to organise Kubernetes manifests and the
resulting deployments. As of the Kustomize website, "Kustomize introduces a template-free way to customize application
configuration that simplifies the use of off-the-shelf applications."

This says it very clear: Kustomize is template-free. The reasoning why Kustomize does not leverage templating and also
will never do so is very reasonable and easy to understand. It tries to avoid the potential overload and complexity that
comes with templating.

In my opinion however, the feared complexity is only a problem if one tries to suit the needs of everyone. This is an
issue that popular [Helm Charts](https://helm.sh) do have for example. If you try to make everyone happy, your templates
must make everything configurable, which eventually leads to Kubernetes manifests having more templating code than
the actual YAML.

Kustomize in that regard, has the advantage that the re-used manifests do not have to take external customization into
account. You, as the "customizer", can decide what needs to be configurable and can achieve this with overlays and
patches.

Looking at Kustomize from that perspective, it is of course very reasonable to keep it fully template-free. There is
however also another perspective that you might want to consider.

## What if I don't care about the needs of others?

That sounds a bit selfish, doesn't it? :) But let's be honest to our self, many times you just need to create a
deployment that suits your own needs right now. This deployment might re-use other components which it needs to
customize, but it doesn't need to be customized by someone else.

Thus, the level of customization that you need to implement is minimal. Maybe you just want to allow changing the
target namespace or some replica count depending on the target environment. Maybe you want to have some components
enabled in one environment and other components disabled in other environments. If you think more about it, you might
also realise that it's not about "customization" anymore but actually about "configuration".

In that case, templating does not bring the risks that come with components that are meant to be re-used and customized.
This is because it is very clear what level of configuration is required and thus the use of templating can be reduced
to exactly that. In the end, you'll only need a few places with something like `{{ my_service.replicas }}` and maybe some
conditional blocks with `{% if my_service.enabled %}...{% endif %}`.

On the other hand, using plain Kustomize to achieve the same result can easily become unnecessarily complex. For example, the
use of bases, overlays and patches tends to create project structures that are hard to grasp when they grow. It can also
force you to change your project structure in "unnatural" (at least that is how it feels for me) ways, because you have
to adapt to the way overlays work. Templating would allow much simpler solutions in the above case.

But...Kustomize doesn't support templating, right?

## Bringing templating to Kustomize

[Kluctl](https://kluctl.io) builds a large set of its features and promised advantages on top of
[templating](https://kluctl.io/docs/kluctl/latest/reference/templating/). The
[Kustomize integration](https://kluctl.io/docs/kluctl/latest/deployments/kustomize/) also allows templating in all
involved resources, including the `kustomization.yaml` itself and all referenced manifests.

Configuration can be provided in different ways:
1. Via [CLI arguments](https://kluctl.io/docs/kluctl/latest/reference/commands/common-arguments/#project-arguments), e.g. `--arg` or `--args-from-file`.
2. Via [Targets](https://kluctl.io/docs/kluctl/latest/kluctl-project/targets/#args), meaning that you can define named targets with fixed args.
3. Via [vars](https://kluctl.io/docs/kluctl/latest/deployments/deployment-yml/#vars-deployment-project) in Kluctl deployments.
4. Via [Environment variables](https://kluctl.io/docs/kluctl/latest/reference/commands/environment-variables/#environment-variables-as-arguments) (through `KLUCTL_ARG_XXX`).

In this blog post, we'll focus on the first option for simplicity. The second and third options are much more
powerful, but require more boilerplate to set up a [Kluctl project](https://kluctl.io/docs/kluctl/latest/kluctl-project/)
and [Kluctl deployments](https://kluctl.io/docs/kluctl/latest/deployments/deployment-yml). The first option also works with
plain Kustomize deployments, which is what I'm going to demonstrate.

Whatever option is used, all "args" are then available in every place by simply using
[Jinja2 variable expressions](https://jinja.palletsprojects.com/en/3.1.x/templates/#variables), e.g. `{{ args.my_arg }}`.

## A simple example

We will use the [podtato-head](https://github.com/podtato-head/podtato-head) project, specifically the
[Kustomize delivery scenario](https://github.com/podtato-head/podtato-head/tree/main/delivery/kustomize) as an example
. But first, we'll need a test cluster. I suggest to simply use [Kind](https://kind.sigs.k8s.io/) and create a fresh
cluster:

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

Ensure that you are on the correct Kubernetes context by calling `kubectl config current-context` and verify that it
points to `kind-kind`. Without using [Kluctl Targets](https://kluctl.io/docs/kluctl/latest/kluctl-project/targets/) with
fixed contexts, Kluctl behaves the same as any other tool in the Kubernetes space when it comes to the current context,
meaning that you have to watch out to not deploy to prod by accident! :)

Then, clone the example project and enter the delivery directory:

```shell
$ git clone https://github.com/podtato-head/podtato-head.git
$ cd podtato-head/delivery/kustomize/base
$ ls -lah
total 120
drwxr-xr-x  17 user  wheel   544B Nov 16 16:28 .
drwxr-xr-x   6 user  wheel   192B Nov 16 16:28 ..
-rw-r--r--   1 user  wheel   558B Nov 16 16:28 configmap-discovery.yaml
-rw-r--r--   1 user  wheel   1.5K Nov 16 16:28 deployment-entry.yaml
-rw-r--r--   1 user  wheel   1.1K Nov 16 16:28 deployment-hat.yaml
-rw-r--r--   1 user  wheel   1.1K Nov 16 16:28 deployment-left-arm.yaml
-rw-r--r--   1 user  wheel   1.1K Nov 16 16:28 deployment-left-leg.yaml
-rw-r--r--   1 user  wheel   1.1K Nov 16 16:28 deployment-right-arm.yaml
-rw-r--r--   1 user  wheel   1.1K Nov 16 16:28 deployment-right-leg.yaml
-rw-r--r--   1 user  wheel   474B Nov 16 16:28 kustomization.yaml
-rw-r--r--   1 user  wheel   447B Nov 16 16:28 service-entry.yaml
-rw-r--r--   1 user  wheel   438B Nov 16 16:28 service-hat.yaml
-rw-r--r--   1 user  wheel   454B Nov 16 16:28 service-left-arm.yaml
-rw-r--r--   1 user  wheel   453B Nov 16 16:28 service-left-leg.yaml
-rw-r--r--   1 user  wheel   456B Nov 16 16:28 service-right-arm.yaml
-rw-r--r--   1 user  wheel   456B Nov 16 16:28 service-right-leg.yaml
-rw-r--r--   1 user  wheel   281B Nov 16 16:28 serviceaccount.yaml
```

As you can see, this is a simple Kustomize deployment, not using any bases or overlays. Let's start using Kluctl by
doing a vanilla deployment first:

```shell
$ kluctl deploy
✓ Loading kluctl project
✓ Initializing k8s client
✓ Rendering templates
✓ Rendering Helm Charts
✓ Building kustomize objects
✓ Postprocessing objects
✓ Getting namespaces
✓ .: Applied 14 objects.
✓ Running server-side apply for all objects

New objects:
  default/ConfigMap/podtato-head-service-discovery
  default/Deployment/podtato-head-entry
  default/Deployment/podtato-head-hat
  default/Deployment/podtato-head-left-arm
  default/Deployment/podtato-head-left-leg
  default/Deployment/podtato-head-right-arm
  default/Deployment/podtato-head-right-leg
  default/Service/podtato-head-entry
  default/Service/podtato-head-hat
  default/Service/podtato-head-left-arm
  default/Service/podtato-head-left-leg
  default/Service/podtato-head-right-arm
  default/Service/podtato-head-right-leg
  default/ServiceAccount/podtato-head
? The diff succeeded, do you want to proceed? (y/N) y
✓ .: Applied 14 objects.
✓ Running server-side apply for all objects

New objects:
  default/ConfigMap/podtato-head-service-discovery
  default/Deployment/podtato-head-entry
  default/Deployment/podtato-head-hat
  default/Deployment/podtato-head-left-arm
  default/Deployment/podtato-head-left-leg
  default/Deployment/podtato-head-right-arm
  default/Deployment/podtato-head-right-leg
  default/Service/podtato-head-entry
  default/Service/podtato-head-hat
  default/Service/podtato-head-left-arm
  default/Service/podtato-head-left-leg
  default/Service/podtato-head-right-arm
  default/Service/podtato-head-right-leg
  default/ServiceAccount/podtato-head
```

`kluctl` shows you the diff it will apply, asks for confirmation, applies the changes and then shows you the applied result. Verify that it got
deployed:

```shell
$ kubect get pod
NAME                                      READY   STATUS    RESTARTS   AGE
podtato-head-entry-7dfd8cdd6d-6mtxd       1/1     Running   0          95s
podtato-head-hat-6bcbf5f957-mfc6r         1/1     Running   0          95s
podtato-head-left-arm-7d9db78544-689tx    1/1     Running   0          95s
podtato-head-left-leg-59f45ffc4-grcjc     1/1     Running   0          95s
podtato-head-right-arm-5444b48b85-427w7   1/1     Running   0          95s
podtato-head-right-leg-f68df999f-g27nz    1/1     Running   0          95s
```

There are multiple things that you might have noticed already:
1. Kustomize does not handle deployments, it just "builds" them and then let's you handle the actual deployment via `kubectl apply -f`.
Kluctl also handles the actual deployment for you. The advantages of this will be clear in a few minutes.
2. Kluctl showed a diff (very simple one in this case, just new objects), then asked for confirmation and then showed what it did (identical to what the diff showed).
The power of the diff feature will become much clearer in a few minutes.

## Introducing some templating

Now let's introduce some templating into the example deployment. For example, let's edit `deployment-entry.yaml` and
change the replicas field to:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podtato-head-entry
...
spec:
  replicas: {{ args.entry_replicas }}
...
```

If you'd try to re-deploy this now, it will throw an error about `args.entry_replicas` being undefined. So let's call
Kluctl with the now required arg:

```shell
$ kluctl deploy -a entry_replicas=3
...

Changed objects:
  default/Deployment/podtato-head-entry

Diff for object default/Deployment/podtato-head-entry
+---------------+--------------------------------------------------------------+
| Path          | Diff                                                         |
+---------------+--------------------------------------------------------------+
| spec.replicas | -1                                                           |
|               | +3                                                           |
+---------------+--------------------------------------------------------------+
? The diff succeeded, do you want to proceed? (y/N)
```

You'll notice that Kluctl again stops and asks for confirmation. But this time, it will actually show you some
meaningful diff. It allows you to verify that Kluctl will apply the intended changes. The diff
that you see is NOT a simple file based diff, but a diff after performing a full-blown server-side apply in dry-run
mode. This means, what you see is what you'll get, no surprises in-between!

Let's confirm with `y`:

```shell
? The diff succeeded, do you want to proceed? (y/N) y
✓ .: Applied 14 objects.
✓ Running server-side apply for all objects

Changed objects:
  default/Deployment/podtato-head-entry

Diff for object default/Deployment/podtato-head-entry
+---------------+--------------------------------------------------------------+
| Path          | Diff                                                         |
+---------------+--------------------------------------------------------------+
| spec.replicas | -1                                                           |
|               | +3                                                           |
+---------------+--------------------------------------------------------------+
```

The actual deployment is performed and the result printed to the user. The result should always be identical to the diff you saw before.

## Let's make something conditional

Now let's make the "hat" of the podtate-head optional. However, as you have previously deployed the project already,
you'll need to delete the hat deployment manually:

```shell
$ kubectl delete deployment.apps/podtato-head-hat service/podtato-head-hat
```

Edit `kustomization.yaml` and put an if/endif around the hat resource:

```yaml
...
resources:
- configmap-discovery.yaml
{% if args.hat_enabled | default(true) %}
- deployment-hat.yaml
- service-hat.yaml # also remove the original entry from the bottom of the file
{% endif %}
- deployment-left-arm.yaml
...
```

You can now deploy with Kluctl while having the hat disabled:

```shell
$ kluctl deploy -a entry_replicas=3 -a hat_enabled=false
```

If you deploy with a hat and then with `hat_enabled=false`, Kluctl will not delete/prune the previously deployed hat.
If you want to have pruning support, you must create a [Kluctl deployment](https://kluctl.io/docs/kluctl/latest/deployments/)
with `commonLabels` enabled, so that Kluctl knows how to identify related objects.

## Using vars files instead of arguments

`args` can also be passed via vars files, which are arbitrary structured YAML files. This is a comparable to helm value
files. Consider the following examples.

test-args.yaml:
```yaml
entry_replicas: 2
hat_enabled: false
```

prod-args.yaml:
```yaml
entry_replicas: 3
```

These values can be used with `--args-from-file`:

```shell
$ kluctl diff --args-from-file=test-args.yaml
$ kluctl diff --args-from-file=prod-args.yaml
```

Based on that, you can easily implement multi-environment deployments. This is however still a poor-mans solution to
multi-environment deployments, with the use of [Kluctl projects](https://kluctl.io/docs/kluctl/latest/kluctl-project/)
being the better solution. Said Kluctl projects allow you to define named [targets](https://kluctl.io/docs/kluctl/latest/kluctl-project/targets/)
which are fixed in their configuration, so that you only have to invoke `kluctl deploy -t test` without needing to
know what the internal details are.

## What's next?

This article has shown how Kluctl can be used on simple/plain Kustomize deployments. The next thing you should consider
is using [Kluctl projects](https://kluctl.io/docs/kluctl/latest/kluctl-project/) and
[Kluctl deployments](https://kluctl.io/docs/kluctl/latest/deployments/) around your Kustomize deployments. It will allow you
to have much more flexible and powerful ways of configuration management. It will also allow you to use the GitOps style
[flux-kluctl-controller](https://github.com/kluctl/flux-kluctl-controller).

Upcoming blog posts will show why the [Helm Integration](https://kluctl.io/docs/kluctl/latest/deployments/helm/) is a good
thing with many advantages and also describe why one would choose Kluctl over plain Helm.
