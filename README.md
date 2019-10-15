# gofixerio &mdash; Simple wrapper for https://fixer.io API

## Installation

```bash
go get https://github.com/LordotU/go-fixerio
```

## Testing

```bash
FIXERIO_API_KEY="your API key here" go test
```

**Note:** all tests are written for *free subscription plan*!

## Usage

The simplest example is:

```go
	import (
		"log"

		"https://github.com/LordotU/go-fixerio"
	)

	func main() {
		fixerio := gofixerio.New("your API key here", "EUR", false)

		latestRates, err := fixerio.GetLatest()
		if err != nil {
			log.Panic(err)
		}

		log.Printf("%+v", latestRates)
	}
```

## Docs

https://godoc.org/github.com/LordotU/go-fixerio
