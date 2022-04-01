---
title: "External Projects"
linkTitle: "External Projects"
weight: 1
description: >
  Using external projects 
---

Kluctl will by default use the Kluctl project itself as default location for deployments, cluster configurations and
sealed secrets. These can however be externalized into other Git repositories, Kluctl will then clone/fetch these external
repositories when necessary.

External Projects allow better reuse of cluster configuration and deployments. As an example, you can have a central
cluster configurations repository for all Kluctl projects, making it easier to manage clusters without having to touch
all Kluctl projects.

Deployment can also be reused by having multiple minimal Kluctl projects that simply point to a single deployment project. 

The following configuration is possible in `.kluctl.yml`

## deployment

Specifies the git project where the [deployment project]({{< ref "docs/reference/deployments" >}}) is located. If it is omitted, the base
directory of the `.kluctl.yml` project is used as the deployment project root.

It has the following form:
```yaml
deployment:
  project:
    url: <git-url>
    ref: <tag-or-branch>
    subdir: <subdir>
```

### project.url
Specifies the git clone url of the project.

### project.ref
This field is optional and specifies which tag/branch to use. If omitted, the repositories default branch is used.

### project.subdir
This field is optional and specifies the subdirectory to use. If omitted, the repository root is used.

## clusters

Specifies the git project where the [cluster configuration]({{< ref "docs/reference/cluster-configs" >}}) is located. If it is omitted, the
`clusters` subdirectory of the `.kluctl.yml` project is used as the clusters config root.

It has the same form as in [deployment]({{< ref "docs/reference/kluctl-project#deployment" >}}), except that it is called `clusters`.

## sealedSecrets

Specifies the git project where the [sealed secrets]({{< ref "docs/reference/cluster-configs" >}}) are located. If it is omitted, the
`.sealed-secrets` subdirectory of the `.kluctl.yml` project is used as the sealed secrets location.

It has the same form as in [deployment]({{< ref "docs/reference/kluctl-project#deployment" >}}), except that it is called `sealedSecrets`.
