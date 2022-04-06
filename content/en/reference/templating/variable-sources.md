---
title: "Variable Sources"
linkTitle: "Variable Sources"
weight: 2
description: >
  Available variable sources.
---

There are multiple places in deployment projects (deployment.yml) where additional variables can be loaded into
future Jinja2 contexts.

The first place where vars can be specified is the deployment root, as documented [here]({{< ref "reference/deployments/deployment-yml#vars-deployment-project" >}}).
These vars are visible for all deployments inside the deployment project, including sub-deployments from includes.

The second place to specify variables is in the deployment items, as documented [here]({{< ref "reference/deployments/deployment-yml#vars-deployment-item" >}}).

The variables loaded for each entry in `vars` are not available inside the `deployment.yml` file itself.
However, each entry in `vars` can use all variables defined before that specific entry is processed. Consider the
following example.

```yaml
vars:
- file: vars1.yml
- file: vars2.yml
```

`vars2.yml` can now use variables that are defined in `vars1.yml`. At all times, variables defined by
parents of the current sub-deployment project can be used in the current vars file.

Different types of vars entries are possible:

### file
This loads variables from a yaml file. Assume the following yaml file with the name `vars1.yml`:
```yaml
my_vars:
  a: 1
  b: "b"
  c:
    - l1
    - l2
```

This file can be loaded via:

```yaml
vars:
  - file: vars1.yml
```

After which all included deployments and sub-deployments can use the jinja2 variables from `vars1.yml`.

### values
An inline definition of variables. Example:

```yaml
vars:
  - values:
      a: 1
      b: c
```

These variables can then be used in all deployments and sub-deployments.

### clusterConfigMap
Loads a configmap from the target's cluster and loads the specified key's value as a yaml file into the jinja2 variables
context.

Assume the following configmap to be deployed to the target cluster:
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-vars
  namespace: my-namespace
data:
  vars: |
    a: 1
    b: "b"
    c:
      - l1
      - l2
```

This configmap can be loaded via:

```yaml
vars:
  - clusterConfigMap:
      name: my-vars
      namespace: my-namespace
      key: vars
```

It assumes that the configmap is already deployed before the kluctl deployment happens. This might for example be
useful to store meta information about the cluster itself and then make it available to kluctl deployments.

### clusterSecret
Same as clusterConfigMap, but for secrets.
