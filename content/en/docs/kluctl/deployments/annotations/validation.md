---
description: Annotations to control validation
github_branch: main
github_repo: https://github.com/kluctl/kluctl
lastmod: "2023-08-26T09:38:51+02:00"
linkTitle: Validation
path_base_for_github_subdir:
    from: .*
    to: docs/kluctl/deployments/annotations/validation.md
title: Validation
weight: 3
---





The following annotations influence the [validate](../../commands/validate.md) command.

### validate-result.kluctl.io/xxx
If this annotation is found on a resource that is checked while validation, the key and the value of the annotation
are added to the validation result, which is then returned by the validate command.

The annotation key is dynamic, meaning that all annotations that begin with `validate-result.kluctl.io/` are taken
into account.

### kluctl.io/validate-ignore
If this annotation is set to `true`, the object will be ignored while `kluctl validate` is run.