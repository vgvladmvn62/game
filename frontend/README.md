# Bullseye frontend

## Overview

This is a frontend side for Bullseye made with React. It relies on [backend server](../backend).

## Prerequisites

- [Docker](https://www.docker.com)

For development using Chrome browser:
- [Redux dev tools for Chrome](http:github.com/zalmoxisus/redux-devtools-extension) with [src and instructions](http:github.com/zalmoxisus/redux-devtools-extension)
- [React dev tools for Chrome](https://chrome.google.com/webstore/detail/react-developer-tools/fmkadmapgofadopljbjfkapdkoienihi?hl=en) which works without additional setup

## Installation

```shell script
$ docker build -t {image name}:{image version} . -f Dockerfile
$ docker run -d -p 80:5000 {image name}:{image version} 
```

Now it will be running on [localhost](http://localhost). You can push a tested version to your Docker repository
and test both backend and frontend using [docker-compose](../backend/deployments/docker-compose.yml).
In order to push built image run:

```shell script
$ docker login
$ docker push {image name}:{image version}
```

## Configuration

Due to React reading environmental variables at build time rather than start time 
there was a necessity to store backend URL in file in order to run on Kubernetes.

You can set backend URL in [config file](public/config/config.js). Default value is `http://localhost:8080`.

## Usage

In order to use and test frontend visit the site where application was deployed and
follow instructions on the screen.