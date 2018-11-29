# tee

#### 简介
tee 在完成gnu tee 功能命令基础上，并增强tee命令。

#### install
```bash
env GOPATH=`pwd` go get -u github.com/guonaihong/tee
```

#### 命令行选项
```console
Usage of tee:
  -A, --max-archive int
    	How many archive files are saved
  -a, --append
    	append to the given FILEs, do not overwrite
  -g, --gzip
    	compressed archived log files
  -s, --max-size string
    	current file maximum write size
```

#### 提供对日志文件的自动归档和裁减
* loop.sh 内容
```bash
while :;
do
    date
done
```

* tee命令保证始终只有10个归档文件(.gz)
```bash
bash loop.sh | tee -A 10 -g -a my.log
```
