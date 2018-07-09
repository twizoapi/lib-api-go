CHANGELOG
=========

All notable changes to this project will be documented in this file.
This project adheres to [Semantic Versioning](http://semver.org/).

## 1.0.0-alpha - ????-??-??
### Compatibility break
- NewVerificationRequest will now also accept a string, and also return an error if any
### Added
- Function to retrieve account balance
- Added backup codes support
- Added time based one time (TOTP) support
- Added biovoice support
- Function to validate api key
- Refactored logger to expose DebugLogger, and remove all other log levels.
### Refactored
- Merged code into more logical files.  

## 0.1.0 - 2017-03-16
### Added
- Initial version
