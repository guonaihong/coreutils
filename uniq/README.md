# uniq

#### 简介
uniq

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/uniq
```

#### diff
与gnu uniq差异部分，相比gnu uniq，本实现不依赖sort命令，可以实现对文件行的去重操作，或者摘出独一无二行，或者重复行。

#### usage

```console
Usage of ./uniq:
  -D, --all-repeated string
        print all duplicate lines delimit-method={none(default),prepend,separate} Delimiting is done with blank lines
  -c, --count
        prefix lines by the number of occurrences
  -d, --repeated
        only print duplicate lines
  -f, --skip-fields int
        avoid comparing the first N fields (default -2147483648)
  -i, --ignore-case
        ignore differences in case when comparing
  -s, --skip-chars int
        avoid comparing the first N characters (default -2147483648)
  -u, --unique
        only print unique lines
  -w, --check-chars int
        compare no more than N characters in lines (default -2147483648)
  -z, --zero-terminated
        end lines with 0 byte, not newline
```
#### Example
(下面的示例来自linux shell scripting cookbook) 和本人构造的测试用例

uniq只能作用于排过序的数据，因此，uniq通常都与sort命令结合使用。
```
cat sorted.txt
bash
foss
hack
```
去重
```
sort unsorted.txt|uniq
或者
sort unsorted.txt|uniq -u
```

要统计各行在文件中出现的次数, 使用下面的命令
```
sort unsorted.txt|uniq -c
1 bash
1 foss
1 hack
```
找出文件中重复的行
```
sort unsorted.txt|uniq -d
hack
```

我们可以结合-s 和 -w选项来指定键
* -s 指定跳过前N个字符
* -w指定用于比较的最大字符数

这个对比健可以作为uniq操作时的索引
```
cat uniq.txt
u:01:gnu
d:04:linux
u:01:bash
u:01:hack
```

为了只测试指定的字符(铁略前两个字符，使用接下来的两个字符)，我们使用-s 2跳过前两个字符，使用-w 2选项指定后续的两个字符:
```
sort uniq.txt | ./uniq -s 2 -w 2
d:04:linux
u:01:bash
```

我们将命令输出作为xargs命令的输入时，最好为输出的各行添加一个0值字符终止符，使用uniq命令的输入作为xargs的数据源时，
同样应当如此。如果没有使用0值字节终止符，那么在默认情况下，xargs命令会用空格来分割参数。例如，来自stdin的文本行"this is a line"
会被xargs视为4个不同的参数。如果使用0值字节终止符，那么\0就被作为定界符，此时，包含空格的行就能够被正确的解析为单个参数。
-z 选项可以生成由0值字节终止的输出
```
uniq -z file.txt
```

下面的命令将删除所有指定的文件，这些文件的名字是从files.txt中读取的。
```
uniq -z file.txt | xargs -0 rm
```
如果某个文件名出现多次,uniq命令只会将这个文件名写入stdout一次，这样就可以避免出现rm: cannot remove FILENAME: No such file or directory

测试 -all-repeated=prepend选项
```
echo -e "\naa\naa\ndd\n11\nbb\nbb\ncc\ncc\nff"|LC_ALL=C ./uniq --all-repeated=prepend

aa
aa

bb
bb

cc
cc
```

测试 -all-repeated=separate选项
```
echo -e "aa\naa\ndd\n11\nbb\nbb\ncc\ncc\nff"|LC_ALL=C ./uniq --all-repeated=separate
aa
aa

bb
bb

cc
cc
```

测试 --group=separate
```
echo  -e "222\n222\n111\n111\n333\naaa\n" |./uniq --group=separate|cat -n
     1  222
     2  222
     3  
     4  111
     5  111
     6  
     7  333
     8  
     9  aaa
    10  
    11  

```

测试 --group=prepend
```
echo  -e "222\n222\n111\n111\n333\naaa\n" |./uniq --group=prepend|cat -n
     1  
     2  222
     3  222
     4  
     5  111
     6  111
     7  
     8  333
     9  
    10  aaa
    11  
    12  
```
测试 --group=both
```
echo  -e "222\n222\n111\n111\n333\naaa\n" |./uniq --group=both|cat -n
     1  
     2  222
     3  222
     4  
     5  111
     6  111
     7  
     8  333
     9  
    10  aaa
    11  
    12  
    13  

```
测试 --group=append
```
echo  -e "222\n222\n111\n111\n333\naaa\n" |./uniq --group=append|cat -n
     1  222
     2  222
     3  
     4  111
     5  111
     6  
     7  333
     8  
     9  aaa
    10  
    11  
    12  

```
