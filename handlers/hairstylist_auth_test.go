package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHairstylistLoginPage(t *testing.T) {
	if err := LoadTemplates("../templates"); err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	HairstylistLoginHandler(response, httptest.NewRequest(http.MethodGet, "/hairstylist/login", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.Code)
	}

	body := response.Body.String()
	for _, want := range []string{
		"Hair Stylist Login",
		`action="/hairstylist/login"`,
		`/barber/forgot-password`,
		"/barber/register",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("page missing %q", want)
		}
	}
}

func TestHairstylistRouterRedirectsToLoginAndSharedAccountPages(t *testing.T) {
	for _, tc := range []struct {
		path     string
		location string
	}{
		{"/hairstylist", "/hairstylist/login"},
		{"/hairstylist/register", "/barber/register?provider=hairstylist"},
		{"/hairstylist/forgot-password", "/barber/forgot-password"},
		{"/hairstylist/reset-password?token=abc", "/barber/reset-password?token=abc"},
	} {
		response := httptest.NewRecorder()
		HairstylistRouter(response, httptest.NewRequest(http.MethodGet, tc.path, nil))

		if response.Code != http.StatusSeeOther {
			t.Errorf("%s status = %d, want 303", tc.path, response.Code)
		}
		if got := response.Header().Get("Location"); got != tc.location {
			t.Errorf("%s location = %q, want %q", tc.path, got, tc.location)
		}
	}
}
