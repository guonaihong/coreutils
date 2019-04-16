# shuf

#### summary
shuf has the same shuf command as ubuntu 18.04

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/shuf/shuf
```

#### usage
```console
Usage of shuf:
  -e, --echo
    	treat each ARG as an input line
  -i, --input-range string
    	treat each number LO through HI as an input line
  -n, --head-count int
    	output at most COUNT lines
  -o, --output string
    	write result to FILE instead of standard output
  -r, --repeat
    	output lines can be repeated
  -random-source string
    	get random bytes from FILE
  -z, --zero-terminated
    	line delimiter is NUL, not newline
```
