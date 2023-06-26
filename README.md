# binutil

[![Go Reference](https://pkg.go.dev/badge/github.com/bbengfort/binutil.svg)](https://pkg.go.dev/github.com/bbengfort/binutil)
[![Tests](https://github.com/bbengfort/binutil/actions/workflows/test.yaml/badge.svg)](https://github.com/bbengfort/binutil/actions/workflows/test.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/bbengfort/binutil)](https://goreportcard.com/report/github.com/bbengfort/binutil)

**Helpers for converting to and from binary and string representations of data**

## Installation

The simplest way to install the `binutil` command is to use [Homebrew](https://github.com/rotationalio/homebrew-tools) on Mac OS X:

```
$ brew install rotationalio/tools/binutil
```

If you have Go v1.19 or later, you can install it directly to your `$GOBIN` as follows:

```
$ go install github.com/bbengfort/binutil/cmd/binutil@latest
```

Finally, you can manually download the appropriate binary from the [releases page](https://github.com/bbengfort/binutil/releases) and extract to a directory in your `$PATH` such as `~/bin`.

Once installed, make sure you can execute the `binutil` command as follows:

```
$ binutil --version
```

## Usage

You can use `binutil` to encode and decode binary strings from the command line. Basic usage is as follows:

```
$ binutil -d [DECODER] -e [ENCODER] [ENCODED BINARY DATA]
```

For example, to convert base64 encoded data to hex encoded data:

```
$ binutil -d base64 -e hex nYOG+hjpWFAGJsaBoLrSeg==
9d8386fa18e958500626c681a0bad27a
```

To see a list of availabled decoders, use the `binutil decoders` command:

```
$ binutil decoders
Registered Decoders:
====================
- ascii
- base64
- base64-raw
- base64-rawurl
- base64-std
- base64-url
- hex
- latin1
- text
- ulid
- utf-8
- uuid
```

**Note:** the above list may vary with new releases of `binutil`!

### Generating Random Data

You can also quickly generate random data with the `binutil rand` command:

```
$ binutil rand -s 42 -e base64
cdMvYhnbiMymr0zTFWpOcEnnhSQhFe4a+FvJu0p3aBYuAiVvqr8crwx7
```

The `-s` specifies the size of the generated data by number of bytes.

This is useful for generating data for tests quickly.

### Generating ULID and UUIDs

One important use case for `binutil` is to manage [ULID](https://github.com/oklog/ulid) and [UUID](https://github.com/google/uuid) in different binary formats (e.g. to insert data into a database via a SQL command). You can generate ULIDs and UUIDs as follows:

```
$ binutil ulid
       ULID        01H3W76DQS0X8PX8ZT62MYAMD1
       Time     2023-06-26T11:04:40.313-05:00
  Hex Bytes  0188f87336f907516ea3fa30a9e551a1
  b64 Bytes          AYj4czb5B1Fuo/owqeVRoQ==
```

The output shows a variety of formats and the timestamp for ULIDs, allowing easy generation of conflict-free IDs for tests and fixtures. To easily copy and paste a UUID:

```
$ binutil uuid -e uuid -n | pbcopy
```