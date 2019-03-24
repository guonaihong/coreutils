# coreutils

## 简介
coreutils项目是gnu coreutils的扩展

## install coreutils
```
env GOPATH=`pwd` go get github.com/guonaihong/coreutils/coreutils
```
If you want to use the cat command
```
./coreutils cat flie
./coreutils cut -d":" -f1 /etc/passwd
./coreutils echo "hello china"
```

## install Compile command separately
```
env GOPATH=`pwd` go run github.com/guonaihong/coreutils/buildall
```
If you want to use the cat command
```
./cat flie
./cut -d":" -f1 /etc/passwd
./echo "hello china"
```

## 已完成命令如下 
* basename [详情](./basename/README.md)
* cat [详情](./cat/README.md)
* cut [详情](./cut/README.md)
* dirname [详情](./dirname/README.md)
* echo [详情](./echo/README.md)
* head [详情](./head/README.md)
* paste [详情](./paste/README.md)
* tee [详情](./tee/README.md)
* tail [详情](./tail/README.md)
* tr [详情](./tr/README.md)
* true [详情](./true/README.md)
* uniq [详情](./uniq/README.md)
* whoami [详情](./whoami/README.md)
* yes [详情](./yes/README.md)
* sleep [详情](./sleep/README.md)
