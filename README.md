# `GNdiff` app takes two files with scientific names and compares them.

## Introduction

It is often useful to compare one checklist to another. This project uses
pretty much the same algorithms as [GNverifier], but does not require an
external database and can be used offline.

## Installation

### Using Homebrew on Mac OS X, Linux, and Linux on Windows ([WSL2])

Homebrew is a popular package manager for Open Source software originally
developed for Mac OS X. Now it is also available on Linux, and can easily
be used on Windows 10, if Windows Subsystem for Linux (WSL) is
[installed][WSL install].

To use `GNdiff` with Homebrew:

1. Install [Homebrew]

2. Open terminal and run the following commands:

```bash
brew tap gnames/gn
brew install gndiff
```

### MS Windows

Download the latest [GNdiff release], unzip.

One possible way would be to create a default folder for executables and place
``GNdiff`` there.

Use ``Windows+R`` keys
combination and type "``cmd``". In the appeared terminal window type:

```cmd
mkdir C:\Users\your_username\bin
copy path_to\gndiff.exe C:\Users\your_username\bin
```

[Add ``C:\Users\your_username\bin`` directory to your ``PATH``][winpath]
environment variable.

Another, simpler way, would be to use ``cd C:\Users\your_username\bin`` command
in ``cmd`` terminal window. The ``GNdiff`` program then will be automatically
found by Windows operating system when you run its commands from that
directory.

### Linux and Mac

Download the latest [GNdiff release], untar, and install binary somewhere
in your path.

```bash
tar xvf gndiff-linux-0.1.1.tar.gz
# or tar xvf gndiff-mac-0.1.1.tar.gz
sudo mv gndiff /usr/local/bin
```

### Compile from source

Install Go according to [installation instructions][go-install]

```bash
go get github.com/gnames/gndiff/gndiff
```

## Usage

If you need to compare a list of names with a data-set that exists as a
[GNverifier data-source] use [GNverifier]:

```bash
gnvrifier names.txt -o -s 12
```

where `-s` option provides Id of required [GNverifier data-source], and `-o`
flag limits restuls to the selected data-source.

### Compare Files

Prepare two files with names. There are 3 possible file formats:

* A simple list of scientific names, one name per line.
* Comma-separated or Tab-separated (CSV/TSV) file with a `ScientificName`
  field. Fields `TaxonID` and `Family` are also ingested if given, any
  capitalization of the fields names is accepted.

The first of the two files should contain names that need to be matched.
The second file should contain reference names.

Run command:

```bash
gndiff source.csv reference.csv
```

### Options and flags

According to POSIX standard flags and options can be given either before or
after file-paths arguments.

#### help

```bash
gndiff -h
# or
gndiff --help
# or
gndiff
```

#### version

```bash
gndiff -V
# or
gndiff --version
```

#### format

Sets the format of the comparison result and can take the following values:

* `csv`: Comma-separated format
* `tsv`: Tab-separated format
* `compact`: JSON as one line
* `pretty`: JSON in a human-readable format with indentations and lines separation.

The default format is CSV.

```bash
gndiff source.txt ref.txt -f pretty
# or
gndiff source.txt ref.txt --format=pretty
```

#### quiet

This flag supresses warnings log, showing only the matching results.

```bash
gndiff source.txt ref.txt -q
# or
gndiff source.txt ref.txt --quiet
```

Please note, that matching result uses `STDOUT`, while log uses `STDERR`,
so a similar result can be achieved by redirecting `STDERR` to `/dev/null`

```bash
gndiff source.txt ref.txt 2> /dev/null
```


[GNdiff release]: https://github.com/gnames/gndiff/releases/latest
[GNdiff]: https://github.com/gnames/gndiff
[GNverifier data-source]: https://verifier.globalnames.org/data_sources
[GNverifier]: https://github.com/gnames/gnverifier
[Homebrew]: https://brew.sh/
[WSL install]: https://docs.microsoft.com/en-us/windows/wsl/install-win10
[go-install]: https://golang.org/doc/install
[winpath]: https://www.computerhope.com/issues/ch000549.htm
