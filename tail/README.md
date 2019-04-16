# tail

#### summary
tail

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/tail/tail
```

#### usage
``` console
Usage of tail:
  -F	same as --follow=name --retry
  -^\d+$, --n, --lines string
    	output appended data as the file grows;
    	an absent option argument means 'descriptor'
  -c, --bytes string
    	output the last NUM bytes; or use -c +NUM to
    	output starting with byte NUM of each file
  -f	output appended data as the file grows;
    	an absent option argument means 'descriptor'
  -f, --follow string
    	output appended data as the file grows;
    	an absent option argument means 'descriptor'
  -max-unchanged-stats string
    	with --follow=name, reopen a FILE which has not
    	changed size after N (default 5) iterations
    	to see if it has been unlinked or renamed
    	(this is the usual case of rotated log files);
    	with inotify, this option is rarely useful
  -pid string
    	with -f, terminate after process ID, PID dies
  -q, --quiet, --silent
    	never print headers giving file names
  -retry
    	keep trying to open a file if it is inaccessible
  -s, --sleep-interval float
    	with -f, sleep for approximately N seconds
    	(default 1.0) between iterations;
    	with inotify and --pid=P, check process P at
    	least once every N seconds
  -v, --verbose
    	always print headers giving file names
  -z, --zero-terminated
    	line delimiter is NUL, not newline

```
#### 如下一些频率不高的选项后面实现 
```console
--pid todo
--retry todo
--max-unchanged-status todo
--follow todo
```
