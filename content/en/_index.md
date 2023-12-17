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
  paragraph: Easily handle deployments of any size, complexity, and across various environments using Kluctl.
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
  title: Take Your Deployments to the Next Level
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
  paragraph: Don’t waste time switching between multiple tools to manage your deployments. Kluctl is one tool to rule them all.
  featureList:
    # feature list item 1
    - title : "Native Git support"
      icon : "images/home/icons/git-icon.svg"
      description : "Easily deploy remote Kluctl projects or externalize parts of your Kluctl project."
    # feature list item 2
    - title : "Multiple environments"
      icon : "images/home/icons/environment-icon.svg"
      description : "Deploy to multiple environments (dev, test, prod etc.) with different configurations."
    # feature list item 3
    - title : "Multiple clusters"
      icon : "images/home/icons/clusters-icon.svg"
      description : "Manage multiple target clusters (in multiple clouds or bare-metal)."
    # feature list item 4
    - title : "Configuration and templating"
      icon : "images/home/icons/settings-icon.svg"
      description : "Kluctl allows to use templating in nearly all places, making it easy to have dynamic configuration."
    # feature list item 5
    - title : "Helm and Kustomize"
      icon : "images/home/icons/kustomize-icon.svg"
      description : "The Helm and Kustomize integrations reuse of 3rd party charts & Kustomisation."
    # feature list item 6
    - title : "Know what went wrong"
      icon : "images/home/icons/wrong-icon.svg"
      description : "Kluctl will show you what part of your deployment failed and why."

################## See Kluctl in action section ##################
actionSection:
  title: See Kluctl in Action
  video: "video/Kluctl demo video.mp4" # relative to the static folder

################## News and updates section ##################
newsSection:
  title: News and Updates
---
