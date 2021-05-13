# Grep

We all know what grep is, come on.
This is a task I had to do as part of my Go training.

## Installation

Just build it.

```bash
go build ./grep.go
```

## Usage

```bash
Usage of ./grep:
  -c    Colorize matches.
  -f string
        File path, mutually exclusive with -t. If both are empty, pipe is used.
  -p string
        Pattern to match - mandatory.
  -r    Enables regexp matching, pattern must be a valid regular expression.
  -t string
        Text string, mutually exclusive with -f. If both are empty, pipe is used.
```

### TODO

This can be extended and optimized in so many ways.. Just needed to hand it in already.