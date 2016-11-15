package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/docker/go-plugins-helpers/authorization"
)

const (
	rulesPath = "/etc/sesame/rules.json"
)

// Rule defines an allowed action by defining method and path pattern of a
// Docker Remote API.
// see https://docs.docker.com/engine/reference/api/docker_remote_api/
type Rule struct {
	Method  string `json:"method"`
	Pattern string `json:"pattern"`
}

// matches checks if given method equals to Rule's method and if the given URI
// matches the Rule's pattern
func (r *Rule) matches(method string, uri string) bool {
	matched := false

	if strings.EqualFold(r.Method, method) {
		// TODO take care of errors!
		matched, _ = regexp.MatchString(r.Pattern, uri)
	}

	return matched
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

func (p *sesame) AuthZReq(req authorization.Request) authorization.Response {
	user := req.User
	method := req.RequestMethod
	uri := req.RequestURI

	if rules, ok := p.rules[user]; ok {

		for _, rule := range rules {
			if rule.matches(method, uri) {
				return authorization.Response{Allow: true}
			}
		}
	} else {
		return authorization.Response{
			Allow: false,
			Msg:   fmt.Sprintf("User '%s' not found!", user),
		}
	}

	return authorization.Response{
		Allow: false,
		Msg:   fmt.Sprintf("User '%s' forbidden to %s on %s!", user, method, uri),
	}
}

func (p *sesame) AuthZRes(req authorization.Request) authorization.Response {
	// Our decision is final!
	return authorization.Response{Allow: true}
}
