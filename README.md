# `GNdiff` app takes two files with scientific names and compares them

<!-- vim-markdown-toc GFM -->

* [Introduction](#introduction)
* [GNdiff Installation](#gndiff-installation)
  * [Using Homebrew on Mac OS X, Linux, and Linux on Windows (WSL2)](#using-homebrew-on-mac-os-x-linux-and-linux-on-windows-wsl2)
  * [Linux and Mac without Homebrew](#linux-and-mac-without-homebrew)
  * [MS Windows](#ms-windows)
  * [Compile from source](#compile-from-source)
* [Usage](#usage)
  * [Compare Files](#compare-files)
  * [Compare using test files](#compare-using-test-files)
  * [Options and flags](#options-and-flags)
    * [help](#help)
    * [version](#version)
    * [format](#format)
    * [quiet](#quiet)
* [Family names as a disambiguation tool](#family-names-as-a-disambiguation-tool)
  * [The same nomenclatural code homonyms (senior, junior homonyms)](#the-same-nomenclatural-code-homonyms-senior-junior-homonyms)
  * [Homonyms from different nomenclatural codes (hemihomonyms)](#homonyms-from-different-nomenclatural-codes-hemihomonyms)

<!-- vim-markdown-toc -->

## Introduction

It is often useful to compare one checklist to another. This project uses
pretty much the same algorithms as [GNverifier], but does not require an
external database and can be used offline.

[GNdiff] is a complementary tool to [GNverifier]. It is made to compare a
checklist with checklists that are not in [GNverifier database][GNverifier
data-sources]. If you need to compare a list of names that are [already in
GNverifier][GNverifier data-sources], either use [GNverifier web app] or
[install][GNverifier Install] it locally and run:

```bash
gnvrifier names.txt -o -s 12,1
```

where `-s` option provides `id`/`ids` of selected [GNverifier data-sources],
and `-o` flag (`--only-preferred`) limits results to data-sources set by `-s`
option.

If both checklists of scientific names are local, use [GNdiff].

## `GNdiff` Installation

### Using Homebrew on Mac OS X, Linux, and Linux on Windows ([WSL2][WSL install])

[Homebrew] is a popular package manager for Open Source software originally
developed for Mac OS X. Now it is also available on Linux, and can easily
be used on MS Windows 10 or 11, if Windows Subsystem for Linux (WSL) is
[installed][WSL install].

Note that [Homebrew] requires some other programs to be installed, like Curl,
Git, a compiler (GCC compiler on Linux, Xcode on Mac). If it is too much,
go to the `Linux and Mac without Homebrew` section.

To use `GNdiff` with Homebrew:

1. Install [Homebrew]

2. Open terminal and run the following commands:

```bash
brew tap gnames/gn
brew install gndiff
```

### Linux and Mac without Homebrew

Download the latest [GNdiff release], untar, and install binary somewhere
in your path.

```bash
tar xvf gndiff-linux-0.1.1.tar.gz
# or tar xvf gndiff-mac-0.1.1.tar.gz
sudo mv gndiff /usr/local/bin
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

### Compile from source

Install Go according to [installation instructions][go-install] and run:

```bash
go install github.com/gnames/gndiff/gndiff@latest
```

## Usage

### Compare Files

Prepare two files with names. There are 3 possible file formats:

* A simple list of scientific names, one name per line.
* Comma-separated or Tab-separated (CSV) file with a `ScientificName`
  field.
* Tab-separated (TCV) file with a `ScientificName` field.

For both CSV and TSV files the fields `TaxonID` and `Family` are also ingested
if given, any capitalization of either of the fields names is accepted.

The `Family` field indicates a family a particular species are assigned to
according to the dataset. Normally this field is not needed, but in case of
tricky homonyms it helps to resolve taxa from each other.

Run command:

```bash
gndiff query.csv reference.csv
```

The first of the two files should contain names that need to be matched.
The second file should contain reference names.

Any combination of these 3 formats would work:

```bash
gndiff file.csv names.txt
gndiff file.tsv file.csv
gndiff name.txt file.tsv
# etc.
```

### Compare using test files

To see how it works you can use tests files from the GNdiff project for example
[ebird.csv] and [ioc-bird.csv]. These files contain avian names from [eBird]
and [eBird] and [IOC Bird] checklists correspondingly.
Open each link and use `Ctrl-S` on Windows/Linux or `⌘-S` on Mac to save
these files on your computer. You can then run:

```bash
gndiff path/to/ebird.csv path/to/ioc-bird.csv
# or
gndiff path/to/ioc-bird.csv path/to/ebird.csv
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
gndiff query.txt ref.txt -f pretty
# or
gndiff query.txt ref.txt --format=pretty
```

#### port (integer)

When `port` is set, `GNdiff` works as a web server with its RESTful API
exposed at the given port.

```bash
gndiff -p 8080
# or
gndiff --port 8080
```

#### quiet

This flag supresses warnings log, showing only the matching results.

```bash
gndiff query.txt ref.txt -q
# or
gndiff query.txt ref.txt --quiet
```

Please note, that matching result uses `STDOUT`, while log uses `STDERR`,
so a similar result can be achieved by redirecting `STDERR` to `/dev/null`

```bash
gndiff query.txt ref.txt 2> /dev/null
```

## Family names as a disambiguation tool

Family sometimes help to distinquish homonyms in names lists. For example,
there are homonyms within one nomenclatural code, and homonyms in between
two nomenclatural codes.

### The same nomenclatural code homonyms (senior, junior homonyms)

For example, in zoology a genus name `Echidna` has 3 homonyms:

1. Moray eel

    Family `Muraenidae` -> Genus `Echidna J. R. Forster`

2. Egg-laying mammal Echidna

    Family `Tachyglossidae` -> `Echidna Cuvier, 1797` (junior homonym)

    Currently genus `Tachyglossus Illiger, 1811`

3. Snake

    Family `Viperidae` -> `Echidna Merrem, 1820` (junior homonym)

    Currently genus `Bitis Gray, 1842`

Such homonyms are not allowed within the same code and eventually they
get corrected. However historical records still contain them and
have to be disambiguated.

### Homonyms from different nomenclatural codes (hemihomonyms)

There are no rules how to deal with homonyms that are treated by different (for
example Botanical and Zoological) nomenclatural codes.

1. Sea Snail (Zoological Nomenclatural Code)

    Family `Ficidae` -> `Ficus variegata Röding, 1798`

2. Red fig (Botanical Nomenclatural Code)

    Family `Moraceae` -> `Ficus variegata Blume`

[GNdiff release]: https://github.com/gnames/gndiff/releases/latest
[GNdiff]: https://github.com/gnames/gndiff
[GNverifier data-sources]: https://verifier.globalnames.org/data_sources
[GNverifier install]: https://github.com/gnames/gnverifier#installation
[GNverifier]: https://github.com/gnames/gnverifier
[GNverifier web app]: https://verifier.globalnames.org
[Homebrew]: https://brew.sh/
[WSL install]: https://docs.microsoft.com/en-us/windows/wsl/install-win10
[go-install]: https://golang.org/doc/install
[winpath]: https://www.computerhope.com/issues/ch000549.htm
[ebird.csv]: https://raw.githubusercontent.com/gnames/gndiff/master/testdata/ebird.csv
[ioc-bird.csv]: https://raw.githubusercontent.com/gnames/gndiff/master/testdata/ioc-bird.csv
[eBird]: https://www.worldbirdnames.org/
[IOC Bird]: https://www.worldbirdnames.org/
