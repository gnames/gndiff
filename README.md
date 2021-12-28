# `gndiff` app takes two files with scientific names, compares them and returns the result.

## Introduction

It is often useful to compare one check-list to another. This project follows
pretty much the same algorithms as [GNverifier], but does not require an
external database and can be used off-line.

## Usage

If you need to compare a list of names with a data-set that is imported as a
[GNverifier data-source] use [GNverifier] instead like this:

```bash
gnvrifier names.txt -o -s 10
```

where `-s` flag provides Id of required [GNverifier data-source]

### Compare Files

Prepare two files with names. There are 3 possible file formats:

* A simple lists of scientific names, one name per line
* Comma-separated or Tab-separated (CSV/TSV) file with a `ScientificName`
  field. Fields `TaxonID` and `Family` are also ingested, if given, any
  capitalization of the fields names is accepted.

The first of the two files should contain names that need to be matched.
The second file should contain reference names.

Run command:

```bash
gndiff source.csv reference.csv
```

### Options and flags

According to POSIX standard flags and options can be given either before or
after name-string or file name.

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


[GNverifier]: https://github.com/gnames/gnverifier
[GNverifier data-source]: https://verifier.globalnames.org/data_sources
