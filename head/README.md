# head

#### summary
head has the same head command as ubuntu 18.04

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/head/head
```

#### usage
```console
Usage of head:
  -^\d+$, --n, --lines int
    	print the first NUM lines instead of the first 10;
    	with the leading '-', print all but the last
    	NUM lines of each file
  -c, --bytes string
    	print the first NUM bytes of each file;
    	 with the leading '-', print all but the last
    	NUM bytes of each file
  -q, --quiet, --silent
    	never print headers giving file names
  -v, --verbose
    	always print headers giving file names
  -z, --zero-terminated
    	line delimiter is NUL, not newline
```
