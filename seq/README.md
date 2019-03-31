# seq

#### 简介
seq 使用gnu的命令行选项

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/seq/seq
```

#### 命令行选项
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
