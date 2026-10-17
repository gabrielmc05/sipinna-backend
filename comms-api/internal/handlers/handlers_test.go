// handlers_test.go es un archivo de prueba para el paquete handlers.
// Contiene pruebas unitarias para las funciones validateReportRequest, newFolio y detectEvidenceType.

package handlers

import (
	"strings"
	"testing"

	"github.com/gabrielmc05/sipinna-backend/comms-api/internal/dboper"
)

func TestValidateReportRequest(t *testing.T) {
	latitude := 0.0
	longitude := 0.0
	request := reportRequest{
		Description:     "Se reporta una situación",
		ImmediateDanger: "No sé",
		Location: locationInput{
			Latitude:    &latitude,
			Longitude:   &longitude,
			TerritoryID: 1,
		},
		Minors: []dboper.Minor{{Gender: "Masculino"}},
	}
	if err := validateReportRequest(request); err != nil {
		t.Fatalf("valid report was rejected: %v", err)
	}

	request.Location.Latitude = nil
	if err := validateReportRequest(request); err == nil {
		t.Fatal("report without latitude was accepted")
	}
}

func TestValidateImmediateDanger(t *testing.T) {
	for _, danger := range []string{"Sí", "No", "No sé"} {
		if !isImmediateDanger(danger) {
			t.Errorf("valid immediate danger %q was rejected", danger)
		}
	}
	for _, danger := range []string{"", "Desconocido", "si"} {
		if isImmediateDanger(danger) {
			t.Errorf("invalid immediate danger %q was accepted", danger)
		}
	}

	latitude, longitude := 0.0, 0.0
	request := reportRequest{
		Description:     "Reporte de prueba",
		ImmediateDanger: "No sé",
		Location:        locationInput{Latitude: &latitude, Longitude: &longitude, TerritoryID: 1},
		Minors:          []dboper.Minor{{Gender: "No identificado"}},
	}
	if err := validateReportRequest(request); err != nil {
		t.Fatalf("valid danger value was rejected: %v", err)
	}
	request.ImmediateDanger = ""
	if err := validateReportRequest(request); err == nil {
		t.Fatal("missing immediate danger value was accepted")
	}
}

func TestNewFolioFormat(t *testing.T) {
	folio, err := newFolio()
	if err != nil {
		t.Fatalf("newFolio() error = %v", err)
	}
	if len(folio) != 15 || !strings.HasPrefix(folio, "REP-") {
		t.Fatalf("newFolio() = %q, want REP-<year>-<six characters>", folio)
	}
}

func TestReportStatusValidation(t *testing.T) {
	for _, status := range []string{"Registrado", "En revisión", "En seguimiento", "Canalizado", "Concluido", "Archivado", "Cancelado", "Reincidente"} {
		if !isHistoryStatus(status) {
			t.Errorf("history status %q was rejected", status)
		}
	}
	for _, status := range []string{"", "inventado"} {
		if isHistoryStatus(status) {
			t.Errorf("invalid history status %q was accepted", status)
		}
	}

	if !isReportStatus("Canalizado a Procuraduría") {
		t.Error("existing report status validation changed")
	}
}

func TestDetectEvidenceType(t *testing.T) {
	contentType, err := detectEvidenceType(strings.NewReader("\x89PNG\r\n\x1a\n" + strings.Repeat("\x00", 16)))
	if err != nil {
		t.Fatalf("detectEvidenceType() error = %v", err)
	}
	if contentType != "image/png" {
		t.Fatalf("detectEvidenceType() = %q, want image/png", contentType)
	}

	for _, fixture := range []struct {
		name string
		data string
		want string
	}{
		{name: "webp", data: "RIFFxxxxWEBP" + strings.Repeat("\x00", 8), want: "image/webp"},
		{name: "webm", data: "\x1a\x45\xdf\xa3" + strings.Repeat("\x00", 8), want: "video/webm"},
		{name: "mp4", data: "xxxxftypisom" + strings.Repeat("\x00", 8), want: "video/mp4"},
	} {
		t.Run(fixture.name, func(t *testing.T) {
			got, err := detectEvidenceType(strings.NewReader(fixture.data))
			if err != nil {
				t.Fatalf("detectEvidenceType() error = %v", err)
			}
			if got != fixture.want {
				t.Fatalf("detectEvidenceType() = %q, want %q", got, fixture.want)
			}
		})
	}

	if _, err := detectEvidenceType(strings.NewReader("not an image")); err == nil {
		t.Fatal("unsupported evidence type was accepted")
	}
}
