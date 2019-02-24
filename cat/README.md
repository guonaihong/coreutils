# cat

#### 简介
cat 使用gnu的命令行选项(已适合到ubuntu 16.04)，在mac电脑上可以获得和linux 一样的体验

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/cat/cat
```

#### 命令行选项
```console
Usage of ./cat:
  -A, --show-all
        equivalent to -vET
  -E, --show-end
        display $ at end of each line
  -T, --show-tabs
        display TAB characters as ^I
  -b, --number-nonblank
        number nonempty output lines, overrides -n
  -e    equivalent to -vE
  -n, --numbe
        number all output line
  -s, --squeeze-blank
        suppress repeated empty output lines
  -t    equivalent to -vT
  -v, --show-nonprinting
        use ^ and M- notation, except for LFD and TAB
```
