# help-bash-go

Clone of [help-bash](https://github.com/berquerant/help-bash) written in Go.

```
❯ ./help-bash -h
help-bash displays the things like documentation comments in .sh file.

USAGE:
  help-bash [OPTIONS] [PATH]


ARGS:
    <PATH>
        A file to display documentations.


ENVIRONMENT VARIABLES:
     HELP_BASH_DEBUG
         If value is 1 then enables debug logs.
         If value is 2 then enables trace logs in addition to debug logs.


OPTIONS:
  -f    Displays top-level (outside of any statements) functions documentations.
  -r    Displays top-level variables documentations.
  -t uint
        Number of rows of which threshold to determine whether to display the top-level comments.
        e.g. the value is 3 then displays 3 or more lines of top-level comments. (default 3)
  -v    Prints version information.
```
