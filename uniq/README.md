# uniq

#### 简介
uniq

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/uniq
```

#### usage

```console
-c, --count Prefix lines with a number representing how many times they occurred.
-d, --repeated  Only print duplicated lines.
```
#### Example
(下面的示例来自linux shell scripting cookbook)

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
