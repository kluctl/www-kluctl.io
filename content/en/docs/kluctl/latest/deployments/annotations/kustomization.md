---
description: Annotations on the kustomization.yaml resource
github_repo: https://github.com/kluctl/kluctl
lastmod: "2023-09-15T15:13:22+02:00"
linkTitle: Kustomize
path_base_for_github_subdir:
    from: .*
    to: main/docs/kluctl/deployments/annotations/kustomization.md
title: Kustomize
weight: 4
---





Even though the `kustomization.yaml` from Kustomize deployments are not really Kubernetes resources (as they are not
really deployed), they have the same structure as Kubernetes resources. This also means that the `kustomization.yaml`
can define metadata and annotations. Through these annotations, additional behavior on the deployment can be controlled.

Example:
```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization

metadata:
  annotations:
    kluctl.io/barrier: "true"
    kluctl.io/wait-readiness: "true"

resources:
  - deployment.yaml
```

### kluctl.io/barrier
If set to `true`, kluctl will wait for all previous objects to be applied (but not necessarily ready). This has the
same effect as [barrier](../../deployments/deployment-yml.md#barriers) from deployment projects.

### kluctl.io/wait-readiness
If set to `true`, kluctl will wait for readiness of all objects from this kustomization project. Readiness is defined
the same as in [hook readiness](../../deployments/readiness.md). Waiting happens after all resources from the current
deployment item have been applied.
