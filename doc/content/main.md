---
weight: 1
---

# Introduction

Welcome to the documentation of **issue-runner**, a toolkit to retrive, set up and run Minimal Working Examples (MWEs).

issue-runner supports MWEs defined in markdown files (such as the first comment in a GitHub issue), and external tarball(s)/zipfile(s)/file(s) can be included. It extracts sources to separate files, (optionally) invokes docker, executes the entrypoint, and cleans up.

Furthermore, a GitHub Action (GHA) is provided, which allows to easily integrate the tool in (scheduled) YAML workflows.
