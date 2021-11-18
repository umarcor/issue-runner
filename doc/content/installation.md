---
weight: 10
title: Installation
---

# Installation

## Set up a GitHub Actions workflow

> A minimal YAML workflow file:

```yml
name: 'issue?'
on:
  issues:
    types: [ opened, edited ]
jobs:
  mwe:
    runs-on: ubuntu-latest
    steps:
    - uses: umarcor/issue-runner@gha-v1
      with:
        token: ${{ secrets.GITHUB_TOKEN }}
        allow-host: false
```

Note that `with` parameters are both optional:

-   `token` is required to report feedback (labelling issues or adding comments automatically).
-   `allow-host` enables/disables running scripts on the host (without a container). For security reasons, this is discouraged and this parameter defaults to `false`.

## CLI tool

The CLI tool is a static binary written in golang, which can optionally use `docker`. It can be downloaded from [github.com/umarcor/issue-runner/releases](https://github.com/umarcor/issue-runner/releases), or it can be built from sources.

```sh
git clone https://github.com/umarcor/issue-runner
cd tool
go build -o issue-runner
```

<!--
```sh
curl -L https://raw.githubusercontent.com/umarcor/issue-runner/master/tool/get.sh | sh -
```

> You can give it a try at [play-with-docker.com](https://labs.play-with-docker.com/). Just create a node and run the command above.
-->
