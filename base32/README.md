# base32

#### summary
base32 has the same base32 command as ubuntu 18.04

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/base32/base32
```

#### usage
```console
Usage of base32:
  -d, --decode
    	decode data
  -i, --ignore-garbage
    	when decoding, ignore non-alphabet characters
  -w, --wrap int
    	wrap encoded lines after COLS character (default 76).
    	Use 0 to disable line wrapping
```
