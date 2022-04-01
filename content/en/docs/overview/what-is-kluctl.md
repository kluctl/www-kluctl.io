---
title: "What is kluctl?"
linkTitle: "What is kluctl?"
weight: 1
description: >
    What is kluctl, it's motivation and history?
---

## What is kluctl?

Kluctl is a CLI driven and declarative tool that unifies operations around Kubernetes deployments.

Deployments might be as simple as a single nginx Helm chart being deployed to a single cluster. It might however also be
indefinitely complex and for example provision everything that is needed for a functional/useful cluster
(e.g. CNI, ingress controllers, monitoring components and so on), allowing you to go from a naked EKS cluster to a fully
functional cluster with a single CLI invocation.

Kluctl is centered around "targets", which can be a cluster or a specific environment (e.g. test, dev, prod, ...) on one
or multiple clusters. Targets can be deployed, diffed, pruned, deleted, and so on. The idea is to have the same set of
operations for every target, no matter how simple or complex the deployment and/or target is.
