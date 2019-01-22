# sort

#### 简介
sort

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/tr
```

#### Example
(下面的示例来自linux shell scripting cookbook)

对一组文件进行排序
```
sort file1.txt file2.txt file3.txt >sorted.txt
sort file1.txt file2.txt file3.txt  -o sorted.txt
```

按照数字顺序进行排序
```
sort -n file.txt
```

按照逆序进行排序
```
sort -r file.txt
```

按照月份进行排序(依照一月，二月，三月)
```
sort -M months.txt
```

合并两个已排序过的文件
```
sort -m sorted1 sorted2
```

检查文件是否排过序
```
sort -C filename
```

```
cat data.txt
1 mac 2000
2 winxp 4000
3 bsd 1000
4 linux 1000
```
有很多方法可以对这段文本排序。目前它是按照序号(第一列)来排序的。我们也可以依据第二列和第三列来排序。
-k指定了排序所依据的字符。如果是单个数字，则指的是列号。-r告诉sort命令按照逆序进行排序。例如
```
#依据第1列，以逆序形式排序
sort -nrk 1 data.txt
4 linux 1000
3 bsd 1000
2 winxp 4000
1 mac 2000

# -nr 表明按照数字顺序，采用逆序形式排序
# 依据第2列进行排序

sort -k 2 data.txt
3 bsd 1000
4 linux 1000
1 mac 2000
2 winxp 4000
```

-k后的整数指定了文本文件中的某一列。列与列之间有空格分隔。如果需要将特定范围内的一组字符（例如，第2列中的第4～5个字符）作为键，应该使用由
点号分隔的两个整数来定义一个字符位置，然后将该范围内的第一个字符和最后一个字符用逗号连接起来:
```
cat data.txt
1 alpha 300
2 beta 200
3 gamma 100

sort -bk 2.3,2.4 data

3 gamma 100
1 alpha 300
2 beta 200
```

把作为排序依据的字符写成数值键。为了提取出这些字符，用其在行内的起止位置作为键的书写格式（在上面的例子中，起止位置是2和3）。
用第一个字符作为键。
```
sort -nk 1,1 dat.txt
```
为了使sort的输出与以\0作为终止符的xargs命令相兼容，采用下面的命令
```
sort -z data.txt|xargs -0
#终止符\0用来确保安全地使用xargs命令
```
有时文本中可能会包含一些像空格之类的多余字符。如果需要忽略标点符号并以字典序排序，可以使用:
sort -bd unsorted.txt
其中，选项-b用于忽略文件中的前导空白行，选项-d用于指明以字典序进行排序。
```
