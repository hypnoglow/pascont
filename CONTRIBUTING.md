# Contributing to pascont

## Run tests manually

Until Go 1.9, `go test ./...` includes `vendor` directory. So, to exclude vendors, run:

    go test $(go list ./... | grep -v '/vendor/')


