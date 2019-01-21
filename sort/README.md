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
