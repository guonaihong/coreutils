# mv

#### 简介
mv

#### install
```
env GOPATH=`pwd` go get -u github.com/guonaihong/coreutils/mv
```

####

#### Examples

```bash
mv myfile.txt myfiles
```
Move the file myfile.txt into the directory myfiles. If myfiles is a file, it will be overwritten. If the file is marked as read-only, but you own the file, you will be prompted before overwriting it.

```bash
mv myfiles myfiles2
```
If myfiles is a file or directory, and myfiles2 is a directory, move myfiles into myfiles2. If myfiles2 does not exist, the file or directory myfiles is renamed myfiles2.

```bash
mv myfile.txt ../
```
Move the file myfile.txt into the parent directory of the current directory.

```bash
mv -t myfiles myfile1 myfile2
```
Move the files myfile1 and myfile2 into the directory myfiles.

```bash
mv myfile1 myfile2 myfiles
```
Same as the previous command.

```bash
mv -n file file2
```

If file2 exists and is a directory, file is moved into it. If file2 does not exist, file is renamed file2. If file2 exists and is a file, nothing happens.
```bash
mv -f file file2
```
If file2 exists and is a file, it will be overwritten.

```bash
mv -i file file2
```

If file2 exists and is a file, a prompt is given:

mv: overwrite 'file2'?
Entering "y", "yes", "Yes", or "Y" will result in the file being overwritten. Any other input will skip the file.

```bash
mv -fi file file2
```
Same as mv -i. Prompt before overwriting. The f option is ignored.

```bash
mv -if file file2
```

Same as mv -f. Overwrite with no prompt. the i option is ignored.
```bash
mv My\ file.txt My\ file\ 2.txt
```

Rename the file "My file.txt" to "My file 2.txt". Here, the spaces in the file name are escaped, protecting them from being interpreted as part of the command.

```bash
mv "My file.txt" "My file 2.txt"

```
Same as the previous command.

```bash
mv "My file.txt" myfiles
```

The result of this command:

If myfiles a directory, My file.txt is moved into myfiles.
If myfiles a file, My file.txt is renamed myfiles, and the original myfiles is overwritten.
If myfiles does not exist, My file.txt is renamed myfiles.

```bash
mv My*.txt myfiles
```
Here, * is a wildcard meaning "any number, including zero, of any character."

If myfiles is a directory: all files with the extension .txt, whose name begins with My, will be moved into myfiles.
If myfiles does not exist or is not a directory, mv reports an error and does nothing.

```bash
my My\ file??.txt myfiles`
```
Here, ? is a wildcard that means "zero or one of any character." It's used twice, so it can match a maximum of two characters.

If myfiles is a directory: any file with zero, one, or two characters between My file and .txt in their name is moved into myfiles.
If myfiles doesn't exist, or is not a directory, mv reports an error and does nothing.

##### Making backups
```bash
mv -b file file2
```

If file2 exists, it will be renamed to file2~.

```bash
mv -b --suffix=.bak file file2
```
If file2 exists, it will be renamed to file2.bak.

```bash
mv --backup=numbered; mv file file2
```
If file2 exists, it will be renamed file2.~1~. If file2.~1~ exists, it will be renamed file2.~2~, etc.

```bash
VERSION_CONTROL=numbered mv -b file file2
```
Same as previous command. The environment variable is defined for this command only.

```bash
export VERSION_CONTROL=numbered; mv -b file file2
```
By exporting the VERSION_CONTROL environment variable, all mv -b commands for the current session will use numbered backups.

```bash
export VERSION_CONTROL=numbered; mv file file2
```
Even though the VERSION_CONTROL variable is set, no backups are created because -b (or --backup) was not specified. If file2 exists, it is overwritten.

* Example comes from(https://www.computerhope.com/unix/umv.htm)
