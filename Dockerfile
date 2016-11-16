FROM alpine:3.4
MAINTAINER Yan Foto <yan.foto@quaintous.com>

LABEL Description="A simple and flexible authorization plugin for Docker"
LABEL Version=0.1.0

COPY build/sesame /usr/bin/sesame

# Volume to load rules from and another to create socket in
VOLUME ["/etc/sesame", "/run/docker/plugins"]

ENTRYPOINT ["sesame"]
