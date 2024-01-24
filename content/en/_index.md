---
weight: 100
date: "2023-05-03T22:37:22+01:00"
draft: false
author: "Alexander Block"
title: "Kluctl homepage"
toc: true
description: "Kluctl is the missing glue to put together large Kubernetes deployments"
publishdate: "2023-05-03T22:37:22+01:00"

################## Hero section ##################
heroSection:
  title: Kluctl - Trust Your Deployments.
  paragraph: Easily handle Kubernetes deployments of any size, complexity, and across various environments using Kluctl.
  button:
    link: docs/
    text: Get Started
  image:
    url: images/home/hero-img.png # relative to the assets folder
    altText: Kluctl CLI
  icon:
    url: images/home/icons/hero-shape.svg # relative to the assets folder
    altText: CI/CD icon

################## Deployment section ##################
deploymentSection:
  enable: false
  title: Take Your Kubernetes Deployments to the Next Level
  paragraph: Kluctl is the missing glue to put together large Kubernetes deployments.
  image:
    url: images/home/deployment-img.png # relative to the assets folder
    altText: Deployment image
  featureList:
    - "**Manage all your deployments:** Control both infrastructure and application deployments using Kluctl." # markdown formatting allowed
    - "**No difference, whether it's simple or complex :** Kluctl makes it easy to manage complex deployments." # markdown formatting allowed

################## Features section ##################
featuresSection:
  title: Features
  featureList:
    - title: "Multiple environments"
      icon: "images/home/icons/environment-icon.svg"
      description: "Deploy to multiple environments (dev, test, prod etc.) and/or clusters with different configurations."
    - title: "Configuration and templating"
      icon: "images/home/icons/settings-icon.svg"
      description: "Kluctl allows to use templating in nearly all places, making it easy to have dynamic configuration."
    - title: "Helm and Kustomize"
      icon: "images/home/icons/kustomize-icon.svg"
      description: "The Helm and Kustomize integrations allow to reuse any 3rd party charts & Kustomisation."
    - title: "Diffs and Dry-Runs"
      icon: "images/home/icons/diff-icon.svg"
      description: "Diffs and dry-runs allow you to always know what will happen and what actually happened."
    - title: "Kluctl GitOps"
      icon: "images/home/icons/gitops-icon.svg"
      description: "Use a GitOps workflow when required, while always being able to go back to a push based workflow."
    - title: "Kluctl Webui"
      icon: "images/home/icons/webui-icon.svg"
      description: "Use the optional Webui for more visibility and control of your GitOps and CLI based deployments."

################## See Kluctl in action section ##################
actionSection:
  title: See Kluctl in Action
  video: "video/Kluctl demo video.mp4" # relative to the static folder

################## News and updates section ##################
newsSection:
  title: News and Updates
---
