package nhkschool

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetAll(t *testing.T) {
	fixture, err := os.ReadFile("testdata/all.tsv")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nfsvideos/all/nhkforschool.tsv" {
			t.Errorf("path = %q", r.URL.Path)
		}
		if got := r.URL.Query().Get("apikey"); got != "test-key" {
			t.Errorf("apikey = %q", got)
		}
		w.Header().Set("Content-Type", "text/tab-separated-values")
		w.Write(fixture)
	}))
	defer srv.Close()

	client, _ := NewClient("test-key", WithBaseURL(srv.URL))
	videos, err := client.GetAll(context.Background())
	if err != nil {
		t.Fatalf("GetAll() error = %v", err)
	}

	if len(videos) != 2 {
		t.Fatalf("len = %d, want 2", len(videos))
	}

	v := videos[0]
	if v.ID != "D0005110412_00000" {
		t.Errorf("ID = %q", v.ID)
	}
	if v.Name != "とじこめられた空気" {
		t.Errorf("Name = %q", v.Name)
	}
	if len(v.Keywords) != 2 || v.Keywords[0] != "空気" {
		t.Errorf("Keywords = %v", v.Keywords)
	}
	if len(v.Grades) != 1 || v.Grades[0] != "小4" {
		t.Errorf("Grades = %v", v.Grades)
	}
}

func TestGetAll_APIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error"))
	}))
	defer srv.Close()

	client, _ := NewClient("test-key", WithBaseURL(srv.URL))
	_, err := client.GetAll(context.Background())
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}
