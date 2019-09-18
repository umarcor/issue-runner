This is the body of an example issue for [VUnit/vunit](https://github.com/VUnit/vunit), where the entrypoint to a docker container is a python script.

```python
#!/usr/bin/env python3

from os.path import join, dirname
from vunit import VUnit
vu = VUnit.from_argv(['-v'])
vu.add_library("lib").add_source_files(join(dirname(__file__), "*.vhd"))
vu.main()

#:image: ghdl/vunit:mcode
```

```vhdl
library vunit_lib;
context vunit_lib.vunit_context;

entity tb_repro is
  generic ( runner_cfg : string );
end entity;

architecture tb of tb_repro is
begin
  main: process
  begin
    test_runner_setup(runner, runner_cfg);
    info("Hello world!");
    test_runner_cleanup(runner);
    wait;
  end process;
end architecture;

--:file: tb_repro.vhd
```

[tb_mwe.vhd.txt](https://github.com/VUnit/vunit/files/2037481/tb_mwe.vhd.txt)
