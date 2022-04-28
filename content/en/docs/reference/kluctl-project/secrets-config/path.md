---
title: "path"
linkTitle: "path"
weight: 1
description: >
  Loads secrets from a local file.
---

A simple local file based source. The path must be relative and multiple places are tried to find the file:

1. Relative to the deployment project root
2. The path provided via [--secrets-dir]({{< ref "docs/reference/commands/seal" >}})

The referenced file must be a valid [secrets file]({{< ref "docs/reference/kluctl-project/secrets-config#format-of-secrets-files" >}}).

Example:
```yaml
secretsConfig:
  secretSets:
    - name: prod
      sources:
        - path: .secrets-non-prod.yml
```
