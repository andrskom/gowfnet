# gowfnet 

![Test](https://github.com/andrskom/gowfnet/workflows/Test/badge.svg)
![Lint](https://github.com/andrskom/gowfnet/workflows/Lint/badge.svg)

Golang implementation of Workflow networks.

## Development

Use `make` for run lint and tests.

## Status

*ALPHA*

_Major version use as second number in semver._

_Wait for version 1.0.0 for stable using._

## State

### Listener overwriting model

You can set the listener for the state in few ways according to next priority list from low to high.

1. Set for the state.
1. Set for the net.

The state store only one listener in time.
If you set a listener for the state when you try to make an operation with a net than had a listener
only the net's listener will be called.  