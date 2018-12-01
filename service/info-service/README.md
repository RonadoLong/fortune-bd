# Video Service

This is the Video service

Generated with

```
micro new shop-micro/video --namespace=shop --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: shop.srv.video
- Type: srv
- Alias: video

## Dependencies

Micro services depend on service discovery. The default is consul.

```
# install consul
brew install consul

# run consul
consul agent -dev
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./video-srv
```

Build a docker image
```
make docker
```