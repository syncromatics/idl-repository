# Interface Definition Language Repository (idl-repository)

A repository and command line interface (CLI) for storing and fetching interface definitions in various formats. ([Protocol Buffers][protobuf], [Avro][avro], [OpenAPI][openapi], etc.)

[protobuf]: https://developers.google.com/protocol-buffers/
[avro]: https://avro.apache.org/
[openapi]: https://swagger.io/docs/specification/about/

## Quickstart

### Initialize a project for the first time

Initializing a project creates a new `idl.yaml` file with the project name and repository URL.

```bash
idl init my-project http://idl-repository.example.com
```

Read more about [`idl init`][idl-init].

### Pull dependencies into your project

Pull IDLs into your project.

```bash
idl pull
```

Read more about [`idl pull`][idl-pull].

### Push project to the repository

Push IDLs in your project to the configured repository.

```bash
idl push
```

Read more about [`idl push`][idl-push].

### Documentation

Read the full documentation for [`idl`][idl].

Read the full documentation for the repository service, [`idl-repository`][idl-repository].

[idl]: docs/idl/idl.md
[idl-init]: docs/idl/idl_init.md
[idl-pull]: docs/idl/idl_pull.md
[idl-push]: docs/idl/idl_push.md
[idl-repository]: docs/idl-repository/idl-repository.md

## Building

[![Travis](https://img.shields.io/travis/syncromatics/idl-repository.svg)](https://travis-ci.org/syncromatics/idl-repository)
[![Docker Build Status](https://img.shields.io/docker/build/syncromatics/idl-repository.svg)](https://hub.docker.com/r/syncromatics/idl-repository/)

Ensure you have your `GOPATH` configured properly. (Typically, you'll want to check this repo out to `$(go env GOPATH)/src/github.com/syncromatics/idl-repository`.) You'll also need Docker to build the Docker image for the repository.

Download dependencies

```bash
go get -t -v ./...
```

Build the CLI and repository server

```bash
make build
```

## Code of Conduct

We are committed to fostering an open and welcoming environment. Please read our [code of conduct](CODE_OF_CONDUCT.md) before participating in or contributing to this project.

## Contributing

We welcome contributions and collaboration on this project. Please read our [contributor's guide](CONTRIBUTING.md) to understand how best to work with us.

## License and Authors

[![GMV Syncromatics Engineering logo](https://secure.gravatar.com/avatar/645145afc5c0bc24ba24c3d86228ad39?size=16) GMV Syncromatics Engineering](https://github.com/syncromatics)

[![license](https://img.shields.io/github/license/syncromatics/idl-repository.svg)](https://github.com/syncromatics/idl-repository/blob/master/LICENSE)
[![GitHub contributors](https://img.shields.io/github/contributors/syncromatics/idl-repository.svg)](https://github.com/syncromatics/idl-repository/graphs/contributors)

This software is made available by GMV Syncromatics Engineering under the MIT license.