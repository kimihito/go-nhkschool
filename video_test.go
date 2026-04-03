package nhkschool

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetVideo(t *testing.T) {
	fixture, err := os.ReadFile("testdata/video.json")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nfsvideo/id/D0005110412_00000" {
			t.Errorf("path = %q, want /nfsvideo/id/D0005110412_00000", r.URL.Path)
		}
		if got := r.URL.Query().Get("apikey"); got != "test-key" {
			t.Errorf("apikey = %q, want %q", got, "test-key")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer srv.Close()

	client, _ := NewClient("test-key", WithBaseURL(srv.URL))
	video, err := client.GetVideo(context.Background(), "D0005110412_00000")
	if err != nil {
		t.Fatalf("GetVideo() error = %v", err)
	}

	if video.ID != "D0005110412_00000" {
		t.Errorf("ID = %q, want %q", video.ID, "D0005110412_00000")
	}
	if video.Name != "とじこめられた空気" {
		t.Errorf("Name = %q, want %q", video.Name, "とじこめられた空気")
	}
	if video.About.NFSSeriesName != "ふしぎエンドレス　理科４年" {
		t.Errorf("NFSSeriesName = %q", video.About.NFSSeriesName)
	}
	if video.Text != nil {
		t.Errorf("Text = %v, want nil", video.Text)
	}
	if len(video.Parts) != 1 {
		t.Fatalf("Parts len = %d, want 1", len(video.Parts))
	}
	if video.Parts[0].Name != "オープニング" {
		t.Errorf("Parts[0].Name = %q", video.Parts[0].Name)
	}
}

func TestGetVideo_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
	}))
	defer srv.Close()

	client, _ := NewClient("test-key", WithBaseURL(srv.URL))
	_, err := client.GetVideo(context.Background(), "INVALID")
	if err == nil {
		t.Fatal("expected error for 404 response")
	}

	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 404 {
		t.Errorf("StatusCode = %d, want 404", apiErr.StatusCode)
	}
}
