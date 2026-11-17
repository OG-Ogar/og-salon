package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"barber-booking/models"
	"barber-booking/storage"
)

func TestExploreProfileMessageLabelFollowsProviderType(t *testing.T) {
	storage.Reset()
	t.Cleanup(storage.Reset)

	cases := []struct {
		name          string
		provider      string
		wantLabel     string
		wantBookLabel string
		wantAbsent    string
	}{
		{"barber", models.ProviderBarber, "Message barber", "Book this barber", "Message stylist"},
		{"stylist", models.ProviderHairstylist, "Message stylist", "Book this stylist", "Message barber"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			provider, err := storage.CreateBarber(models.Barber{
				Email:        tc.name + "_profile@example.com",
				FullName:     tc.name + " profile",
				ShopName:     tc.name + " shop",
				ProviderType: tc.provider,
				Status:       models.BarberActive,
			})
			if err != nil {
				t.Fatal(err)
			}
			id := strconv.Itoa(provider.ID)
			recorder := httptest.NewRecorder()
			ExploreBarberHandler(recorder, httptest.NewRequest(http.MethodGet, "/explore/"+id, nil), id)
			if recorder.Code != http.StatusOK {
				t.Fatalf("status %d: %s", recorder.Code, recorder.Body.String())
			}
			body := recorder.Body.String()
			if !strings.Contains(body, tc.wantLabel) {
				t.Fatalf("missing %q", tc.wantLabel)
			}
			if !strings.Contains(body, tc.wantBookLabel) {
				t.Fatalf("missing %q", tc.wantBookLabel)
			}
			if strings.Contains(body, tc.wantAbsent) {
				t.Fatalf("unexpected %q", tc.wantAbsent)
			}
		})
	}
}
