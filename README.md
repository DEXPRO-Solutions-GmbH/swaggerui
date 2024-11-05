# SwaggerUI Handler For Gin

[![Go Reference](https://pkg.go.dev/badge/github.com/DEXPRO-Solutions-GmbH/swaggerui.svg)](https://pkg.go.dev/github.com/DEXPRO-Solutions-GmbH/swaggerui)
[![Go Report Card](https://goreportcard.com/badge/github.com/DEXPRO-Solutions-GmbH/swaggerui)](https://goreportcard.com/report/github.com/DEXPRO-Solutions-GmbH/swaggerui)

This project implements a gin Handler which allows you to
expose your OpenAPI spec via SwaggerUI.

## Installing

```shell
go get github.com/DEXPRO-Solutions-GmbH/swaggerui
```

## Maintenance

### SwaggerUI upgrade

Every then and now we want to upgrade the builtin SwaggerUI.

Do do that, go to https://github.com/swagger-api/swagger-ui and download the latest release.

Unzip the `/dist` directory from that release and replace all contents of the local `/dist` directory.

Then, apply the git patch `patch-swagger-ui-initializer.patch`.
