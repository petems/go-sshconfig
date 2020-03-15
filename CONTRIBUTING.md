## Contributing

1. Fork the repo.

1. Create a separate branch for your change.

1. Run the tests. We only take pull requests with passing tests, and
   documentation.

1. Add a test for your change. Only refactoring and documentation
   changes require no new tests. If you are adding functionality
   or fixing a bug, please add a test.

1. Squash your commits down into logical components. Make sure to rebase
   against the current master.

1. Push the branch to your fork and submit a pull request.

Please be prepared to repeat some of these steps as our contributors review
your code.

## Syntax and style

Run `lint` to detect style issues and perform fixes:

    make lint

Or

    golangci-lint run ./... 

## Running the unit tests

The unit test suite covers most of the code, as mentioned above please
add tests if you're adding new functionality.

To run your all the unit tests:

    make lint

Or without lint:

    go test ./...