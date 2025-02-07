// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

type HTTPRequestTags struct {
	action      string
	controller  string
	framework   string
	route       string
	application string
	driver      string
}

func NewHTTPRequestTags(action, controller, framework, route, application, driver string) *HTTPRequestTags {
	return &HTTPRequestTags{action, controller, framework, route, application, driver}
}

func (h *HTTPRequestTags) Action() string {
	return h.action
}

func (h *HTTPRequestTags) Controller() string {
	return h.controller
}

func (h *HTTPRequestTags) Framework() string {
	return h.framework
}

func (h *HTTPRequestTags) Route() string {
	return h.route
}

func (h *HTTPRequestTags) Application() string {
	return h.application
}

func (h *HTTPRequestTags) Driver() string {
	return h.driver
}
