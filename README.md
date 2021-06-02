# gowfnet

**It was big mistake to put some logic into state.
That's why this project will be archived.
See a new implementation in gopetri.** 

![Test](https://github.com/andrskom/gowfnet/workflows/Test/badge.svg)
![Lint](https://github.com/andrskom/gowfnet/workflows/Lint/badge.svg)

The Golang implementation of Workflow networks.

## Development

Use `make` for run lint and tests.

## Status

The only minimal implementation is finished.

It's tested and works well. The Docs for this is in progress.

You can see examples in e2e tests.

### Listener overwriting model

You can set the listener for the state in few ways according to next priority list from low to high.

1. Set for the state.
1. Set for the net.


The state store only one listener in time.
If you set a listener for the state when you try to make an operation with a net than had a listener
only the net's listener will be called.  