---
weight: 13
title: ToDo
---

# ToDo

- [ ] Complete section 'Internal execution steps' of the README.

- [ ] Properly handle exit codes / results.

- [ ] Rethink the format/name of temporal directories created for each MWE.

- Action:
  - [ ] Support labelling issues as `reproducible` or `fixed?`.
  - [ ] Support editing an existing comment published by github-actions bot, instead of adding a new one each time.
  - [ ] Support writing the logs and/or refs to jobs/artifacts in the body of the comment.
  - [ ] Support a `with` option/parameter to optionally make test artifacts available through a CI job.

- CLI:
  - [ ] Add *golden* files.
  - [ ] Write `get.sh` script and uncomment related info in section 'Installation > CLI tool' of the README.

- CI:
  - [ ] Implement publishing regular tagged (i.e. non-tip) releases.

- [podman.io](https://podman.io/)
  - [ ] Check if podman is supported for local execution.
  - [ ] Check if podman is supported in GHA environments (even though docker is installed by default).

- Security:
  - [ ] Ensure that github_token is not accessible for MWEs.
  - [ ] Which are the risks of executing MWEs on the host?

- [arcanis/sherlock](https://github.com/arcanis/sherlock)
  - Close obsolete issues after each release

<!--
TODO: publish issue-runner as a scratch-based docker image
Check requisites, e.g. "indocker":
 - Check if the socket is available
 - Check if there is some mechanism to share data between sibling containers
-->
