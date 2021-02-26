# agent

The Seashell Agent

------------------

Agent meant for executing on edge devices and synchronizing their configurations with those defined in the Seashell Platform.

## Requirements
- Golang 1.16+

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
`TODO`

## Contributing
- Fork it
- Download your fork (git clone https://github.com/your_username/agent && cd agent)
- Create your feature branch (git checkout -b my-new-feature)
- Make changes and stage them (git add .)
- Commit your changes (git commit -m 'Add some feature')
- Push to the branch (git push origin my-new-feature)
- Create new pull request


## Roadmap
- [ ] Code coverage

## License
The Seashell Agent is released under the Apache 2.0 license. See LICENSE.txt
