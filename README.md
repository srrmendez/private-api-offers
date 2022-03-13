# private api offers

Private api offers server

## Content table

- [Swagger API Documentation](#swagger-api-documentation)
- [Build](#build)

## Swagger API Documentation

- To generate swagger documentation first run command

```bash
swag init --parseDependency -g app.go 
```

## Build

- To build run command:

```bash
GOOS=<so> GOARCH=<arch> go build
```