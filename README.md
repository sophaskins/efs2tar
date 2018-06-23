# efs2tar

`efs2tar` is a tool that converts SGI EFS-formatted filesystem images (ie, the result of `dd`-ing a whole device in to a file) in to tarballs. It was based entirely on NetBSD's `sys/fs/efs` ([source](http://cvsweb.netbsd.org/bsdweb.cgi/src/sys/fs/efs/?only_with_tag=MAIN)).

## Example usage

```
$ go install github.com/sophaskins/efs2tar
$ efs2tar -in ~/my-sgi-disc.iso -out ~/my-sgi-disc.tar
```

The Internet Archive has [several discs](https://archive.org/search.php?query=sgi&and%5B%5D=mediatype%3A%22software%22&page=2) in its collections that are formatted with EFS.


## "Edge cases" not covered
* any type of file other than directories and normal files (which is to say, links in particular do not work)
* partition layouts other than what you'd expect to see on an SGI-produced CDROM
* any sort of error handling...at all
* that includes verifying magic numbers
* preserving the original file permissions
* I've only tested this on like, one CD