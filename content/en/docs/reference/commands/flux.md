---
title: "flux"
linkTitle: "flux"
weight: 10
description: >
    flux command
---

## Top Command
<!-- BEGIN SECTION "flux" "Usage" false -->
Usage: kluctl flux [command] [flags]

Interact with [flux-kluctl-controller](https://github.com/kluctl/flux-kluctl-controller) KluctlDeployment objects.

<!-- END SECTION -->

## Sub Commands
The following sets of arguments are available:
1. [reconcile]({{< ref "./flux#reconcile" >}})
1. [suspend]({{< ref "./flux#supend" >}})
1. [resume]({{< ref "./flux#resume" >}})

In addition, the following arguments are available:
<!-- BEGIN SECTION "flux" "Misc arguments" false -->
```bash
Usage: kluctl flux [command]

Flux sub-commands

Global arguments:
      --cpu-profile string   Enable CPU profiling and write the result to the given path
      --debug                Enable debug logging
      --no-update-check      Disable update check on startup

Commands:
  reconcile   Reconcile KluctlDeployment
  resume      Resume KluctlDeployment
  suspend     Suspend KluctlDeployment

Use "kluctl flux [command] --help" for more information about a command.
```
<!-- END SECTION -->

## reconcile
<!-- BEGIN SECTION "reconcile" "Misc arguments" false -->
Usage: kluctl flux reconcile [flags]

Reconcile KluctlDeployment and it's Source.
```bash
Global arguments:
      --cpu-profile string   Enable CPU profiling and write the result to the given path
      --debug                Enable debug logging
      --no-update-check      Disable update check on startup

Flux arguments:
  EXPERIMENTAL: Subcommands for interaction with flux-kluctl-controller

  -k, --kluctl-deployment string   Name of the KluctlDeployment to interact with
  -n, --namespace string           Namespace where KluctlDeployment is located
      --no-wait                    Don't wait for objects readiness'
      --with-source                --with-source will annotate Source object as well, triggering pulling
```
<!-- END SECTION -->

## suspend

<!-- BEGIN SECTION "suspend" "Misc arguments" false -->
Usage: kluctl flux suspend [flags]

Suspend KluctlDeployment from reconciling.
```bash
Global arguments:
      --cpu-profile string   Enable CPU profiling and write the result to the given path
      --debug                Enable debug logging
      --no-update-check      Disable update check on startup

Flux arguments:
  EXPERIMENTAL: Subcommands for interaction with flux-kluctl-controller

  -k, --kluctl-deployment string   Name of the KluctlDeployment to interact with
  -n, --namespace string           Namespace where KluctlDeployment is located
      --no-wait                    Don't wait for objects readiness'
      --with-source                --with-source will annotate Source object as well, triggering pulling
```
<!-- END SECTION -->

## resume
<!-- BEGIN SECTION "resume" "Misc arguments" false -->
Usage: kluctl flux resume [flags]

Resume KluctDeployment reconciliation.
```bash
Global arguments:
      --cpu-profile string   Enable CPU profiling and write the result to the given path
      --debug                Enable debug logging
      --no-update-check      Disable update check on startup

Flux arguments:
  EXPERIMENTAL: Subcommands for interaction with flux-kluctl-controller

  -k, --kluctl-deployment string   Name of the KluctlDeployment to interact with
  -n, --namespace string           Namespace where KluctlDeployment is located
      --no-wait                    Don't wait for objects readiness'
      --with-source                --with-source will annotate Source object as well, triggering pulling
```
<!-- END SECTION -->