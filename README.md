# Rapid [![Build Status][ci-img]][ci] [![GoDoc][godoc-img]][godoc]

Rapid is a Go library for property-based testing.

Rapid checks that properties you define hold for a large number
of automatically generated test cases. If a failure is found, rapid
automatically minimizes the failing test case before presenting it.

Property-based testing emphasizes thinking about high level properties
the program should satisfy rather than coming up with a list
of individual examples of desired behavior (test cases).
This results in concise and powerful tests that are a pleasure to write.

Design and implementation of rapid are heavily inspired by
[Hypothesis](https://github.com/HypothesisWorks/hypothesis), which is itself
a descendant of [QuickCheck](https://hackage.haskell.org/package/QuickCheck).

## Features

- Idiomatic Go API
  - Designed to be used together with `go test` and the `testing` package
  - Works great with
    [testify/require](https://godoc.org/github.com/stretchr/testify/require) and
    [testify/assert](https://godoc.org/github.com/stretchr/testify/assert)
- Automatic minimization of failing test cases
- No dependencies outside of the Go standard library

### Planned features

- Automatic persistence of failing test cases

## Examples

Example [function](./example_function_test.go) and
[state machine](./example_statemachine_test.go) tests are provided.
They both fail. Making them pass is a good way to get first real experience
of working with rapid.

## Comparison

Rapid aims to bring to Go the power and convenience Hypothesis brings to Python.

Compared to [gopter](https://godoc.org/github.com/leanovate/gopter), rapid:

- has a much simpler API
- does not require any user code to minimize failing test cases
- uses more sophisticated algorithms for data generation

Compared to [testing/quick](https://golang.org/pkg/testing/quick/), rapid:

- provides much more control over test case generation
- supports state machine ("stateful" or "model-based") testing
- automatically minimizes any failing test case
- has to settle for `rapid.Check` being the main exported function
  instead of much more stylish `quick.Check`
 
## Status

Rapid is alpha software. Important pieces of functionality are missing;
API breakage and bugs should be expected.

If rapid fails to find a bug you believe it should, or the failing test case
that rapid reports does not look like a minimal one,
please [open an issue](https://github.com/flyingmutant/rapid/issues).

## License

Rapid is licensed under the [Mozilla Public License version 2.0](./LICENSE). 

[ci-img]: https://travis-ci.org/flyingmutant/rapid.svg?branch=master
[ci]: https://travis-ci.org/flyingmutant/rapid
[godoc-img]: https://godoc.org/github.com/flyingmutant/rapid?status.svg
[godoc]: https://godoc.org/github.com/flyingmutant/rapid
