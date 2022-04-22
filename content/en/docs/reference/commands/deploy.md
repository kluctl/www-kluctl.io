---
title: "deploy"
linkTitle: "deploy"
weight: 10
description: >
    deploy command
---

## Command
<!-- BEGIN SECTION "deploy" "Usage" false -->
Usage: kluctl deploy

Deploys a target to the corresponding cluster

This command will also output a diff between the initial state and the state after deployment. The format of this diff
is the same as for the 'diff' command. It will also output a list of prunable objects (without actually deleting them).

<!-- END SECTION -->

## Arguments
The following sets of arguments are available:
1. [project arguments]({{< ref "./common-arguments#project-arguments" >}})
1. [image arguments]({{< ref "./common-arguments#image-arguments" >}})
1. [inclusion/exclusion arguments]({{< ref "./common-arguments#inclusionexclusion-arguments" >}})

In addition, the following arguments are available:
<!-- BEGIN SECTION "deploy" "Misc arguments" true -->
```
Misc arguments:
  Command specific arguments.

  -y, --yes                                Suppresses 'Are you sure?' questions and proceeds as if you would answer
                                           'yes' ($KLUCTL_YES).
      --dry-run                            Performs all kubernetes API calls in dry-run mode ($KLUCTL_DRY_RUN).
      --force-apply                        Force conflict resolution when applying. See documentation for details
                                           ($KLUCTL_FORCE_APPLY)
      --replace-on-error                   When patching an object fails, try to replace it. See documentation for more
                                           details ($KLUCTL_REPLACE_ON_ERROR).
      --force-replace-on-error             Same as --replace-on-error, but also try to delete and re-create objects. See
                                           documentation for more details ($KLUCTL_FORCE_REPLACE_ON_ERROR).
      --abort-on-error                     Abort deploying when an error occurs instead of trying the remaining
                                           deployments ($KLUCTL_ABORT_ON_ERROR)
      --hook-timeout=5m                    Maximum time to wait for hook readiness. The timeout is meant per-hook.
                                           Timeouts are in the duration format (1s, 1m, 1h, ...). If not specified, a
                                           default timeout of 5m is used ($KLUCTL_HOOK_TIMEOUT).
  -o, --output-format=OUTPUT-FORMAT,...    Specify output format and target file, in the format 'format=path'. Format
                                           can either be 'text' or 'yaml'. Can be specified multiple times. The actual
                                           format for yaml is currently not documented and subject to change
                                           ($KLUCTL_OUTPUT_FORMAT).
      --render-output-dir=STRING           Specifies the target directory to render the project into. If omitted, a
                                           temporary directory is used ($KLUCTL_RENDER_OUTPUT_DIR).
      --no-wait                            Don't wait for objects readiness' ($KLUCTL_NO_WAIT)

```
<!-- END SECTION -->

### --parallel
kluctl runs deployments sequentially and in-order by default. This options allows kluctl to perform all deployments
in parallel, which speeds up the deployment significantly.

Due to the nature of parallel deployments, no guarantees can't be made in regard to deployment order. This means for
example, that objects that are meant to be deployed into a namespace might be deployed before the namespace is deployed,
resulting in failure.

### --force-apply
kluctl implements deployments via [server-side apply](https://kubernetes.io/reference/using-api/server-side-apply/)
and a custom automatic conflict resolution algorithm. This algurithm is an automatic implementation of the
"[Don't overwrite value, give up management claim](https://kubernetes.io/reference/using-api/server-side-apply/#conflicts)"
method. It should work in most cases, but might still fail. In case of such failure, you can use `--force-apply` to
use the "Overwrite value, become sole manager" strategy instead.

Please note that this is a risky operation which might overwrite fields which were initially managed by kluctl but were
then overtaken by other managers (e.g. by operators). Always use this option with caution and perform a dry-run
before to ensure nothing unexpected gets overwritten.

### --replace-on-error
In some situations, updating Kubernetes objects is not possible, for example when modified fields are read-only. Jobs
are a good example where this might be the case. In such cases, you can use `--replace-on-error` to instruct kluctl to
retry an update by deleting and then recreating the object.

Please note that this is a potentially risky operation, especially when an object carries some kind of important state.

### --abort-on-error
kluctl does not abort a command when an individual object fails can not be updated. It collects all errors and warnings
and outputs them instead. This option modifies the behaviour to immediately abort the command.
