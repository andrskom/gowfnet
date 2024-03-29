# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

##[Unreleased]
### Changed
- Version of go to 1.20
- Linter to v1.55
- color 


##[1.0.1] - 2020-10-04
### Fixes
- License.

##[1.0.0] - 2020-10-04
### Added
- Config pkg with a minimal config based on an interface.
- Validator component amd error.
- Validators for a base flow.
- Dead places validator.
- Non-finish places validator.
- Duplicated places in places validator.
- Duplicated places in transitions validator.
- Error in state pkg with tests.
- State in state pkg with tests.
- New implementation of net and tests for it.
- Listeners for state and net.
- Coverage for github action.
- E2e test.
### Changed
- Linter version to 1.24
- GO version to 1.14
- Use gomock instead 
### Removed
- Old implementations.

## [0.4.0] -2020-03-04
### Fixed
- New linter errors.
### Changed
- Use ErrStack instead Error in state.
- Error's workflow of state.

## [0.3.4] - 2019-09-17
### Fixed
- If listener does not set, AfterPlaced will not be run.

## [0.3.3] - 2019-09-16
### Fixed
- JSON serialization of err.

## [0.3.2] - 2019-09-16
### Fixed
- Bug in build net function.

## [0.3.1] - 2019-09-16
### Fixed
- Err in call of listener in net.

## [0.3.0] - 2019-09-16
### Added
- Net in listener.

## [0.2.1] - 2019-09-16
### Added
- Add state's getters.

## [0.2.0] - 2019-09-16
### Added
- Global listener for all nets in registry.

## [0.1.1] - 2019-09-15
### Added
- Method set listener.

## [0.1.0] - 2019-09-15
### Added
- Implementation of network.
- Automatic transition listener.
- Part of unit test.
- Registry of nets.
