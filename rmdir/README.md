# rmdir

#### summary
Rmidr has the same rmdir command as ubuntu 18.04

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/rmdir/rmdir
```

#### usage
```console
Usage of rmdir:
  -ignore-fail-on-non-empty
    	ignore each failure that is solely because a directory
    	is non-empty
  -p, --parent
    	remove DIRECTORY and its ancestors; e.g., 'rmdir -p a/b/c' is
    	similar to 'rmdir a/b/c a/b a'
```
