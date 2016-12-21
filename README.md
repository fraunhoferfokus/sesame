# Sesame
[![Build Status](https://travis-ci.org/fraunhoferfokus/sesame.svg?branch=master)](https://travis-ci.org/fraunhoferfokus/sesame)
[![Go Report Card](https://goreportcard.com/badge/github.com/fraunhoferfokus/sesame)](https://goreportcard.com/report/github.com/fraunhoferfokus/sesame)
[![License](https://img.shields.io/github/license/fraunhoferfokus/sesame.svg)](https://github.com/fraunhoferfokus/sesame/blob/master/LICENSE)

Sesame is a Docker [access authorization plugin](http://docs-stage.docker.com/engine/extend/plugins_authorization/) with
focous on simplicity and flexibility that is easy to configure and adapt.

## Quick start
Assuming that Sesame is to run on the same host as Docker engine, the easiest way is to start it with your rules (i.e. `rules.json`)
and let docker engine use it as an authorization plugin. But first fetch the latest release of `sesame`
[here](https://github.com/fraunhoferfokus/sesame/releases/latest) and then:

```bash
# Start plugin (i.e. downloaded binary) with rules.json
./sesame rules.json
# Let Docker run with TLS enabled and use sesame as authorization plugin
# see https://docs.docker.com/engine/admin/#/configuring-docker for more
dockerd \
  --tlsverify \
  --tlscacert=~/.docker/ca.pem \
  --tlscert=~/.docker/cert.pem \
  --tlskey=~/.docker/key.pem \
  --authorization-plugin=sesame
```

If no parameters (e.g. `rules.json`) is given, sesame tries to load the rules from `/etc/sesame/rules.json`. For more on rules, see
[**Rules**](#rules) section.

***NOTE:*** Sesame can also run as a Docker container. See [**Docker Image**](#docker-image) section for more.

## Why this plugin?
Docker engine changes rapidly and so does its [remote API](https://docs.docker.com/engine/reference/api/docker_remote_api/). To keep up
with this changes, this plugin focuses on simplicity and flexibility so in case of introducing any breaking changes in the remote API,
you would only need to update your rules and enhance them with matching rules to the new specification. You wouldn't need to wait for
the plugin to adapt to the new changes.

## Rules
Docker clients communicate with Docker engine over [remote API](https://docs.docker.com/engine/reference/api/docker_remote_api/) which
is (mostly) an HTTP REST API. A rule is a combination of an HTTP method template (e.g. `GET|PUT`) and a URI template (e.g. `/containers/.*`)
to define ressources of remote API. Rules use regular expressions to define method and URI templates, allowing them to match multiple API
resources.

To define a rule, for example, matching starting, pausing, and unpausing of a container called `MYCONTAINER` can be done using a single
rule:

```json
{
    "method": "POST",
    "pattern": "/containers/MYCONTAINER/(start|pause|unpause)$"
}
```

The `rules.json` file associates users with an array of rules:

```json
{
    "someUser": [{
       "method": "GET",
       "pattern": "/version$"
    }],
    "otherUser": [
    {
        "method": "GET",
        "pattern": "/containers/json$"
    },
    {
        "method": "POST",
        "pattern": "/build"
    }
    ]
}
```

The example above allows `someUser` to get Docker engine's version and `otherUser` to list containers and build images. To see more examples
see `testdata/rules.json` and take a look at `plugin_test.go`.

## Docker Image
Sesame can also be used inside a container. Currently there are no ready-made Sesame images available (but planned) on Docker hub. However,
an image can directly be built from the source:

```bash
# Clone the repository
https://github.com/fraunhoferfokus/sesame.git && cd sesame
# Build the binaries
make build
# Build docker image
docker build -t fraunhofer/sesame .
```

**Note** that this method requires having `git`, `make`, `go >= 1.7`, and obviously `docker` installed. You can also download the latest
binary (as mentioned in quick start) and build an image without being required to clone the repository and build from the source using the
following `Dockerfile`:

```Dockerfile
FROM alpine:3.4
MAINTAINER Yan Foto <yan.foto@quaintous.com>

LABEL Description="A simple and flexible authorization plugin for Docker"

COPY sesame /usr/bin/sesame

# Volume to load rules from and another to create socket in
VOLUME ["/etc/sesame", "/run/docker/plugins"]

ENTRYPOINT ["sesame"]
```

To start the container, rules should be mounted at `/etc/sesame/rules.json` and the socket is provided at
`/run/docker/plugins/sesame.sock` of the running container. A typical procedure would look as follows:

```bash
# Start sesame as a daemon inside a container
docker run -d --restart=always --name sesame \
    -v `pwd`/rules.json:/etc/sesame/rules.json \
    -v /run/docker/plugins:/run/docker/plugins \
    fraunhofer/sesame
# Let Docker run with TLS enabled and use sesame as authorization plugin
# see https://docs.docker.com/engine/admin/#/configuring-docker for more
dockerd \
  --tlsverify \
  --tlscacert=~/.docker/ca.pem \
  --tlscert=~/.docker/cert.pem \
  --tlskey=~/.docker/key.pem \
  --authorization-plugin=sesame
```

This example assumes that `rules.json` is available under working directory.