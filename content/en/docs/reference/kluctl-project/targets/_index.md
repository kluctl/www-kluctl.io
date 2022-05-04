---
title: "targets"
linkTitle: "targets"
weight: 4
description: >
  Required, defines targets for this kluctl project.
---

Specifies a list of targets for which commands can be invoked. A target puts together environment/target specific
configuration and the target cluster. Multiple targets can exist which target the same cluster but with differing
configuration (via `args`). Target entries also specifies which secrets to use while [sealing]({{< ref "docs/reference/sealed-secrets" >}}).

Each value found in the target definition is rendered with a simple Jinja2 context that only contains the target itself
and cluster configuration. The rendering process is retried 10 times until it finally succeeds, allowing you to reference
the target itself in complex ways. This is especially useful when using [dynamic targets]({{< ref "./dynamic-targets" >}}).

Target entries have the following form:
```yaml
targets:
...
  - name: <target_name>
    cluster: <cluster_name>
    args:
      arg1: <value1>
      arg2: <value2>
      ...
    dynamicArgs:
      - name: <arg_name>
        pattern: <regex_pattern>
      ...
    sealingConfig:
      secretSets:
        - <name_of_secrets_set>
...
```

The following fields are allowed per target:

## name
This field specifies the name of the target. The name must be unique. It is referred in all commands via the
[-t]({{< ref "docs/reference/commands/common-arguments" >}}) option.

## cluster
This field specifies the name of the target cluster. The cluster must exist in the [cluster configuration]({{< ref "docs/reference/cluster-configs" >}})
specified via [clusters]({{< ref "../external-projects#clusters" >}}).

## args
This fields specifies a map of arguments to be passed to the deployment project when it is rendered. Allowed argument names
are configured via [deployment args]({{< ref "docs/reference/deployments/deployment-yml#args" >}})

## dynamicArgs
This field specifies a list of CLI arguments that can be passed to kluctl when performing any commands on the target. These
arguments are passed with `-a arg_name=arg_value` when for example calling `kluctl deploy -t target_name`.

Each entry has the following fields:

### name
The name of the argument.

### pattern
This field is optional and specifies a regex pattern that the arguments value must match.

## sealingConfig
This field configures how sealing is performed when the [seal command] ({{< ref "docs/reference/commands/seal" >}}) is invoked for this target.
It has the following form:

```yaml
targets:
...
- name: <target_name>
  ...
  sealingConfig:
    dynamicSealing: <true_or_false>
    args:
      arg1: <override_for_arg1>
    secretSets:
      - <name_of_secrets_set>
```

### dynamicSealing
This field specifies weather sealing should happen per [dynamic target]({{< ref "./dynamic-targets" >}}) or only once. This
field is optional and defaults to `true`.

### args
This field allows adding extra arguments to the target args. These are only used while sealing and may override
arguments which are already configured for the target.

### secretSets
This field specifies a list of secret set names, which all must exist in the [secretsConfig]({{< ref "../secrets-config" >}}).
