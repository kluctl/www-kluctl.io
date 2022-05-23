---
title: "secretsConfig"
linkTitle: "secretsConfig"
weight: 4
description: >
  Optional, defines where to load secrets from.
---

This configures how secrets are retrieved while sealing. It is basically a list of named secret sets which can be
referenced from targets.

It has the following form:
```yaml
...
secretsConfig:
  secretSets:
    - name: <name>
      vars:
        - ...
...
```

Each `secretSets` entry has the following fields.

### name
This field specifies the name of the secret set. The name can be used in targets to refer to this secret set.

### vars
A list of variables sources. Check the documentation of
[variables sources]{{< ref "docs/reference/templating/variable-sources" >}} for details.

Each variables source must have a root dictionary with the name `secrets` and all the actual secret values
below that dictionary. Every other root key will be ignored.

Example variables file:

```yaml
secrets:
  secret: value1
  nested:
    secret: value2
    list:
      - a
      - b
...
```
