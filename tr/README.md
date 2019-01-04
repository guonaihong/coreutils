#### Example
（来自linux shell cookbook）  
tr只能通过stdin接收输入。调用格式如下：
tr [options] set1 set2
来自stdin的输入字符会按照位置从set1映射到set2，然后将输出写入stdout。set1和set2是字符类或字符组。
如果两个字符组的长度不相等，那么set2会不断复制其最后一个字符，直到长度与set1相同。
如果set2的长度大于set1，那么在set2中超出set1长度的那部分字符则全部被忽略。

要将输入中的字符由大写转换成小写，可以使用
```shell
echo "HELLO WORLD" |tr 'A-Z' 'a-z'
```
'A-Z'和'a-z'都是字符组。我们可以按照需要追加字符或字符类来构造自己的字符组。
'ABD-}'、'aA., '、'a-ce-x'以及'a-c0-9'等均是合法的集合。定义集合也很简单，不需要书写一长串连续的字符序列，
只需要使用“起始字符-终止字符”这种格式就行了。这种写法也可以和其他字符或者字符类结合使用。如果“起始字符-终止字符”
不是有效连续字符序列，那么它就会被视为含有3个元素的集合（起始字符、-和终止字符）。你也可以使用像'\t'、'\n'这种特殊字符或者其他ASCII字符

在tr中利用集合的概念，可以轻松地将字符从一个集合映射到另一个集合中。
```shell
echo 12345|tr '0-9' '9876543210'
87654
```
```shell
echo 87654|tr '9876543210' '0-9'
```
tr命令可以用来加密。ROT13是一个著名的加密算法。在ROT13算法中，字符会被移动13个位置，因此文本加密和解密都使用相同一个函数
```shell
echo "tr came, tr saw, tr conquered."|tr 'a-zA-Z' 'n-za-mN-ZA-M'
```
还可以将制表符转换成单个空格:
```shell
tr '\t' ' ' < file.txt
```
##### 用tr删除字符
tr有一个选项-d，可以通过指定需要被删除的字符集合，将出现在stdin中的特定字符清除掉:
```shell
cat file.txt | tr -d '[set1]'
#只使用set1，不使用set2
```
例如:
```shell
echo "Hello 123 world 456"|tr -d '0-9'
Hello world
```
将stdin数字删除并打印删除后的结果
##### 字符组补集
可以利用选项-c来使用set1的补集。下面的命令中,set2是可选的:
```shell
tr -c [set1] [set2]
```
如果只给出了set1，那么tr会删除所有不在set1中的字符。如果也给出了set2，tr会将不在set1中的字符转换成set2中的字符。
如果使用了-c选项，set1和set2必须都给出。如果-c与-d选项同时出现，你只能使用set1,其他所有的字符都会被删除。
下面的例子会从输入文本中删除不在补集中的所有字符:
```shell
echo hello 1 char 2 next 4|tr -d -c '0-9\n'
124
```
接下来的例子会将不在set1中的字符替换成空格
```shell
echo hello 1 char 2 next 4| tr -c '0-9' ' '
```
1 2 4
##### 用tr压缩字符
tr命令能够完成很多文本处理任务。例如，它可以删除字符串中重复出现的字符。基本实现形式如下:
```shell
tr -s '[需要被压缩的一组字符]'
```
如果你习惯在点号后面放置两个空格，你需要在不删除重复字母的情况下去掉多余的空格:
```
echo "GNU is not UNIX. Recursive right ?"|tr -s ' '
```
tr命令还可以用来删除多余的换行符:
```
cat multi_blanks.txt|tr -s '\n'
line1
line2
line3
line4
```
上面的例子展示了如何使用tr删除多余的'\n'字符。接下来让我们用一种巧妙的方式将数字列表进行相加
```
seq 10|echo $[ `tr '\n' '+' ` 0]
```
##### 字符类
tr可以将不同的字符类作为集合使用，所支持的字符类如下所示。
* alnum: 字母和数字
* alpha: 字母
* cntrl: 控制(非打印)字符
* digit: 数字
* graph: 图形字符
* lower: 小写字母
* print: 可打印字符
* punct: 标点符号
* space: 空白字符
* upper: 大写字母
* xdigit: 十六进制字符
可以按照下面的方式选择所需的字符类    
tr [:class:] [:class:]
