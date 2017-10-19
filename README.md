# Development Setup

By default, `go get` does not work for private GitHub repos. This will fix that.
1. `git config --global url."git@github.com:".insteadOf "https://github.com/"`
2. `go get github.com/8thlight/vulcanizedb`

# Running the Tests

`go test`
