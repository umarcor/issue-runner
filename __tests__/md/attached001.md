The attached [:mwe:genericIssue.zip](https://github.com/ghdl/ghdl/files/2041677/genericIssue.zip) zipfile contains a MWE that can be executed with:

``` bash :image:ghdl/ghdl:buster-mcode
ghdl -a --std=08 genericIssue/LinkedListPkg.vhd
ghdl -a --std=08 genericIssue/Test_LinkedListPkg.vhd
ghdl -e --std=08 Test_LinkedListPkg
```

This is a valid file [:mwe:tb_mwe.vhd.txt](https://github.com/VUnit/vunit/files/2037481/tb_mwe.vhd.txt), along with some other [reference](https://example.com) and a couple of [:mwe:tricky](https://example.com/data) [:mwe:cases](https://example.com/data.csv). A last [:mwe:link.txt](./textfile.txt).
