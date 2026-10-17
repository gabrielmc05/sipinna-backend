package main

import "testing"

func TestReportRoutesAreRegistered(t *testing.T) {
	want := map[string]bool{
		"GET /api/v1/reportes/historico":         false,
		"POST /api/v1/reportes/asignar/:folio":   false,
		"GET /api/v1/reportes/:folio/evidencias": false,
		"GET /api/v1/territorios":                false,
		"GET /api/v1/territorios/:id":            false,
		"GET /api/v1/usuario/alimentadores":      false,
		"GET /api/v1/mapa/reportes":              false,
		"GET /api/v1/instituciones":              false,
		"GET /api/v1/instituciones/:id":          false,
		"POST /api/v1/instituciones":             false,
		"PATCH /api/v1/instituciones/:id":        false,
		"PATCH /api/v1/instituciones/:id/estado": false,
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
