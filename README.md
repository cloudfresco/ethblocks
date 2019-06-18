[![Build Status](https://travis-ci.org/cloudfresco/ethblocks.svg?branch=master)](https://travis-ci.org/cloudfresco/ethblocks) [![Go Report Card](https://goreportcard.com/badge/github.com/cloudfresco/ethblocks)](https://goreportcard.com/report/github.com/cloudfresco/ethblocks)


# Ethblocks -- Analyze Ethereum Blocks

## Features

	* Lightweight, easy to use Go library for getting blocks and transactions
	* Save blocks and transactions to MySQL (Postgres to come)
	* Uses structs/functions from go-ethereum as much as possible

## Requirements
We use the latest version of Go.

## Usage

Import the library:

```go
import "github.com/cloudfresco/ethblocks/svc"
```

Then call any of the functions.  See [examples](https://github.com/cloudfresco/ethblocks/tree/master/examples).

## License

The Ethblocks source files are distributed under the [MIT License](https://github.com/cloudfresco/ethblocks/blob/master/LICENSE).
