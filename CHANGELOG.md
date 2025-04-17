# Changelog

## [Unreleased]

### Added
- Linux-style tab completion for files/directories
  - Trigger with space+Tab
  - Shows multiple matches
  - Adds / suffix for directories
- New subcommands implementation:
  - cat: Display file contents
  - ls: List directory contents
  - curl: Make HTTP requests
  - wget: Download files
  - clear: Clear terminal
- New interactive mode features:
  - Command history navigation
  - Better error handling
  - Context-aware input processing

### Changed
- Refactored input handling system into modular components
- Improved terminal interaction and responsiveness
- Restructured command processing pipeline

### Fixed
- Fixed unresponsive input issues in interactive mode
- Improved Ctrl+C handling and interrupt recovery
- Fixed various edge cases in command parsing
