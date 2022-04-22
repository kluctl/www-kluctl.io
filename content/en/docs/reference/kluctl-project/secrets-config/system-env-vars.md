---
title: "systemEnvVars"
linkTitle: "systemEnvVars"
weight: 2
description: >
  Loads secrets from system environment variables.
---

Load secrets from environment variables. Children of `systemEnvVars` can be arbitrary yaml, e.g. dictionaries or lists.
The leaf values are used to get a value from the system environment.

Example:
```yaml
secretsConfig:
  secretSets:
    - name: prod
      sources:
        - systemEnvVars:
            var1: ENV_VAR_NAME1
            someDict:
              var2: ENV_VAR_NAME2
            someList:
              - var3: ENV_VAR_NAME3
```

The above example will make 3 secret variables available: `secrets.var1`, `secrets.someDict.var2` and
`secrets.someList[0].var3`, each having the values of the environment variables specified by the leaf values.
