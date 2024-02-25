
# Getting started

## Install Golang

Make sure you have Go 1.13 or higher installed.

https://golang.org/doc/install

## Environment Config

Set-up the standard Go environment variables according to latest guidance (see https://golang.org/doc/install#install).


## Install Dependencies
From the project root, run:
```
go build ./...
go test ./...
go mod tidy
```

## Testing
From the project root, run:
```
go test ./...
```
or
```
go test ./... -cover
```
or
```
go test -v ./... -cover
```
depending on whether you want to see test coverage and how verbose the output you want.

## API Postman
https://api.postman.com/collections/9798442-00efe7b7-e03e-48c9-81ac-e3a61685bc0f?access_key=PMAT-01HQG4ER6PEVJMR5JN6AY2WGWK