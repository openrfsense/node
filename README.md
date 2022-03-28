# OpenRFSense Node
The node software is designed to run as a daemon on remote/embedded devices, which serve as RF monitoring stations. This daemon is responsible for:
- Handling commands sent by the backend through MQTT and responding accordingly
- Showing a simple web interface exposed on localhost, easily accessible via the temporary hotspot (see [openrfsense/image](https://github.com/openrfsense/image))

## Table of contents <!-- omit in toc -->
- [OpenRFSense Node](#openrfsense-node)
    - [Configuration](#configuration)
      - [YAML](#yaml)
      - [Environment variables](#environment-variables)

### Configuration
> ⚠️ The configuration is WIP

The `config` module from [openrfsense/common](https://github.com/openrfsense/common) is used. As such, configuration values can be loaded from several YAML files, and finally environment variables.

#### YAML

#### Environment variables
Environment variables are defined as follows: `ORFS_SECTION_SUBSECTION_KEY=value`. They are loaded after any other configuration file, so they cam be used to overwrite any configuration value.