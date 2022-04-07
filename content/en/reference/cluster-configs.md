---
title: "Cluster Configs"
linkTitle: "Cluster Configs"
weight: 2
description: >
    Cluster configurations.
---

Cluster configurations are a set of yaml files, located in the clusters configuration directory. Each cluster
configuration defines the name of the cluster, which kubectl context belongs to it and which arbitrary configuration
is attached to it.

[Targets]({{< ref "reference/kluctl-project/targets" >}}) refer to these configurations by name. Multiple targets
can reference the same cluster.

## Naming of cluster configuration files

The name of the cluster configuration file must match the name of the cluster, including the yaml file extension. For
example, the cluster configuration for the cluster `test.example.com` must be placed into the file `test.example.com.yml`.

## Location of clusters configurations

The clusters configuration  directory is by default `<kluctl-project-dir>/clusters`. There are 3 ways to override 
the default cluster configurations directory:

1. Via [clusters]({{< ref "reference/kluctl-project/external-projects#clusters" >}}) in `.kluctl.yml`. This is the preferred way
and allows to put the cluster configurations into a separate git repository, which in turn allows you to reuse cluster
configurations for multiple kluctl projects.
2. Via `--local-clusters` when invoking CLI commands. This is useful when you need to override the settings from `.kluctl.yml` due to locally uncommitted changes inside the cluster configurations.
3. Via the environment variable `KLUCTL_LOCAL_CLUSTERS`.

## Minimal cluster config

Let's assume you want to add a cluster with the name "test.example.com". You'd create the file
`<kluctl-project-dir>/clusters/test.example.com.yml` and fill it with the following minimal configuration:

```yaml
cluster:
  name: test.example.com
  context: test.example.com
```

The `context` refers to the kubeconfig context that is later used to connect to the cluster. This means, that you must
have that same cluster configured in your kubeconfig, referred by the given context name. The name and context do not
have to match, but it is recommended to use the same name.

## Using cluster configuration in templates

In every place that allows to use [templates]({{< ref "reference/templating" >}}), the target's cluster
configuration is available in the form of template variables. This means, that you for example can
use `{{ cluster.name }}` to get the cluster name into one of your deployment configurations/resources.

## Custom cluster configuration

The configuration also allows adding as much custom configuration as you need below the `cluster` dictionary.
For example:

```yaml
cluster:
  name: test.example.com
  context: test.example.com
  ingress_config:
    replicas: 1
    external_auth:
      url: https://example.com/auth
```

This is then also available in all templates (e.g. `{{ cluster.ingress_config.replicas }}`). In the above example,
it would allow to configure a test cluster differently from the production cluster when it comes to ingress controller
configuration.
