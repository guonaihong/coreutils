# shuf

#### 简介
shuf 使用gnu的命令行选项(已适配到ubuntu 16.04)

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/shuf/shuf
```

#### 命令行选项
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
