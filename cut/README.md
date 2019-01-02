# cut

#### 简介
cut

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/cut
```

#### 命令行选项
```console
```

#### Example
(以下内容来自linux shell Scripting Cookbook)

这条命令将显示第2列和第3列
```shell
cut -f 2,3 filename
```

cut也能从stdint中读取输入文本。
\t是字段或列的默认分割符。对于分割符的行。会将该行照原样打印出来。如果不想打印出这种不包含分割符的行，
则可以使用cut的 -s选项。

我们也可以使用--complement选项对提取的字段进行补集运算。假设有多个字段，
你希望打印出除第3列之外的所有列，则可以使用:
```shell
cut -f3 --complement student_data
```

要指定分割符，使用-d选项
```
cat delimited_data.txt
No;Name;Mark;Percent
1;Sarath;45;90

cut -f 2 -d ';' delimited_data.txt
```

假设我们不依赖分割符，但需要通过将字段定义为一个字符范围来进行字段提取
```
N- 从第N个字节，字符或字段到行尾
N-M 从第N个字节，字符或字段到第M个(包括第M个在内)字节、字符或字段
-M 第1个字节，字符或字段到第M个(包括第M个在内)字节、字符或字段
```

打印第1到5个字符
```
cut  -c1-5 filename
```

打印前2个字符
```
cut filename -c -2
```
