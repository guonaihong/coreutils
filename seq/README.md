# seq

#### summary
seq has the same seq command as ubuntu 18.04

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/seq/seq
```

#### usage
```console
Usage of seq:
  -f, --format string
    	use printf style floating-point FORMAT
  -s, --separator string
    	use STRING to separate numbers (default: \n)
  -w, --equal-width string
    	equalize width by padding with leading zeroes
```

#### todo feature(%a or %A)
```
seq -f "%a" 100
seq -f "%A" 100
```
