# nippo-cli
The tool to power my nippo.

## Install `nippo` command
```
go install github.com/c18t/nippo-cli/nippo@latest
```

## Usage
### Setup
```shell
nippo init
```

### Build
```shell
nippo build
```

### Publish
```shell
nippo deploy
```

## Setting up your development environment
```console
// host
$ (echo UID=$(id -u) & echo GID=$(id -g)) > .env
$ docker compose up -d
$ docker compose exec nippo-cli bash

// container
$ go run nippo/nippo.go
$ make
```

## License
[MIT](./LICENSE)

## Author
ɯ̹t͡ɕʲi
- [github / c18t](https://github.com/c18t)
