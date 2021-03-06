package main

import (
	"fmt"
	"testing"

	auth "github.com/docker/go-plugins-helpers/authorization"
)

func mustDo(what bool, t *testing.T, res auth.Response) {
	if res.Allow != what {
		t.Errorf("Should allow: '%t',\n Response: %#v", what, res)
	}
}

func mustAllow(t *testing.T, res auth.Response) {
	mustDo(true, t, res)
}

func mustDeny(t *testing.T, res auth.Response) {
	mustDo(false, t, res)
}

func TestInvalidFiles(t *testing.T) {
	// Non-existing files
	t.Run("Rules files does not exist", func(t *testing.T) {
		_, err := newPlugin("./testdata/doesntexist.json")
		if err == nil {
			t.Error("Should NOT accept non-existing files!")
		}
	})

	// Invalid JSON rules
	t.Run("Rules files is not a valid JSON", func(t *testing.T) {
		_, err := newPlugin("./testdata/rules.invalid.json")
		if err == nil {
			t.Error("Should NOT accept invalid JSON files!")
		}
	})
}

func TestAuthZReq(t *testing.T) {
	p, err := newPlugin("./testdata/rules.json")
	if err != nil {
		t.Error(err)
	}

	// Should allow listing containers (e.g. ps) for 'readonly' user but nothing more!
	t.Run("Only read Containers", func(t *testing.T) {
		req := auth.Request{User: "readonly", RequestMethod: "GET"}

		req.RequestURI = "/v1.24/containers/json"
		res := p.AuthZReq(req)
		mustAllow(t, res)

		req.RequestURI = "/v1.24/version"
		res = p.AuthZReq(req)
		mustDeny(t, res)

		req.RequestMethod = "POST"
		req.RequestURI = "/v1.24/containers/MINE/start"
		res = p.AuthZReq(req)
		mustDeny(t, res)
	})

	// Should allow building images for user 'buildonly' but nothing more!
	t.Run("Only build images", func(t *testing.T) {
		req := auth.Request{User: "buildonly", RequestMethod: "POST"}

		req.RequestURI = "/v1.24/build"
		res := p.AuthZReq(req)
		mustAllow(t, res)

		req.RequestMethod = "GET"
		req.RequestURI = "/images/MINE/history"
		res = p.AuthZReq(req)
		mustDeny(t, res)
	})

	// Should allow multiple operations for user 'onecontainer' for container 'MYCONTAINER'
	t.Run("Multiple container operations", func(t *testing.T) {
		req := auth.Request{User: "onecontainer", RequestMethod: "POST"}
		for _, op := range []string{"start", "stop", "restart", "kill", "update", "pause", "unpause"} {
			req.RequestURI = fmt.Sprintf("/v1.24/containers/MYCONTAINER/%s", op)
			res := p.AuthZReq(req)
			mustAllow(t, res)

			req.RequestURI = fmt.Sprintf("/v1.24/containers/OTHERCONTAINER/%s", op)
			res = p.AuthZReq(req)
			mustDeny(t, res)
		}
	})

	t.Run("Multiple method operations", func(t *testing.T) {
		req := auth.Request{User: "allcontainers", RequestMethod: "GET"}

		req.RequestURI = "/v1.24/containers/json"
		res := p.AuthZReq(req)
		mustAllow(t, res)

		req.RequestMethod = "POST"
		req.RequestURI = "/v1.24/containers/create"
		res = p.AuthZReq(req)
		mustAllow(t, res)
	})

	t.Run("Non-existing user", func(t *testing.T) {
		req := auth.Request{
			User:          "bob",
			RequestMethod: "GET",
			RequestURI:    "/v1.24/containers/MINE/start",
		}

		res := p.AuthZReq(req)
		mustDeny(t, res)
	})
}

func TestAuthZRes(t *testing.T) {
	p, err := newPlugin("./testdata/rules.json")
	if err != nil {
		t.Error(err)
	}

	// Should always allow!
	req := auth.Request{User: "WHOEVER", RequestMethod: "INVALID", RequestURI: "WHEREVER"}
	res := p.AuthZRes(req)
	mustAllow(t, res)
}
