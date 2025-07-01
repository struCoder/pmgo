# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Modern web interface for process management
- RESTful HTTP API with OpenAPI documentation
- Configuration file support (YAML, TOML, JSON)
- Process health checks and monitoring
- Structured logging with multiple formats
- Cross-platform binary builds
- Docker support
- Authentication and authorization
- Process metrics and monitoring
- Notification system (webhook, Slack, email)
- Process restart policies and limits
- Live reload for development
- Comprehensive test suite
- Performance optimizations

### Changed
- Upgraded to Go 1.21+ with modern Go patterns
- Replaced kingpin with cobra for CLI
- Modernized project structure
- Updated all dependencies to latest versions
- Improved error handling and logging
- Enhanced configuration management
- Better process isolation and security
- Optimized memory usage and performance

### Deprecated
- Legacy RPC-based communication (replaced with HTTP API)
- Old configuration format (still supported for backward compatibility)

### Removed
- CircleCI configuration (replaced with GitHub Actions)
- Outdated dependencies

### Fixed
- Memory leaks in process monitoring
- Race conditions in process management
- Improper signal handling
- Configuration file parsing issues

### Security
- Added authentication and authorization
- Improved process isolation
- Secure defaults for configuration

## [v0.6.1] - 2020-06-27

### Added
- Basic process management functionality
- RPC-based remote communication
- Simple CLI interface
- Process monitoring and auto-restart
- Configuration file support

### Fixed
- Various bug fixes and stability improvements

## [v0.6.0] - 2020-01-15

### Added
- Initial release of PMGO
- Process lifecycle management
- Basic monitoring capabilities
- CLI interface

[Unreleased]: https://github.com/struCoder/pmgo/compare/v0.6.1...HEAD
[v0.6.1]: https://github.com/struCoder/pmgo/compare/v0.6.0...v0.6.1
[v0.6.0]: https://github.com/struCoder/pmgo/releases/tag/v0.6.0
