# chgrp

#### summary
chgrp has the same chgrp command as ubuntu 18.04

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/chgrp/chgrp
```

#### usage
```console
Usage of chgrp:
  -H	if a command line argument is a symbolic link
    	to a directory, traverse it
  -L	traverse every symbolic link to a directory
    	encountered
  -P	do not traverse any symbolic links (default)
  -R, --recursive
    	operate on files and directories recursively
  -c, --changes
    	like verbose but report only when a change is made
  -dereference
    	affect the referent of each symbolic link (this is
    	the default), rather than the symbolic link itself
  -f, --quiet, --silent
    	suppress most error messages
  -from string
    	change the owner and/or group of each file only if
    	its current owner and/or group match those specified
    	here.  Either may be omitted, in which case a match
    	is not required for the omitted attribute
  -h, --no-dereference
    	affect symbolic links instead of any referenced file
    	(useful only on systems that can change the
    	ownership of a symlink)
  -no-preserve-root
    	do not treat '/' specially (the default)
  -preserve-root
    	fail to operate recursively on '/'
  -reference string
    	use RFILE's owner and group rather than
    	specifying OWNER:GROUP values
  -v, --verbose
    	output a diagnostic for every file processed
```
