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
      sources:
        - ...
...
```

Each `secretSets` entry has the following fields.

### name
This field specifies the name of the secret set. The name can be used in targets to refer to this secret set.

### sources
This field specifies a list of secret sources. Check the sub-sections for the supported secret sources.
