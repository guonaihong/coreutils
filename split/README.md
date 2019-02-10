# split

#### 简介
split

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/split
```

#### Example
(下面的示例来自linux shell scripting cookbook) 和本人构造的测试用例

有时候必须把文件分割成多个更小的片段。比如把大的pcm音频文件切成一堆小的音频文件。
split命令可以用来分割文件。该命令接受文件名作为参数，然后创建出一系列体积更小的文件，
其中依据字母序排在首位的那个文件对应于原始文件的第一部分，排在次位的文件对应于原始文件
的第二部分，以此类推。
例如，通过指定分割大小，可以将100KB的文件分成一系列10KB的小文件。在split命令中，除了k(KB)
我们还可以使用M(MB)，G(GB), c(byte), w(word)

```
split -b 10k data.file
ls
data.file xaa xab xac xad xae xaf xag xah xai xaj
```

上而后命令将data.file分割成了10个大小为10KB的文件。这些新文件以xab，xac，xad的形式命名。
split默认使用字母后缀。如果想使用数字后缀，需要使用-d选项。此外，-a length可以指定后缀长度:
```
split -b 10k data.file -d -a 4
```

```
ls
data.file x0009 x0019 x0029 x0039 x0049 x0059 x0069 x0079
```

为分割后的文件指定文件名前缀
```
split -b 10k data.file -d -a 4 split_file

data.file split_file0002 split_file003
```

如果不想按照数据块大小，而是根据行数来分割文件的话，可以使用-l no_of_lines:
```
split -l 10 data.file
```
