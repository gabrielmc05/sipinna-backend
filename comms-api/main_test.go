package main

import "testing"

func TestReportRoutesAreRegistered(t *testing.T) {
	want := map[string]bool{
		"GET /api/v1/reportes/historico":         false,
		"POST /api/v1/reportes/asignar/:folio":   false,
		"GET /api/v1/reportes/:folio/evidencias": false,
	}
	for _, route := range newRouter().Routes() {
		key := route.Method + " " + route.Path
		if _, exists := want[key]; exists {
			want[key] = true
		}
	}
	for route, registered := range want {
		if !registered {
			t.Errorf("route %q was not registered", route)
		}
	}
}
