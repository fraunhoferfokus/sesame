package main

import (
  "fmt"
  "io/ioutil"
	"encoding/json"
)

const (
  rulesPath = "/etc/sesame/rules.json"
)

// Rule defines an allowed action by defining method and path pattern of a
// Docker Remote API.
// see https://docs.docker.com/engine/reference/api/docker_remote_api/
type Rule struct {
  Method string   `json:"method"`
  Pattern string  `json:"pattern"`
}

// sesame implements the Plugin inteface of Docker authorization API and manages
// authorization usign RuleSets
type sesame struct {
  rules map[string][]Rule
}

// newPlugin creates a new Sesame plugin.
// It first loads the rules and then registers the unix socket.
func newPlugin() (*sesame, error) {
  var plugin sesame

  // Read the rules and decode them
  content, err := ioutil.ReadFile(rulesPath)
  if err != nil {
    return nil, err
  }

  err = json.Unmarshal(content, &plugin.rules)
  if err != nil {
    return nil, err
  }

  // We're gut to go!
  return &plugin, nil
}
