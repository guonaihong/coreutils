# join

#### 简介
join

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/join
```

#### example
(下面的)
(下面的内容来自网络)
join用于连接公共字段的两个文件的行。默认连接字段是由空格分隔的第一个字段。
请看下面两个文件foodtypes.txt和foods.txt
```
cat foodtypes.txt
1 Protein
2 Carbohydrate
3 Fat

cat foods.txt
1 Cheese 
2 Potato
3 Butter
```

这两个文件共享第一个字段，所以可以连接
```
join foodtypes foods.txt
1 Protein Cheese
2 Carbohydrate Potato
3 Fat Butter
```
要使用不同的字段连接文件，可以传递-1和-2选项以进行连接。 在以下示例中，有两个文件wine.txt和reviews.txt。
```
cat wine.txt
Red Beaunes France
White Reisling Germany
Red Riocha Spain

cat reviews.txt
Beaunes Great!
Reisling Terrible!
Riocha Meh
```
可以通过指定应该用于加入文件的字段来连接这些文件。 这两个文件的共同点是葡萄酒的名称。 
在wine.txt中，这是第二个字段。 在reviews.txt中，这是第一个字段。 
通过指定这些字段，可以使用-1和-2连接文件。

在这些文件上运行联接会导致错误，因为文件未排序。
```
join -1 2 -2 1 wine.txt reviews.txt
join: wine.txt:3: is not sorted: Red Beaunes France
join: reviews.txt:2: is not sorted: Beaunes Great!
Riocha Red Spain Meh
Beaunes Red France Great!
```
用sort命令排序之后再join
```
join -1 2 -2 1 <(sort -k 2 wine.txt) <(sort reviews.txt)
Beaunes Red France Great!
Reisling White Germany Terrible!
Riocha Red Spain Meh
```

如何指定分割符
要使用join指定分割符，可以使用-t选项。特别是join csv文件
```
cat names.csv
1,John Smith,London
2,Arthur Dent, Newcastle
3,Sophie Smith,London

cat transactions.csv
£1234,Deposit,John Smith
£4534,Withdrawal,Arthur Dent
£4675,Deposit,Sophie Smith
```
使用-t指定分割符
```
join -1 2 -2 3 -t , names.csv transactions.csv
John Smith,1,London,£1234,Deposit
Arthur Dent,2, Newcastle,£4534,Withdrawal
Sophie Smith,3,London,£4675,Deposit
```
```
259/5000
如何指定输出格式
要指定join的输出格式，请使用-o选项。 这允许定义将在输出中显示的字段的顺序，或仅显示某些字段。

在前面的例子中输出我们如下。
```
John Smith,1,London,£1234,Deposit
```
要指定顺序，将字段列表传递给-o。 对于这个例子，这是-o 1.1,1.2,1.3,2.2,2.1。 这将按所需顺序格式化输出。
```
join -1 2 -2 3 -t , -o 1.1,1.2,1.3,2.2,2.1 names.csv transactions.csv
1,John Smith,London,Deposit,£1234
2,Arthur Dent, Newcastle,Withdrawal,£4534
3,Sophie Smith,London,Deposit,£4675
```
#####参考资料
* https://shapeshed.com/unix-join/
