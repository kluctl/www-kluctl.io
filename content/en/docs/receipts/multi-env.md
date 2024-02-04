---
title: "Deploying multiple times"
linkTitle: "Deploying multiple times"
weight: 20
---

This receipt will guide you on how to deploy the same deployment multiple times to the same (via namespaces) or
different clusters.

## Use specific targets

The easiest way to achieve this is to define [targets]({{< ref "docs/kluctl/kluctl-project/targets" >}}) in
your `.kluctl.yaml`. Each target should then use [args]({{< ref "docs/kluctl/kluctl-project/targets#args" >}}) to define
a small set configuration values for the specific target.

Each target should relate to the target environment and/or cluster that it needs to be deployed to. For example, one
could be named `prod` while another is named `test`, meaning that you can either deploy to the `prod` or to the `test`
environment. It's also useful to set the [context]({{< ref "docs/kluctl/kluctl-project/targets#context" >}}) field
on each target, so that you can't accidentally deploy the `prod` target to the `test` cluster.

`args` should be minimalistic to avoid bloating up the `.kluctl.yaml`. It should be used as the "entrypoint" into
the actual configuration, which is then loaded from inside the root [deployment.yaml]({{< ref "docs/kluctl/deployments/deployment-yml" >}})
via `vars`. See [advanced configuration]({{< ref "docs/receipts/advanced-configuration" >}}) for details on this.

Example targets definition:

```yaml
targets:
  - name: prod
    context: prod.example.com
    args:
      environment_name: prod
  - name: test
    context: test.example.com
    args:
      environment_name: test

# Warning, this discriminator is only ok if targets are only deployed once per cluster. See next chapter for details.
discriminator: "my-project-{{ target.name }}"

args:
  - name: environment_name
```

Example CLI invocations:
```shell
$ kluctl deploy -t prod
$ kluctl deploy -t test
```

## Use more dynamic targets

As an alternative to very specific targets, you could also define targets that are more dynamic so that each target can
be deployed multiple times, but to different Kubernetes contexts or even namespaces. You can also mix such targets,
for example have one `prod` target that is just like described in the previous chapter, and one `non-prod` target
that can be used to deploy to multiple non-production clusters.

The dynamic targets then need a way so that they can be differentiated. The easiest way is to use different contexts,
which means you deploy it to different clusters. Another way is to introduce `args` that serve to differentiate, e.g.
an arg names `environment_name` which can then be used to deploy the same workloads to different namespaces, add prefixes
to global resources, create unique ingresses, and so on.

If such an argument is introduced, you would then invoke the CLI with the argument being set.

Another thing to take into account is the required uniqueness of [discriminators]({{< ref "docs/kluctl/kluctl-project/targets#discriminator" >}})
to make [delete]({{< ref "docs/kluctl/commands/delete" >}}) and [prune]({{< ref "docs/kluctl/commands/prune" >}}) work
properly. If you miss this crucial part or make a mistake, you might end up deleting resources that were not meant to
be deleted. The uniqueness must be ensured inside the boundaries of individual clusters.

Example targets definition:
```yaml
targets:
  - name: prod
    context: prod.example.com
    args:
      environment_type: prod
      environment_name: prod
  - name: non-prod
    args:
      environment_type: non-prod
      # environment_name must be passed via CLI

# This will ensure that the discriminator is unique, even if the same target is deployed multiple times
discriminator: my-project-{{ target.name }}-{{ args.environment_type }}-{{ args.environment_name }}

# This is a bad example of a discriminator. It will lead to the discriminator being equal for every environment deployed to the same cluster.
# discriminator: "my-project-{{ target.name }}"

args:
  - name: environment_type
  - name: environment_name
```

Example CLI invocations:
```shell
$ kluctl deploy -t prod # deploys to prod context
$ kluctl deploy -t non-prod -a environment_name=test-env1 # deploys to currently active context 
$ kluctl deploy -t non-prod -a environment_name=test-env2 # deploys to currently active context
$ kluctl deploy -t non-prod -a environment_name=test-env3 --context test2.exmaple.com
```

## Too long discriminators

Right now, Kluctl is internally using a single label to store discriminators in Kubernetes. This has some serious
limitations in regard to the length of the discriminators, which is 63 characters. This means, that the discriminator
template shown in the above example can easily lead to errors. This issue will be fixed when https://github.com/kluctl/kluctl/issues/468
is implemented.

Until then, you might need to use some form of shortening, e.g. by using a shortened hash of some string. Example
for this:

```yaml
discriminator: my-project-{{ target.name }}-{{ args.environment_type }}-{{ (args.environment_name | sha256)[:8] }}
```

## Using namespaces and more

So far, we have only shown how to define and use the `targets` feature to implement multiple target environments.
This works out-of-the-box when you target different clusters per target, but will need some additional work when
deploying to the same cluster. In that case, you are required to use different namespaces for each environment.

This can be easily achieved by using the mentioned `environment_name` inside resources. Combined with templating, it can
be used to create dynamic namespaces, prefix resource names and create unique ingresses.

Example project:

```
my-project/
├── .kluctl.yaml
├── deployment.yaml
├── namespaces/
│   └── namespace.yaml
└── apps
    ├── deployment.yaml
    ├── app1/
    │   ├── resource1.yaml
    │   └── resource2.yaml
    └── app2/
        ├── resource1.yaml
        └── resource2.yaml
```

##### .kluctl.yaml

See above.

##### deployment.yaml

```yaml
deployments:
  - path: namespaces
  - barrier: true # ensure namespaces are applied before we continue
  - include: apps
```

###### namespaces/namespace.yaml

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: {{ args.environment_name }}
```

###### apps/deployment.yaml

```yaml
deployments:
  - path: app1
  - path: app2

# This instructs Kluctl to set the specified namespace on all resources, including resources from `app1` and `app2`,
# that do not have a namespace set explicitly.
overrideNamespace: {{ args.environment_name }}
```

###### apps/app1/resource1.yaml

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-cm
  # no namespace needed here, as it is set via the `overrideNamespace` from `apps/deployment.yaml`
data:
  # just an example to show that you can also use the `args` here.
  environment_name: {{ args.environment_name }}
```
