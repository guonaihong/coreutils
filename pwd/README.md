# pwd

#### summary
pwd has the same pwd command as ubuntu 18.04

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/pwd/pwd
```

#### usage
```console
Usage of pwd:
  -L, --logical
    	print the value of $PWD if it names the current working
    	directory
  -P, --physical
    	print the physical directory, without any symbolic links
```
