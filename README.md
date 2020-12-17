# re
`re` is a super simple recursive search and replace tool. It traverses the current 
directory tree (including its subfolders) and replaces a search string in all files that
it finds it in.

You can include and exclude files by extension. 

## Installation
Binaries can be found on the [releases page](https://github.com/binwiederhier/re/releases). 
Alternatively, you can build it yourself by running `go build` (or `make`).

## Usage
```
Syntax: re [options] SEARCH REPLACEMENT [DIR ...]

Options:
  -e string
    	Comma-separated list of excluded files, wildcards supported (default ".bzr,CVS,.git,.hg,.svn")
  -f	Apply changes
  -i string
    	Comma-separated list of included files, wildcards supported, e.g "*.js,*.html,*index.*"
``` 

## Examples

```bash
# Shows files that would have been replaced
$ re CacheDir ClipboardDir

# Replace CacheDir with ClipboardDir recursively
$ re -f CacheDir ClipboardDir

# Replace "Foo(Bar)" with "Foo(Baz)" recursively in all scala files
$ re -f -i "*.scala" "Foo(Bar)" "Foo(Baz)"  
```

## License 
Made by [Philipp Heckel](https://heckel.io), distributed under the [Apache License 2.0](LICENSE).