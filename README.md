<h1 align="center">
    Seashell
<br></h1>

<h5 align="center">
The Seashell CLI
</h5>

------------------

<p align="center">
  <a href="https://goreportcard.com/report/github.com/seashell/agent"><img src="https://goreportcard.com/badge/github.com/seashell/agent" alt="Go report: A+"></a>
  <img alt="GitHub" src="https://img.shields.io/github/license/seashell/seashell">
</p>

TODO: description

## Requirements
- Golang 1.16+
- Systemd
- Nomad 1.0.0+
- Consul 1.9.0+

## Build

System requirements:
- Golang 1.14+
- Node 10.17.0+
- yarn 1.12.3+

```
$ go generate
$ go build
```

Alternatively, you can build with `make`, for example:

```bash
$ make dev
```

To see help on building with make, run:

```bash
$ make help
```

Build for all architectures with

```bash
make release
```

## Usage

```bash
seashell agent --config=<config_file>
```

An example configuration can be found in `/dist/seashell.hcl`

## Overview

## Supported modules

## Configuration

- `log_level` :

- `name` : Device name for identifying it in Nomad and Consul

- `data_dir` :

- `advertise_addr` : 

## API

The Seashell agent exposes a simple REST API that allows for simple system information queries.

- `GET /status` :reports modules' `systemd` services current active state and substate.

Sample response:

```bash
$ curl -X GET localhost:5345/status
    Nomad: active running
    Consul: failed failed
```

#### Coming soon  :clock1:
- Etcd as a storage backend
- RPC API for clients nodes to interact with the server
- Fine-grained authorization
- CLI improvements


## Contributing
- Fork it
- Download your fork (git clone https://github.com/your_username/seashell && cd seashell)
- Create your feature branch (git checkout -b my-new-feature)
- Make changes and stage them (git add .)
- Commit your changes (git commit -m 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create new pull request


## Roadmap
- [ ] Website
- [ ] Code coverage

## License
Seashell is released under the Apache 2.0 license. See LICENSE.txt
