# Bullseye backend

## Overview

This part contains implementation of Bullseye server which handles all requests from
user, manages configuration, communicates with Enterprise Commerce instance and send requests
to Bullseye Raspberry PI client using MQTT.

## Prerequisites

- [Docker](https://www.docker.com)
- [helm](https://helm.sh)
- [Kyma](https://kyma-project.io/)

Docker is used for development ( provided support for [docker-compose](https://docs.docker.com/compose/)) and
managing images. Deployment on Kyma cluster can be done using prepared [helm chart](./deployments/chart/bullseye).
 
## Installation

Refer to [installation](./docs/installation.org) document.

## Configuration

Refer to [configuration](./docs/configuration.org) document.

## Usage

After running server it should be accessible on `localhost:8080` by default (configuration
can be changed). Usage of server is strictly connected with providing an appropriate
configuration which has been explained in sections above.