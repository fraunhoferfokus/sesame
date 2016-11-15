package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/authorization"
)

const (
	// pluginSocket denotes where the plugin is actually
	// loaded.
	pluginSocket = "/run/docker/plugins/sesame.sock"
	// Rules JSON is loaded from this path if not provided explicitly
	defaultRulesPath = "/etc/sesame/rules.json"
)

func main() {
	// Rules are either defined explicitly as first argument or are
	// to be found on the default path
	rulesPath := defaultRulesPath
	args := os.Args[1:]
	if len(args) > 0 {
		rulesPath = args[0]
	}

	sesame, err := newPlugin(rulesPath)
	if err != nil {
		log.Fatalf("Could not initiate plugin! (err: %s)", err)
	}

	h := authorization.NewHandler(sesame)
	err = h.ServeUnix("root", pluginSocket)
	if err != nil {
		log.Fatalf("Could not initiate handler! (err: %s)", err)
	}
}
