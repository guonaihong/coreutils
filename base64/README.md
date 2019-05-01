# base64

#### summary
base64 has the same base64 command as ubuntu 18.04

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/base64/base64
```

#### usage
```console
Usage of base64:
  -d, --decode
    	decode data
  -i, --ignore-garbage
    	when decoding, ignore non-alphabet characters
  -w, --wrap int
    	wrap encoded lines after COLS character (default 76).
    	Use 0 to disable line wrapping
```
