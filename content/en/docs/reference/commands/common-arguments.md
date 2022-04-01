---
title: "Common Arguments"
linkTitle: "Common Arguments"
weight: 1
description: >
    Common arguments
---

A few sets of arguments are common between multiple commands. These arguments are still part of the command itself and
must be placed *after* the command name.

## Global arguments

These arguments are available for all commands.

<!-- BEGIN SECTION "deploy" "Global arguments" true -->
```
Global arguments:
  -h, --help                Show context-sensitive help.
  -v, --verbosity="info"    Log level (debug, info, warn, error, fatal, panic) ($KLUCTL_VERBOSITY).
      --no-update-check     Disable update check on startup ($KLUCTL_NO_UPDATE_CHECK)

```
<!-- END SECTION -->

## Project arguments

These arguments are available for all commands that are based on a Kluctl project.
They control where and how to load the kluctl project and deployment project.

<!-- BEGIN SECTION "deploy" "Project arguments" true -->
```
Project arguments:
  Define where and how to load the kluctl project and its components from.

  -p, --project-url=STRING              Git url of the kluctl project. If not specified, the current directory will be
                                        used instead of a remote Git project ($KLUCTL_PROJECT_URL)
  -b, --project-ref=STRING              Git ref of the kluctl project. Only used when --project-url was given
                                        ($KLUCTL_PROJECT_REF).
  -c, --project-config=STRING           Location of the .kluctl.yml config file. Defaults to $PROJECT/.kluctl.yml
                                        ($KLUCTL_PROJECT_CONFIG)
      --local-clusters=STRING           Local clusters directory. Overrides the project from .kluctl.yml
                                        ($KLUCTL_LOCAL_CLUSTERS)
      --local-deployment=STRING         Local deployment directory. Overrides the project from .kluctl.yml
                                        ($KLUCTL_LOCAL_DEPLOYMENT)
      --local-sealed-secrets=STRING     Local sealed-secrets directory. Overrides the project from .kluctl.yml
                                        ($KLUCTL_LOCAL_SEALED_SECRETS)
      --from-archive=STRING             Load project (.kluctl.yml, cluster, ...) from archive. Given path can either be
                                        an archive file or a directory with the extracted contents
                                        ($KLUCTL_FROM_ARCHIVE).
      --from-archive-metadata=STRING    Specify where to load metadata (targets, ...) from. If not specified, metadata
                                        is assumed to be part of the archive ($KLUCTL_FROM_ARCHIVE_METADATA).
      --cluster=STRING                  Specify/Override cluster ($KLUCTL_CLUSTER)
  -t, --target=STRING                   Target name to run command for. Target must exist in .kluctl.yml
                                        ($KLUCTL_TARGET).
  -a, --arg=ARG,...                     Template argument in the form name=value ($KLUCTL_ARG)

```
<!-- END SECTION -->

## Image arguments

These arguments are available on some target based commands.
They control image versions requested by `images.get_image(...)` [calls]({{< ref "docs/reference/deployments/images#imagesget_image" >}}).

<!-- BEGIN SECTION "deploy" "Image arguments" true -->
```
Image arguments:
  Control fixed images and update behaviour.

  -F, --fixed-image=FIXED-IMAGE,...    Pin an image to a given version. Expects
                                       '--fixed-image=image<:namespace:deployment:container>=result'
                                       ($KLUCTL_FIXED_IMAGE)
      --fixed-images-file=STRING       Use .yml file to pin image versions. See output of list-images sub-command or
                                       read the documentation for details about the output format
                                       ($KLUCTL_FIXED_IMAGES_FILE)
  -u, --update-images                  This causes kluctl to prefer the latest image found in registries, based on the
                                       'latest_image' filters provided to 'images.get_image(...)' calls. Use this flag
                                       if you want to update to the latest versions/tags of all images. '-u' takes
                                       precedence over '--fixed-image/--fixed-images-file', meaning that the latest
                                       images are used even if an older image is given via fixed images
                                       ($KLUCTL_UPDATE_IMAGES).

```
<!-- END SECTION -->

## Inclusion/Exclusion arguments

These arguments are available for some target based commands.
They control inclusion/exclusion based on tags and deployment item pathes.

<!-- BEGIN SECTION "deploy" "Inclusion/Exclusion arguments" true -->
```
Inclusion/Exclusion arguments:
  Control inclusion/exclusion.

  -I, --include-tag=INCLUDE-TAG,...    Include deployments with given tag ($KLUCTL_INCLUDE_TAG).
  -E, --exclude-tag=EXCLUDE-TAG,...    Exclude deployments with given tag. Exclusion has precedence over inclusion,
                                       meaning that explicitly excluded deployments will always be excluded even if an
                                       inclusion rule would match the same deployment ($KLUCTL_EXCLUDE_TAG).
      --include-deployment-dir=INCLUDE-DEPLOYMENT-DIR,...
                                       Include deployment dir. The path must be relative to the root deployment project
                                       ($KLUCTL_INCLUDE_DEPLOYMENT_DIR).
      --exclude-deployment-dir=EXCLUDE-DEPLOYMENT-DIR,...
                                       Exclude deployment dir. The path must be relative to the root deployment project.
                                       Exclusion has precedence over inclusion, same as in --exclude-tag
                                       ($KLUCTL_EXCLUDE_DEPLOYMENT_DIR)

```
<!-- END SECTION -->
