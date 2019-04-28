# uname

#### summary
uname has the same uname command as ubuntu 18.04

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/uname/uname
```

#### usage
```console
Usage of uname:
  -a, --all
    	print all information, in the following order,
    	except omit -p and -i if unknown:
  -i, --hardware-platform
    	print the hardware platform (non-portable)
  -m, --machine
    	print the machine hardware name
  -n, --nodename
    	print the network node hostname
  -o, --operating-system
    	print the operating system
  -p, --processor
    	print the processor type (non-portable)
  -r, --kernel-release
    	print the kernel release
  -s, --kernel-name
    	print the kernel name
  -v, --kernel-version
    	print the kernel version
```
