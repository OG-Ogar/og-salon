package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"barber-booking/models"
	"barber-booking/storage"
)

func TestHomeFeedAuthorPhotoLinksToProfile(t *testing.T) {
	storage.Reset()
	t.Cleanup(storage.Reset)

	customer, err := storage.CreateCustomer(models.Customer{FullName: "Feed Author"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := storage.CreateSocialPost(context.Background(), models.SocialPost{
		AuthorRole:      models.SenderCustomer,
		AuthorID:        customer.ID,
		AuthorName:      customer.FullName,
		Text:            "A post from the feed author",
		BackgroundColor: "#fff0c2",
		TextStyle:       "normal",
	}); err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	CustomerDashboardHandler(response, httptest.NewRequest(http.MethodGet, "/customer/dashboard", nil), customer, sessionData{CSRFToken: "test"})
	if response.Code != http.StatusOK {
		t.Fatalf("status %d: %s", response.Code, response.Body.String())
	}

	body := response.Body.String()
	mustContain(t, "author photo link", body, `class="social-author-avatar"`)
	mustContain(t, "author photo link", body, `href="/customer/profile"`)
	mustContain(t, "author photo link", body, `class="social-author-name"`)
	mustContain(t, "author photo link", body, "A post from the feed author")
}

func TestHomeFeedShowsProviderRoleUnderAuthorName(t *testing.T) {
	storage.Reset()
	t.Cleanup(storage.Reset)

	customer, err := storage.CreateCustomer(models.Customer{FullName: "Feed Viewer"})
	if err != nil {
		t.Fatal(err)
	}
	barber, err := storage.CreateBarber(models.Barber{FullName: "Feed Barber", Email: "feed-barber@example.com", ProviderType: models.ProviderBarber})
	if err != nil {
		t.Fatal(err)
	}
	stylist, err := storage.CreateBarber(models.Barber{FullName: "Feed Stylist", Email: "feed-stylist@example.com", ProviderType: models.ProviderHairstylist})
	if err != nil {
		t.Fatal(err)
	}

	for _, author := range []models.Barber{barber, stylist} {
		if _, err := storage.CreateSocialPost(context.Background(), models.SocialPost{
			AuthorRole:      models.SenderBarber,
			AuthorID:        author.ID,
			AuthorName:      author.FullName,
			Text:            "A post from " + author.FullName,
			BackgroundColor: "#fff0c2",
			TextStyle:       "normal",
		}); err != nil {
			t.Fatal(err)
		}
	}

	response := httptest.NewRecorder()
	CustomerDashboardHandler(response, httptest.NewRequest(http.MethodGet, "/customer/dashboard", nil), customer, sessionData{CSRFToken: "test"})
	if response.Code != http.StatusOK {
		t.Fatalf("status %d: %s", response.Code, response.Body.String())
	}

	body := response.Body.String()
	mustContain(t, "barber role label", body, "Feed Barber</a></strong><span>Barber &middot;")
	mustContain(t, "stylist role label", body, "Feed Stylist</a></strong><span>Hair Stylist &middot;")
}

func TestSocialComposerOmitsSoundtrackPicker(t *testing.T) {
	storage.Reset()
	t.Cleanup(storage.Reset)

	customer, err := storage.CreateCustomer(models.Customer{FullName: "Sound Author"})
	if err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	SocialPostComposerHandler(response, httptest.NewRequest(http.MethodGet, "/customer/posts/new", nil), models.SenderCustomer, customer.FullName, customer.ProfileImage, sessionData{CSRFToken: "test"})
	if response.Code != http.StatusOK {
		t.Fatalf("status %d: %s", response.Code, response.Body.String())
	}

	body := response.Body.String()
	mustContain(t, "composer", body, `data-social-composer`)
	for _, unwanted := range []string{`data-sound-picker`, `name="music"`, `/static/music/`} {
		if strings.Contains(body, unwanted) {
			t.Fatalf("composer still exposes soundtrack UI %q", unwanted)
		}
	}
}

func TestHomeFeedOmitsSoundtrackPlayer(t *testing.T) {
	storage.Reset()
	t.Cleanup(storage.Reset)

	customer, err := storage.CreateCustomer(models.Customer{FullName: "Music Author"})
	if err != nil {
		t.Fatal(err)
	}
	// An old post may still carry a music ID in storage. The feed must not
	// render a broken player for it now that soundtracks are removed.
	if _, err := storage.CreateSocialPost(context.Background(), models.SocialPost{
		AuthorRole:      models.SenderCustomer,
		AuthorID:        customer.ID,
		AuthorName:      customer.FullName,
		BackgroundColor: "#fff0c2",
		TextStyle:       "normal",
		Image:           "/uploads/posts/photo.jpg",
		Music:           "smooth-criminal",
	}); err != nil {
		t.Fatal(err)
	}

	response := httptest.NewRecorder()
	CustomerDashboardHandler(response, httptest.NewRequest(http.MethodGet, "/customer/dashboard", nil), customer, sessionData{CSRFToken: "test"})
	if response.Code != http.StatusOK {
		t.Fatalf("status %d: %s", response.Code, response.Body.String())
	}

	body := response.Body.String()
	mustContain(t, "feed post", body, `id="post-`)
	for _, unwanted := range []string{`data-post-music`, `data-post-audio`, `/static/music/`} {
		if strings.Contains(body, unwanted) {
			t.Fatalf("feed still exposes soundtrack UI %q", unwanted)
		}
	}
}
