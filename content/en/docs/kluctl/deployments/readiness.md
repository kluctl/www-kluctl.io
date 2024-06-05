---
description: Definition of readiness.
github_branch: main
github_repo: https://github.com/kluctl/kluctl
lastmod: "2023-09-27T11:33:53+02:00"
linkTitle: Readiness
path_base_for_github_subdir:
    from: .*
    to: docs/kluctl/deployments/readiness.md
title: Readiness
weight: 7
---





There are multiple places where kluctl can wait for "readiness" of resources, e.g. for hooks or when `waitReadiness` is
specified on a deployment item. Readiness depends on the resource kind, e.g. for a Job, kluctl would wait until it
finishes successfully.

## Control via Annotations

Multiple [annotations](./annotations/) control the behaviour when waiting for readiness of resources. These are
the following annoations:

- [kluctl.io/wait-readiness in resources](./annotations/all-resources.md#kluctliowait-readiness)
- [kluctl.io/wait-readiness in kustomization.yaml](./annotations/kustomization.md#kluctliowait-readiness)
- [kluctl.io/is-ready](./annotations/all-resources.md#kluctliois-ready)
- [kluctl.io/hook-wait](./annotations/hooks.md#kluctliohook-wait)
