# Bullseye Raspberry Pi client

## Overview

This part contains implementation of Raspberry Pi client which communicates through
MQTT with Bullseye.

## Prerequisites

- [Golang](https://golang.org/) (latest version)

This client should run on Raspberry Pi with an internet access to retrieve requests from MQTT broker.

## Installation

Build application using provided dev tools .

```shell script
$ make build
```

It creates executable with name defined inside the [Makefile](./Makefile) (by default - **bullseye-rpi-develop**).
Run this application on Raspberry Pi later. 

## Configuration

It might be useful to run client on system startup for administrator convenience.

## Usage

There are no additional steps required after installation. Client is self-sufficient - after retrieving
requests from MQTT it performs actions on hardware.