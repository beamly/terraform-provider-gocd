# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [0.6.5] - 11-12-2017
### Changed
 - Fixed bugs around loading config profiles
 - Discontinued ARM releases
 - Widespread refactoring of interal CLI structures to make testing more viable

## [0.3.10] - 15-08-2017
### Changed
 - Disabled fpm build
 - Expanding test coverage
 - Removed version as a specific arg to SDK for pipeline configs.

## [0.3.8] - 14-08-2017
### Changed
 - Fixed a bug where fpm was not building deb's on release

## [0.3.7] - 12-08-2017
### Changed
 - Fixed a bug where version was not being set on build
 - Reduced the amount of duplicate code.

## [0.3.5] - 09-07-2017
### Added
 - Expanded test cases for jobs and tasks
 - Added Pipeline status querying and pipeline history structure
 - Added Basic get-configuration command.
 - Added plugin retrieval

### Changed
 - Fixed a bug where if the config file was not found the cli tool would crash.

## [0.3.3] - 07-07-2017
### Added
 - Added missing structure definitions for pipeline history

## [0.3.2] - 05-07-2017
### Changed
 - Refactored tests so that duplicate code is handled by common functions
 - All testcases and integration tests passing.

## [0.3.0] - 04-07-2017
### Changed
 - Did a lot of code cleanup, mainly around the build process, linting, and made sure godoc build works.

## [0.2.9] - 04-07-2017
### Changed
 - Fixed linting issues


## [0.2.8] - 04-08-2017
### Changed
 - Fixed bad paths in build frameworks, which caused builds on travis to fail.
 - Fixed a bug where builds were not being correctly generated.

## [0.2.3] - 30-07-2017
### Added
 - Added `pipeline config` creation in sdk and cli tool.
 - Added `pipeline config` update in sdk and cli tool.
 - Changed cli package name.

### Changed
 - Moved gocd package into subdirectory.
 - Added `doc.go` for godoc.com.

## [0.2.1] - 30-07-2017
### Changed
 - Added `goreleaser` to travis build
 - Made sure deploy occurs on golang 1.8

## [0.1.12] - 30-07-2017

### Changed
 - Fixed Travis release process

## [0.1.0] - 30-07-2017
### Added
 - Initial Release
 - CLI tool
   - Create, List, Get, Update build agents
   - Create, List, Get config templates