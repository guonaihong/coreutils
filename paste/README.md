# paste

#### 简介
paste 在部分完成gnu paste 功能命令基础上，主要解决gnu paste一个很诡异的问题

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/paste
```

#### 命令行选项
```console
Usage of paste:
  -d, --delimiters string
    	reuse characters from LIST instead of TABs (default "\t")
```

#### gnu paste 诡异问题
* gnu paste
```
paste -d "/" session_id.log session_id.log
# 输出
/c64e787e-30ae-45b5-b8ee-7037f8b92515
/c64e787e-30ae-45b5-b8ee-7037f8b92515
/c64e787e-30ae-45b5-b8ee-7037f8b92515
```

* paste
```
paste -d "/" session_id.log session_id.log

# 输出
c64e787e-30ae-45b5-b8ee-7037f8b92515/c64e787e-30ae-45b5-b8ee-7037f8b92515
c64e787e-30ae-45b5-b8ee-7037f8b92515/c64e787e-30ae-45b5-b8ee-7037f8b92515
c64e787e-30ae-45b5-b8ee-7037f8b92515/c64e787e-30ae-45b5-b8ee-7037f8b92515
```
