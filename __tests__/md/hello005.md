This is the body of an example issue with multiple code blocks, where the entrypoint to a docker container is a shell/bash script.

```sh :image: alpine
#!/usr/bin/env sh
echo "Hello world!"
./demo.sh
source hello.sh
```

Some text in the middle...

```sh
#!/usr/bin/env sh
echo "Hello world again!"
#:file: demo.sh
```

More text...

~~~sh :file: hello.sh
#!/usr/bin/env sh
echo "Hello world last!"
~~~

This is some `example code snippet` in text.
