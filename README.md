<p align="center">
  <a title="Codacy" href="https://app.codacy.com/manual/eine/issue-runner/dashboard"><img alt="Codacy grade" src="https://img.shields.io/codacy/grade/66830b37677941949d500400e2c7d1c8?longCache=true&label=quality&logo=codacy&style=flat-square"></a><!--
  -->
  <a title="Go Report Card" href="https://goreportcard.com/report/github.com/eine/issue-runner"><img src="https://goreportcard.com/badge/github.com/eine/issue-runner?longCache=true&style=flat-square"></a><!--
  -->
  <a title="godoc.org" href="https://godoc.org/github.com/eine/issue-runner/tool"><img src="http://img.shields.io/badge/godoc-reference-5272B4.svg?longCache=true&style=flat-square"></a><!--
  -->
  <a title="Dependency Status" href="https://david-dm.org/eine/issue-runner"><img src="https://img.shields.io/david/eine/issue-runner.svg?longCache=true&style=flat-square&label=deps"></a><!--
  -->
  <a title="DevDependency Status" href="https://david-dm.org/eine/issue-runner?type=dev"><img src="https://img.shields.io/david/dev/eine/issue-runner.svg?longCache=true&style=flat-square&label=devdeps"></a><!--
  -->
</p>

**issue-runner** is a toolkit to retrive, set up and run Minimal Working Examples (MWEs). MWEs are defined in markdown files (such as the first comment in a GitHub issue), and external tarball(s)/zipfile(s)/file(s) can be included. It extracts sources to separate files, (optionally) invokes docker, executes the entrypoint, and cleans up.

The main use case for this toolkit is to be added to a GitHub Actions (GHA) workflow in order to monitor the issues in a repository and optionally report status/results by:

- labelling issues as `reproducible` or `fixed?`,
- adding a comment to the issue with logs and/or refs to jobs/artifacts,
- and/or making test artifacts available through a CI job

Nonetheless, the CLI tool can also be used to set up and test any MWE or issue locally.
