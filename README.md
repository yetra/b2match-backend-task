# b2match-backend-task

## Description

An API for event management. 

## Installation

### Requirements

* [Go](https://go.dev/) 1.20
* [SQLite](https://www.sqlite.org/index.html)

Install the requirements from `go.mod` in the project's root directory:
```shell
go mod tidy
```

## Usage

Start the server on `localhost:8085` using the following commad:
```shell
go run main.go
```

Go to the `/swagger/index.html` page to see the API documentation.

## Tests

Tests are located in the `tests` package.

To run them from the command line, `cd` into the `tests` directory and execute the following command:
```shell
go test
```

For manual testing, an exported [Insomnia](https://insomnia.rest/) JSON is provided in the `extras` directory.
