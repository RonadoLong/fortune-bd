# Home-Service Service

This is the Home-Service service

Generated with

```
micro new shop-micro/service/home-service --namespace=shop --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: shop.srv.home-service
- Type: srv
- Alias: home-service

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
./home-service-srv
```

Build a docker image
```
make docker
```