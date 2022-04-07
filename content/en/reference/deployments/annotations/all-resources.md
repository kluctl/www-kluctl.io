---
title: "All resources"
linkTitle: "All resources"
weight: 1
description: >
  Annotations on all resources
---

The following annotations control the behavior of the `deploy` and related commands.

## Control deploy behavior

The following annotations control deploy behavior, especially in regard to conflict resolution.

### kluctl.io/force-apply
If set to "true", the whole resource will be force-applied, meaning that all fields will be overwritten in case of
field manager conflicts.

### kluctl.io/force-apply-field
Specifies a [JSON Path](https://goessner.net/articles/JsonPath/) for fields that should be force-applied. Matching
fields will be overwritten in case of field manager conflicts.

If more than one field needs to be specified, add `-xxx` to the annotation key, where `xxx` is an arbitrary number.

## Control deletion/pruning

The following annotations control how delete/prune is behaving.

### kluctl.io/skip-delete
If set to "true", the annotated resource will not be deleted when [delete]({{< ref "reference/commands/delete" >}}) or
[prune]({{< ref "reference/commands/prune" >}}) is called.

### kluctl.io/skip-delete-if-tags
If set to "true", the annotated resource will not be deleted when [delete]({{< ref "reference/commands/delete" >}}) or
[prune]({{< ref "reference/commands/prune" >}}) is called and inclusion/exclusion tags are used at the same time.

This tag is especially useful and required on resources that would otherwise cause cascaded deletions of resources that
do not match the specified inclusion/exclusion tags. Namespaces are the most prominent example of such resources, as
they most likely don't match exclusion tags, but cascaded deletion would still cause deletion of the excluded resources.

## Control diff behavior

The following annotations control how diffs are performed.

### kluctl.io/diff-name
This annotation will override the name of the object when looking for the in-cluster version of an object used for
diffs. This is useful when you are forced to use new names for the same objects whenever the content changes, e.g.
for all kinds of immutable resource types.

Example (filename job.yml):
```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: myjob-{{ load_sha256("job.yml", 6) }}
  annotations:
    kluctl.io/diff-name: myjob
spec:
  template:
    spec:
      containers:
      - name: hello
        image: busybox
        command: ["sh",  "-c", "echo hello"]
      restartPolicy: Never
```

Without the `kluctl.io/diff-name` annotation, any change to the `job.yml` would be treated as a new object in resulting
diffs from various commands. This is due to the inclusion of the file hash in the job name. This would make it very hard
to figure out what exactly changed in an object.

With the `kluctl.io/diff-name` annotation, kluctl will pick an existing job from the cluster with the same diff-name
and use it for the diff, making it a lot easier to analyze changes. If multiple objects match, the one with the youngest
`creationTimestamp` is chosen.

Please note that this will not cause old objects (with the same diff-name) to be prunes. You still have to regularely
prune the deployment.

### kluctl.io/ignore-diff
If set to "true", the whole resource will be ignored while calculating diffs.

### kluctl.io/ignore-diff-field
Specifies a [JSON Path](https://goessner.net/articles/JsonPath/) for fields that should be ignored while calculating
diffs.

If more than one field needs to be specified, add `-xxx` to the annotation key, where `xxx` is an arbitrary number.

## Control downscale behavior

The following annotations control how [downscale]({{< ref "reference/commands/downscale" >}}) is behaving.

### kluctl.io/downscale-patch
Describes how a [downscale]({{< ref "reference/commands/downscale" >}}) on a resource can be done via a json patch. This is useful on
CRD based resources where no automatic downscale can be performed by kluctl.

Example (filename job.yml):
```yaml
apiVersion: kibana.k8s.elastic.co/v1
kind: Kibana
metadata:
  name: kibana
  annotations:
    kluctl.io/downscale-patch: |
      - op: replace
        path: /spec/count
        value: 0
spec:
  version: 7.14.1
  count: 1
```

If more than one patch needs to be specified, add `-xxx` to the annotation key, where `xxx` is an arbitrary number.

### kluctl.io/downscale-ignore
If set to "true", the resource will be ignored while [downscale](./commands.md#downscale) is executed.

### kluctl.io/downscale-delete
If set to "true", the resource will be deleted while [downscale](./commands.md#downscale) is executed.
