**issue-runner** is a toolkit to retrive, set up and run Minimal Working Examples (MWEs). MWEs are defined in a markdown file (such as the first comment in a GitHub issue), and external tarball(s)/zipfile(s)/file(s) can be included. The main use case for this toolkit is to be added to a GitHub Actions (GHA) workflow in order to monitor the issues in a repository and optionally report status/results by:

- labelling issues as `reproducible` or `fixed?`,
- adding a comment to the issue with logs and/or refs to jobs/artifacts,
- and/or making test artifacts available through a CI job

Nonetheless, the CLI tool can also be used to set up and test any MWE or issue locally.

# Installation

## Set up a GitHub Actions workflow

The following block shows a minimal YAML workflow file:

```yml
name: 'issue?'
on:
  issues:
    types: [ opened, edited ]
jobs:
  mwe:
    runs-on: ubuntu-latest
    steps:
    - uses: 1138-4EB/issue-runner@gha-v1
      with:
        token: ${{ secrets.GITHUB_TOKEN }}
        allow-host: false
```

Note that `with` parameters are both optional:

- `token` is required to report feedback (labelling issues or adding comments automatically).
- `allow-host` enables/disables running scripts on the host (without a container). For security reasons, this is discouraged and this parameter defaults to `false`.

## CLI tool

The CLI tool is a static binary written in golang, which can optionally use `docker`. It can be downloaded from [github.com/1138-4EB/issue-runner/releases](https://github.com/1138-4EB/issue-runner/releases), or it can be built from sources:

```
git clone https://github.com/1138-4EB/issue-runner
cd tool
go build -o issue-runner
```

<!--
```sh
curl -L https://raw.githubusercontent.com/1138-4EB/issue-runner/master/tool/get.sh | sh -
```

> You can give it a try at [play-with-docker.com](https://labs.play-with-docker.com/). Just create a node and run the command above.
-->

# Usage

## Supported markdown syntax

**issue-runner** scans the (markdown) body to extract:

- Code blocks with either the body or the language definition string matching `:file:.*`:

~~~md
```sh :file: hello.sh
#!/usr/bin/env sh
echo "Hello world!"
```
~~~

~~~md
```sh
#!/usr/bin/env sh
echo "Hello world!"
#:file: hello.sh
```
~~~

Note that, in the latter, `:file:` is prepended with a comment symbol that depends on the target language.

- Attached/linked text files, tarballs or zipfiles with the name to the reference matching `:mwe:.*`:

~~~md
[:mwe:filename.ext.txt](URL)
[:mwe:filename.tar.gz](URL)
~~~

Since GitHub allows uploading files with a limited set of extensions, issue-runner expects the user to append `.txt` to attached source filenames. This extra extension is trimmed. The exception to rule above are tarballs, zipfiles or txt files. No extra extension needs to be appended. The content of these is automatically extracted.

## Entrypoint

One, and only one, of the code blocks should contain `:image:.*` instead of `:file:.*`. That file will be the entrypoint to an OCI container. For example:

~~~md
```sh :image: debian:buster-slim
echo "Hello world!"
```
~~~

~~~md
```py
#!/usr/bin/env python3

print('Hello world!')

#:image: python:slim-buster
```
~~~

Alternatively, if no `:image:` is defined, the file which is named `run` will be used as the entrypoint to execute the MWE on the host.

## CLI

> NOTE: automatically labelling/commenting features are not included in the CLI tool. These features are implemented in the GitHub Action only.

At least one of the following references needs to be provided:

- Path or URL to markdown file: `issue-runner path/to/markdown/file.md`
- Full URL to a GitHub issue: `issue-runner 'https://github.com/user/repo/issues/number'`
- Short reference of a GitHub issue (see [GitHub Help: Autolinked references and URLs](https://help.github.com/articles/autolinked-references-and-urls/#issues-and-pull-requests)): `issue-runner 'user/repo#number'`

---

Providing a list of identifiers is also supported. For example:

```sh
issue-runner \
  'https://raw.githubusercontent.com/1138-4EB/issue-runner/master/examples/vunit_py.md' \
  test/vunit_sh.md \
  'VUnit/vunit#337' \
  'ghdl/ghdl#579' \
  'ghdl/ghdl#584'
```

> NOTE: multiple references can be managed as a single MWE with flag `-m|--merge`.

---

MWEs defined in a single body can be read through *stdin*. For example:

```sh
cat ./__tests__/md/hello001.md | ./issue-runner -
# or
./issue-runner -y - < ./__tests__/md/hello003.md
```

# Development

**issue-runner** is composed of:

- A JavaScript action with TypeScript compile time support, unit testing with Jest and using the GHA Toolkit.
- A CLI tool written in golang.
- A GHA workflow.

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

`master` branch of this repository contains the TypeScript sources of the action. However, these need to be compiled. A job in workflow `push.yml` is used to update branch `gha-tip` after each push that passes the tests. This kind of *auto-updated* branches need to be manually created the first time:

```bash
git checkout -b <release branch name>

rm -Rf node_modules
git rm -rf *.json *config.js *.lock .github .gitignore .v0 __tests__ src examples cli.sh
git add dist
git mv dist/* ./
git add .

git commit -am <release message>
git push origin <release branch name>
```

> NOTE: this procedure is based on https://github.com/actions/typescript-action#publish-to-a-distribution-branch

## Set up CI

A 'Deploy key' pair needs to be configured in order to automatically `git push` branch `gha-tip`:

```sh
ssh-keygen -t ed25519
```
- Repository 'Settings':
  - Add public key to 'Deploy keys'.
  - Add private key as 'GHA_DEPLOY_KEY' in 'Secrets'.

## Continuous integration

After each commit is pushed to `master`:

- The action is built, tested and published to branch `gha-tip`.
- `golangci-lint` and `goreleaser` are executed in subdir `tool`, and `test.sh` is executed.
- Action [1138-4EB/tip@gha-tip](https://github.com/1138-4EB/issue-runner/tip) is used to update tag `tip` and tool artifacts are published as a pre-release named `tip`.

> NOTE: version `1138-4EB/issue-runner@gha-tip` of this action will automatically retrieve the CLI tool from [github.com/1138-4EB/issue-runner/releases/tag/tip](https://github.com/1138-4EB/issue-runner/releases/tag/tip).

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
