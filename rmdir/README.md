# rmdir

#### 简介
rmdir 使用gnu的命令行选项(已适配到ubuntu 16.04)

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/rmdir/rmdir
```

#### 命令行选项
```console
Usage of rmdir:
  -ignore-fail-on-non-empty
    	ignore each failure that is solely because a directory
    	is non-empty
  -p, --parent
    	remove DIRECTORY and its ancestors; e.g., 'rmdir -p a/b/c' is
    	similar to 'rmdir a/b/c a/b a'
```
