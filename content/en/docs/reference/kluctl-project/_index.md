---
title: "Kluctl project (.kluctl.yaml)"
linkTitle: ".kluctl.yaml"
weight: 1
description: >
    Kluctl project configuration, found in the .kluctl.yaml file.
---

The `.kluctl.yaml` is the central configuration and entry point for your deployments. It defines where the actual
[deployment project]({{< ref "docs/reference/deployments" >}}) is located,
where [sealed secrets]({{< ref "docs/reference/sealed-secrets" >}}) and unencrypted secrets are localed and which targets are available to
invoke [commands]({{< ref "docs/reference/commands" >}}) on.

## Example

An example .kluctl.yaml looks like this:

```yaml
# This is optional. If omitted, the same directory where `.kluctl.yaml` is located will be used as root deployment
# See "External Projects" for details
deployment:
  project:
    url: https://github.com/kluctl/kluctl-example

# This is optional. If omitted, `<baseDirOfKluctlYml>/.sealed-secrets` will be used
# See "External Projects" for details
sealedSecrets:
  project:
    url: https://github.com/kluctl/kluctl-example
    subDir: .sealed-secrets

targets:
  # test cluster, dev env
  - name: dev
    cluster: test.example.com
    args:
      environment_name: dev
    sealingConfig:
      secretSets:
        - non-prod
  # test cluster, test env
  - name: test
    cluster: test.example.com
    args:
      environment_name: test
    sealingConfig:
      secretSets:
        - non-prod
  # prod cluster, prod env
  - name: prod
    cluster: prod.example.com
    args:
      environment_name: prod
    sealingConfig:
      secretSets:
        - prod

# This is only required if you actually need sealed secrets
secretsConfig:
  secretSets:
    - name: prod
      vars:
        # This file should not be part of version control!
        - file: .secrets-prod.yaml
    - name: non-prod
      vars:
        # This file should not be part of version control!
        - file: .secrets-non-prod.yaml
```

## Allowed fields

Please check the sub-sections of this section to see which fields are allowed at the root level of `.kluctl.yaml`.

## Separating kluctl projects and deployment projects

As seen in the `.kluctl.yaml` documentation, deployment projects can exist in other repositories then the kluctl project.
This is a desired pattern in some circumstances, for example when you want to share a single deployment project with
multiple teams that all manage their own clusters. This way each team can have its own minimalistic kluctl project which
points to the deployment project and the teams clusters configuration.

This way secret sources can also differ between teams and sharing can be reduced to a minimum if desired.
