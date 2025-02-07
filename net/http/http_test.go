package http

import "testing"

func TestNewHTTPRequestTags(t *testing.T) {
	rt := NewHTTPRequestTags("a", "c", "f", "r", "app", "d")

	if rt.action != "a" {
		t.Errorf("rt.action - got: %s, want: %s", rt.action, "a")
	}

	if rt.controller != "c" {
		t.Errorf("rt.controller - got: %s, want: %s", rt.controller, "c")
	}

	if rt.framework != "f" {
		t.Errorf("rt.framework - got: %s, want: %s", rt.framework, "f")
	}

	if rt.route != "r" {
		t.Errorf("rt.route - got: %s, want: %s", rt.route, "r")
	}

	if rt.application != "app" {
		t.Errorf("rt.application - got: %s, want: %s", rt.application, "app")
	}

	if rt.driver != "d" {
		t.Errorf("rt.driver - got: %s, want: %s", rt.driver, "d")
	}
}
