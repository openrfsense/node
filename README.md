# OpenRFSense Node
The node software is designed to run as a daemon on remote/embedded devices, which serve as RF monitoring stations. This daemon is responsible for:
- Handling commands sent by the backend through MQTT and responding accordingly
- Showing a simple web interface exposed on `localhost`, easily accessible via the temporary hotspot (see [`openrfsense/image`](https://github.com/openrfsense/image))

## Table of contents <!-- omit in toc -->
- [OpenRFSense Node](#openrfsense-node)
    - [Configuration](#configuration)
      - [YAML](#yaml)
      - [Environment variables](#environment-variables)
    - [NATS](#nats)
      - [System statistics/metrics](#system-statisticsmetrics)

### Configuration
> ⚠️ The configuration is WIP

The `config` module from [`openrfsense/common`](https://github.com/openrfsense/common) is used. As such, configuration values are loaded from a YAML file first, then from environment variables.

#### YAML
See the example [`config.yml`](./config.yml) file for now, as the configuration is very prone to change.

#### Environment variables
Environment variables are defined as follows: `ORFS_SECTION_SUBSECTION_KEY=value`. They are loaded after any other configuration file, so they cam be used to overwrite any configuration value.

### NATS
Nodes use [NATS](https://nats.io/) to exchange messages, under the `node.` root subject. Messages are encoded with JSON and relayed in NATS' own wire format. The subject structure can be generalized as follows:
- General, network-wide or broadcast messages:
  - Backend sends a request on `node.$channel`
  - Nodes send NATS replies on a single arbitrary/unique inbox
- Single-node requests:
  - Backend sends a request on `node.$id.$channel`
  - The node sends a NATS reply on an arbitrary/unique inbox

#### System statistics/metrics
> This section will probably get moved, but it felt right to include it in this readme

