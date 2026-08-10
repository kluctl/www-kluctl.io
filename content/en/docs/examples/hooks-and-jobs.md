---
title: "Hooks and jobs"
linkTitle: "hooks and jobs"
weight: 5
description: >
    Basic example on running jobs using hooks.
---

## Introduction

When trying to update a `Jobs` image kubernetes will throw an error. This is because `Jobs` only have some [mutable fields](https://kubernetes.io/docs/concepts/workloads/controllers/job/#mutable-scheduling-directives), making them hard to patch/update.
> The fields in a Job's pod template that can be updated are node affinity, node selector, tolerations, labels, annotations and scheduling gates.

This problem can be solved by using [hooks]({{% ref "docs/kluctl/deployments/hooks.md" %}}) since hooked resources are deleted by default right before creation. Thus bypassing the entire initial patch/update problem.



## Problem Scenario
Let's say we want to deploy the following perl `Job` using kluctl.

```yaml 
# job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: perl:5.34.0
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4
```

```sh
$ kluctl deploy
...
New objects:
  default/Job/pi
...
$ kubectl get 
NAME   STATUS     COMPLETIONS   DURATION   AGE
pi     Complete   1/1           23s        1m5s
```
Cool we got the initial `Job` deployed.

If we now change the `command` (or any other immutable field)

```yaml 
# job.yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: perl:5.34.0
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(20)"]
      restartPolicy: Never
  backoffLimit: 4
```
And try to deploy again we get the following error.

```sh
$ kluctl deploy
Errors:
  ...
  field is immutable
  ...
```


## Using hooks
In the following example we've marked a `Job` to run after the deployment. We've also specified kluctl not to wait for its completion.

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: pi
  annotations:
    kluctl.io/hook: post-deploy # execute after deployment is performed.
    kluctl.io/hook-wait: false  # do not wait for Job to finish.
spec:
  template:
    spec:
      containers:
      - name: pi
        image: perl:5.34.0
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4
```
