package main

import (
	"github.com/docker/go-plugins-helpers/authorization"
	log "github.com/Sirupsen/logrus"
)

const (
	// pluginSocket denotes where the plugin is actually
	// loaded.
	pluginSocket = "/run/docker/plugins/sesame.sock"
)


func main() {
	sesame, err := newPlugin()
	if err != nil {
		log.Fatalf("Could not initiate plugin! (err: %s)", err)
	}

	h := authorization.NewHandler(sesame)
	err = h.ServeUnix("root", pluginSocket)
	if err != nil {
		log.Fatalf("Could not initiate handler! (err: %s)", err)
	}
	log.Infof("Listening on %s", pluginSocket)
}
