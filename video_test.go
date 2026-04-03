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

func TestListByCSCode(t *testing.T) {
	fixture, err := os.ReadFile("testdata/list.json")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nfsvideos/cscode/8260243111100000" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("apikey"); got != "test-key" {
			t.Errorf("apikey = %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer srv.Close()

	client, _ := NewClient("test-key", WithBaseURL(srv.URL))
	resp, err := client.ListByCSCode(context.Background(), "8260243111100000", nil)
	if err != nil {
		t.Fatalf("ListByCSCode() error = %v", err)
	}

	if resp.TotalCount != 2 {
		t.Errorf("TotalCount = %d, want 2", resp.TotalCount)
	}
	if len(resp.Videos) != 2 {
		t.Fatalf("Videos len = %d, want 2", len(resp.Videos))
	}
	if resp.Videos[0].ID != "D0005110412_00000" {
		t.Errorf("Videos[0].ID = %q", resp.Videos[0].ID)
	}
}

func TestListByCSCode_WithOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if got := q.Get("includelower"); got != "1" {
			t.Errorf("includelower = %q, want %q", got, "1")
		}
		if got := q.Get("contenttype"); got != "1" {
			t.Errorf("contenttype = %q, want %q", got, "1")
		}
		if got := q.Get("perpage"); got != "5" {
			t.Errorf("perpage = %q, want %q", got, "5")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"totalCount":0,"page":1,"perPage":5,"result":[]}`))
	}))
	defer srv.Close()

	client, _ := NewClient("test-key", WithBaseURL(srv.URL))
	includeLower := 1
	contentType := 1
	perPage := 5
	_, err := client.ListByCSCode(context.Background(), "8260243111100000", &ListOptions{
		IncludeLower: &includeLower,
		ContentType:  &contentType,
		PerPage:      &perPage,
	})
	if err != nil {
		t.Fatalf("ListByCSCode() error = %v", err)
	}
}

// Tests for Task 5 guard fixes (added during code review)
func TestSearchByKeyword_NilParams(t *testing.T) {
	client, _ := NewClient("test-key")
	_, err := client.SearchByKeyword(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error for nil params")
	}
}

func TestSearchByKeyword_EmptyKeywords(t *testing.T) {
	client, _ := NewClient("test-key")
	_, err := client.SearchByKeyword(context.Background(), &KeywordParams{Keywords: ""})
	if err == nil {
		t.Fatal("expected error for empty keywords")
	}
}
