---
weight: 12
title: Development
---

# Development

**issue-runner** is composed of:

- A JavaScript action with TypeScript compile time support, unit testing with Jest and using the GHA Toolkit.
- A CLI tool written in golang.
- A GHA workflow.

## Build and test inside a container

The CLI can be developed in an official `golang` (or `golang:alpine`) container. See [Usage]() for information about required binds.

```sh
docker run --rm -it \
  -v /$(pwd)://src \
  -w //src \
  -v //var/run/docker.sock://var/run/docker.sock \
  -v issues://volume \
  golang bash
```

## Internal execution steps

Each time an issue is created or edited:

<!--
issue metadata
spawn container mounting files in /src
report the results back to GitHub
-->

- ...
- A temporal subdir is created.
- Source files defined or referred in the first message of the issue are saved to separate files.
- If a file with key `:image:` exists, it is saved as `run`.
- `run` is executed either on the host or inside a container.
- ...

## Create a new release branch

```sh
cp action.yml dist/
git checkout --orphan <BRANCH>
git rm --cached -r .
git add dist
git clean -fdx
git mv dist/* ./

git commit -am <release message>
git push origin <release branch name>
```

`master` branch of this repository contains the TypeScript sources of the action. However, these need to be compiled. A job in workflow `push.yml` is used to update branch `gha-tip` after each push that passes the tests. This kind of *auto-updated* branches need to be manually created the first time.

## Continuous integration

After each commit is pushed to `master`:

- The action is built, tested and published to branch `gha-tip`.
- `golangci-lint` and `goreleaser` are executed in subdir `tool`, and `test.sh` is executed.
- Action [eine/tip@gha-tip](https://github.com/eine/issue-runner/tip) is used to update tag `tip` and tool artifacts are published as a pre-release named `tip`.

<aside class="notice">
Version `eine/issue-runner@gha-tip` of this action will automatically retrieve the CLI tool from [github.com/eine/issue-runner/releases/tag/tip](https://github.com/eine/issue-runner/releases/tag/tip).
</aside>
