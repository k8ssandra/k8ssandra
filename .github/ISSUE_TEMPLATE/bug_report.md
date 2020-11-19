---
name: Bug report
about: Create a report to help us improve
title: ''
labels: bug, needs-triage
assignees: ''

---

## Bug Report

<!--
Thanks for filing an issue! Before hitting the button, please answer these questions.
Fill in as much of the template below as you can. 
-->

**Describe the bug**
A clear and concise description of what the bug is.

**To Reproduce**
Steps to reproduce the behavior:
1. Go to '...'
2. Click on '....'
3. Scroll down to '....'
4. See error

**Expected behavior**
A clear and concise description of what you expected to happen.

**Screenshots**
If applicable, add screenshots to help explain your problem.

**Environment (please complete the following information):**

* Helm charts version info
<!-- list installed charts and their versions from all namespaces -->
<!-- Replace the command with its output -->
`$ helm ls -A` 

* Helm charts user-supplied values
<!-- For each k8ssandra chart involved list user-supplied values -->
<!-- Replace the commands with its output -->
`$ helm get values RELEASE_NAME` 

* Kubernetes version information:
<!-- Replace the command with its output -->
`kubectl version`

* Kubernetes cluster kind:
<!-- Insert how you created your cluster: kind, kops, bootkube, etc. -->

**Additional context**
<!-- Add any other context about the problem here. -->
