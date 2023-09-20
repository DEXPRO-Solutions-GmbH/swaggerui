# Changelog

## Next

### Added

- Add `/swaggerui` as alternative to `swagger-ui`
- Add `WithReplaceServerUrls` option as alternative to `WithAddServerUrls`

### Fixed

- 

## 1.0.1

### Fixed

- Add Documentation of central types to improve package usability.

## 1.0.0

### Added

- Initial project commit
- Handler for serving OpenAPI spec via SwaggerUI
- Support for middlewares to modify Spec on ever request
- Middleware to replace OIDC url in a security component
- Middleware to set server urls based on the incoming request
