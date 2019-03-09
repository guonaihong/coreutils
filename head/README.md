# head

#### 简介
head 使用gnu的命令行选项(已适配到ubuntu 16.04)

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/head/head
```

#### 命令行选项
```console
Usage of ./head:
  -^\d+$, --n, --lines int
        print the first NUM lines instead of the first 10;with the leading '-', print all but the lastNUM lines of each file
  -c, --bytes int
        print the first NUM bytes of each file; with the leading '-', print all but the last NUM bytes of each file
  -q, --quiet, --silent
        never print headers giving file names
  -v, --verbose
        always print headers giving file names
  -z, --zero-terminated
        line delimiter is NUL, not newline
```
