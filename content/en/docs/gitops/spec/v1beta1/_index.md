---
description: gitops.kluctl.io/v1beta1 documentation
github_branch: main
github_repo: https://github.com/kluctl/kluctl
lastmod: "2023-10-17T00:30:26+02:00"
linkTitle: v1beta1 specs
path_base_for_github_subdir:
    from: .*
    to: docs/gitops/spec/v1beta1/README.md
title: v1beta1 specs
weight: 10
---



# gitops.kluctl.io/v1beta1

This is the v1beta1 API specification for defining continuous delivery pipelines
of Kluctl Deployments.

## Specification

- [KluctlDeployment CRD](kluctldeployment.md)
    + [Spec fields](kluctldeployment.md#spec-fields)
    + [Reconciliation](kluctldeployment.md#reconciliation)
    + [Kubeconfigs and RBAC](kluctldeployment.md#kubeconfigs-and-rbac)
    + [Credentilas](kluctldeployment.md#credentials)
    + [Secrets Decryption](kluctldeployment.md#secrets-decryption)
    + [Status](kluctldeployment.md#status)
