# OpenRFSense Node
The node software is designed to run as a daemon on remote/embedded devices, which serve as RF monitoring stations. This daemon is responsible for:
- Handling commands sent by the backend through MQTT and responding accordingly
- Showing a simple web interface exposed on `localhost`, easily accessible via the temporary hotspot (see [`openrfsense/image`](https://github.com/openrfsense/image))

## Table of contents <!-- omit in toc -->
- [OpenRFSense Node](#openrfsense-node)
    - [Configuration](#configuration)
      - [YAML](#yaml)
      - [Environment variables](#environment-variables)
    - [MQTT](#mqtt)
      - [System statistics/metrics](#system-statisticsmetrics)

### Configuration
> ⚠️ The configuration is WIP

The `config` module from [`openrfsense/common`](https://github.com/openrfsense/common) is used. As such, configuration values are loaded from a YAML file first, then from environment variables.

#### YAML
See the example [`config.yml`](./config.yml) file for now, as the configuration is very prone to change.

#### Environment variables
Environment variables are defined as follows: `ORFS_SECTION_SUBSECTION_KEY=value`. They are loaded after any other configuration file, so they cam be used to overwrite any configuration value.

### MQTT
Nodes use the MQTT 3.1.1 protocol as a lightweight way to receive requests from the backend and respond accordingly. The communication flow is simple: a request is sent by the backend to an appropriate request channel and a response is sent by the nodes on another channel. Channels are structured in the following way (following [Steve's suggestions for MQTT topic structure](http://www.steves-internet-guide.com/mqtt-topic-payload-design-notes/)):

- Top-level (broadcast) channel with path `CHANNEL/SUBCHANNEL/`:
  - Requests (from the backend) go to `node/METHOD/CHANNEL/SUBCHANNEL/`
  - Responses (from the node) go to `node/CHANNEL/SUBCHANNEL/`
- Device-specific channel with path `CHANNEL/SUBCHANNEL/`:
  - Requests (from the backend) go to `node/METHOD/NODE_ID/CHANNEL/SUBCHANNEL/`
  - Responses (from the node) go to `node/NODE_ID/CHANNEL/SUBCHANNEL/`

Where `METHOD` can be one of (following HTTP standard methods or more to be defined):
- `post`
- `get`

#### System statistics/metrics
> This section will probably get moved, but it felt right to include it in this readme

