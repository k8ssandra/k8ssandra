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
The new roadmap is implemented as a [Git project](https://github.com/orgs/k8ssandra/projects/6) under the main K8ssandra GitHub org. The roadmap is broken out by quarter, with a lane for each of the next 4 quarters, and a lane for other items that, while important, are not yet prioritized. 

Chris Bradford shared a quick overview of this new roadmap and contents at the most recent meeting of the Cassandra Kubernetes Special Interest Group (SIG). You can see the video here:

Going forward, we hope to use this roadmap to manage our work. We'll convert items on the project board into issues as work is started on them.

We'd love to have your feedback on the items in the project and their prioritization, as well as any refinements to how we review and document changes.
