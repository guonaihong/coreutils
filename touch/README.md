# touch

#### summary
touch has the same touch command as ubuntu 18.04

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/touch/touch
```

#### usage
```console
Usage of touch:
  -a	change only the access time
  -c, --no-create
    	do not create any files
  -d, --date string
    	parse STRING and use it instead of current time
  -debug
    	debug mode
  -h, --no-dereference
    	affect each symbolic link instead of any referenced
    	file (useful only on systems that can change the
    	timestamps of a symlink)
  -m	change only the modification time
  -r, --referenced string
    	use this file's times instead of current time
  -t string
    	use [[CC]YY]MMDDhhmm[.ss] instead of current time
  -time string
    	change the specified time:
    	WORD is access, atime, or use: equivalent to -a
    	WORD is modify or mtime: equivalent to -m
```

#### todo
gnu touch -d  选项是用yacc实现，如果要复该投入时间巨大，后面实现
