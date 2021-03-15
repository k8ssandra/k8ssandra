---
date: 2021-03-15
title: "K8ssandra Roadmap Update"
linkTitle: "K8ssandra Roadmap Update"
description: >
  We've updated our approach to managing the K8ssandra roadmap and would love to get your input.
author: Chris Bradford ([@bradfordcp](https://twitter.com/bradfordcp)), Jeff Carpenter ([@jscarp](https://twitter.com/jscarp))

---

Prioritization and roadmaps are often among the more challenging areas for any open source project to manage. On the K8ssandra project, we're making some improvements to how we manage our roadmap for new features. 

# Legacy Roadmap
If you've poked around the website a bit, you may have noticed the [roadmap]({{< relref "docs/roadmap">}}) page that we've been using to list a non-prioritized list of changes that have been proposed by project contributors. 

As part of our post-1.0 release cleanup work, we've identified that a more flexible, robust approach was needed, and that it would be much easier to track in Git. The old roadmap page will be available for a short time but eventually phased out.
 
# New Roadmap 
The new roadmap is implemented as a [GitHub project](https://github.com/orgs/k8ssandra/projects/6) under the K8ssandra [organization](https://github.com/k8ssandra/). The roadmap is broken out by quarter, with a lane for each of the next 4 quarters, and a lane for other items that, while important, are not yet prioritized. Currently, each item in this list is a "note", which boils down to a single text field. Over the coming days, these will be transitioned into issues within the roadmap [repository](https://github.com/k8ssandra/roadmap) which allow for richer integration with related resources inside Github (milestones, labels, etc) and a cleaner UI. 

Chris Bradford shared a quick overview of this new roadmap and contents at the most recent meeting of the Cassandra Kubernetes Special Interest Group (SIG). You can see the video here:

{{< youtube id="RXAdWLKw450" yt_start="2315" >}}

Going forward, we hope to use this roadmap to manage our work and define epics that span repository boundaries.

We'd love to have your feedback on the items in the project and their prioritization, as well as any refinements to how we review and document changes.
