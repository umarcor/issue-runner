# issue-runner

These are a couple of scripts that allow to execute Minimal Working Examples (MWEs) defined in the body of GitHub issues.
For example, in order to run [VUnit/vunit#337](https://github.com/VUnit/vunit/issues/337):

```
curl -L https://raw.githubusercontent.com/1138-4EB/issue-runner/master/runner.sh | sh -s VUnit/vunit#337
```

NOTE: you can give it a try at [play-with-docker.com](https://labs.play-with-docker.com/). Just create a node and run the command above.

- A subdir named `VUnit--vunit--337` is created.
- The source files defined or referred in the first message of the issue are saved to separate files: `run.sh`, `tb_repro.vhd` and `tb_mwe.vhd`.
- `run.sh` is executed.

Other examples: [ghdl/ghdl#579](https://github.com/ghdl/ghdl/issues/579), [ghdl/ghdl#584](https://github.com/ghdl/ghdl/issues/584)...

## Parser

- The script parses the body of the message to search for `#>> filename.ext` lines. All the content between such a token and the next one is saved to the file. Therefore, a las dummy token, `#>> anything` must be added.
- GitHub allows uploading files with a limited set of extensions. So, `issue-runner` scripts expect the user to append `.txt` to the attached source filenames. This extra extension is trimmed.
  - The exception to rule above are tarballs and zipfiles. No extra extension needs to be appended. The content of these is automatically extracted.

## Requirements

- The frontend is a bash script. The first and single argument must be a raw issue reference (see [GitHub Help: Autolinked references and URLs](https://help.github.com/articles/autolinked-references-and-urls/#issues-and-pull-requests)), the URL of a file or the path to a local file.
- Python and/or docker are required. The script automatically detects if `python` or `python3` are available in the PATH. A docker `python:slim-stretch` container is used if none is found.
