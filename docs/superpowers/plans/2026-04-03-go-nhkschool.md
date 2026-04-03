# go-nhkschool Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a Go client library + CLI for NHK for School API v2, covering all 4 endpoints (content list, content detail, keyword search, bulk data).

**Architecture:** Flat package `nhkschool` with Functional Options client, TDD with `httptest`, CLI via `flag` + subcommands. Standard library only — zero external dependencies.

**Tech Stack:** Go 1.24, standard library (`net/http`, `net/http/httptest`, `encoding/json`, `encoding/csv`, `flag`)

---

## File Structure

| File | Responsibility |
|------|---------------|
| `nhkschool.go` | `Client` struct, `NewClient`, `Option` type, `WithHTTPClient`, `WithBaseURL`, internal `do` helper for HTTP requests |
| `types.go` | All response types: `Video`, `About`, `CurriculumStandard`, `Part`, `VideoSummary`, `ListResponse`, `ListOptions`, `KeywordParams` |
| `error.go` | `APIError` struct and its `Error()` method |
| `video.go` | `GetVideo`, `ListByCSCode`, `SearchByKeyword` methods |
| `all.go` | `GetAll` method (TSV parsing) |
| `nhkschool_test.go` | Tests for `NewClient` |
| `video_test.go` | Tests for `GetVideo`, `ListByCSCode`, `SearchByKeyword` |
| `all_test.go` | Tests for `GetAll` |
| `testdata/video.json` | JSON fixture for `GetVideo` |
| `testdata/list.json` | JSON fixture for `ListByCSCode` |
| `testdata/keyword.json` | JSON fixture for `SearchByKeyword` |
| `testdata/all.tsv` | TSV fixture for `GetAll` |
| `cmd/nhkschool/main.go` | CLI entry point with subcommands |
| `.gitignore` | Build artifacts, IDE files, OS files, design docs |

---

### Task 1: Project scaffold — replace existing code with new module

**Files:**
- Delete: `nfs.go`, `nfs_test.go`, `go.sum`
- Modify: `go.mod`
- Verify: `.gitignore` (already created)

- [ ] **Step 1: Delete old source files**

```bash
cd /Users/kimihito/ghq/github.com/kimihito/nhk-for-school-go
rm nfs.go nfs_test.go go.sum
```

- [ ] **Step 2: Rewrite go.mod**

Replace the contents of `go.mod` with:

```
module github.com/kimihito/go-nhkschool

go 1.24
```

- [ ] **Step 3: Create directory structure**

```bash
mkdir -p testdata cmd/nhkschool
```

- [ ] **Step 4: Verify clean state**

```bash
go mod tidy
```

Expected: no errors, no dependencies in `go.sum`.

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "chore: scaffold go-nhkschool module, remove old code"
```

---

### Task 2: Error type

**Files:**
- Create: `error.go`
- Test: `error_test.go` (lightweight, test Error() output)

- [ ] **Step 1: Write the failing test**

Create `error_test.go`:

```go
package nhkschool

import (
	"errors"
	"testing"
)

func TestAPIError_Error(t *testing.T) {
	err := &APIError{StatusCode: 404, Body: "Not Found"}
	want := "nhkschool: HTTP 404: Not Found"
	if got := err.Error(); got != want {
		t.Errorf("Error() = %q, want %q", got, want)
	}
}

func TestAPIError_ErrorsAs(t *testing.T) {
	orig := &APIError{StatusCode: 429, Body: "Rate limited"}
	var wrapped error = orig

	var apiErr *APIError
	if !errors.As(wrapped, &apiErr) {
		t.Fatal("errors.As failed to match *APIError")
	}
	if apiErr.StatusCode != 429 {
		t.Errorf("StatusCode = %d, want 429", apiErr.StatusCode)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test -run TestAPIError -v
```

Expected: FAIL — `APIError` not defined.

- [ ] **Step 3: Write implementation**

Create `error.go`:

```go
package nhkschool

import "fmt"

// APIError is returned when the API responds with a non-2xx status code.
type APIError struct {
	StatusCode int
	Body       string
}

func (e *APIError) Error() string {
	return fmt.Sprintf("nhkschool: HTTP %d: %s", e.StatusCode, e.Body)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test -run TestAPIError -v
```

Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add error.go error_test.go
git commit -m "feat: add APIError type"
```

---

### Task 3: Response types

**Files:**
- Create: `types.go`

- [ ] **Step 1: Create types.go**

```go
package nhkschool

// Video represents the full detail of a video from the content detail API.
type Video struct {
	ID             string               `json:"id"`
	ContentType    string               `json:"contentType"`
	Name           string               `json:"name"`
	About          About                `json:"about"`
	Description    string               `json:"description"`
	Text           *string              `json:"text"`
	URL            string               `json:"url"`
	ThumbnailURL   string               `json:"thumbnailUrl"`
	Grades         []string             `json:"grades"`
	SubjectAreas   []string             `json:"subjectAreas"`
	Subjects       []string             `json:"subjects"`
	Curriculum     []CurriculumStandard `json:"curriculumStandard"`
	Keywords       []string             `json:"keywords"`
	Duration       string               `json:"duration"`
	UploadDate     string               `json:"uploadDate"`
	DateModified   string               `json:"dateModified"`
	DatePublished  string               `json:"datePublished"`
	Expires        *string              `json:"expires"`
	RegionsAllowed *string              `json:"regionsAllowed"`
	UsageInfo      string               `json:"usageInfo"`
	Bitrate        string               `json:"bitrate"`
	Height         int                  `json:"height"`
	Width          int                  `json:"width"`
	Parts          []Part               `json:"hasPart"`
}

// About contains the series information for a video.
type About struct {
	NFSSeriesName string `json:"nfsSeriesName"`
}

// CurriculumStandard represents a curriculum standard code mapping.
type CurriculumStandard struct {
	Version string `json:"curriculumStandardVersion"`
	NfsID   string `json:"curriculumStandardNfsId"`
	Code    string `json:"curriculumStandardCode"`
}

// Part represents a chapter/segment within a video.
type Part struct {
	ClipNumber   int     `json:"clipNumber"`
	StartOffset  float64 `json:"startOffset"`
	EndOffset    float64 `json:"endOffset"`
	ThumbnailURL string  `json:"thumbnailUrl"`
	Name         string  `json:"name"`
}

// VideoSummary represents a row from the bulk data TSV API.
type VideoSummary struct {
	ID             string
	ContentType    string
	Name           string
	NFSSeriesName  string
	Description    string
	URL            string
	ThumbnailURL   string
	Grades         []string
	SubjectAreas   []string
	Subjects       []string
	CurriculumCode []string
	Keywords       []string
	Duration       string
	DatePublished  string
	Expires        string
	RegionsAllowed string
}

// ListResponse is the paginated response from list and keyword search APIs.
type ListResponse struct {
	Videos     []*Video
	TotalCount int
	Page       int
	PerPage    int
}

// ListOptions are optional parameters for ListByCSCode.
type ListOptions struct {
	IncludeLower *int
	ResultOrder  *int
	ContentType  *int
	Page         *int
	PerPage      *int
}

// KeywordParams are the parameters for SearchByKeyword.
type KeywordParams struct {
	Keywords     string
	Grades       []string
	SubjectAreas []string
	Subjects     []string
	ResultOrder  *int
	ContentType  *int
	Page         *int
	PerPage      *int
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 3: Commit**

```bash
git add types.go
git commit -m "feat: add response and parameter types"
```

---

### Task 4: Client with NewClient and options

**Files:**
- Create: `nhkschool.go`
- Test: `nhkschool_test.go`

- [ ] **Step 1: Write the failing tests**

Create `nhkschool_test.go`:

```go
package nhkschool

import (
	"net/http"
	"testing"
)

func TestNewClient(t *testing.T) {
	c, err := NewClient("my-key")
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if c.apiKey != "my-key" {
		t.Errorf("apiKey = %q, want %q", c.apiKey, "my-key")
	}
	if c.baseURL != "https://api.nhk.or.jp/school/v2" {
		t.Errorf("baseURL = %q, want default", c.baseURL)
	}
	if c.httpClient != http.DefaultClient {
		t.Error("httpClient should default to http.DefaultClient")
	}
}

func TestNewClient_EmptyAPIKey(t *testing.T) {
	_, err := NewClient("")
	if err == nil {
		t.Fatal("NewClient('') should return an error")
	}
}

func TestNewClient_WithBaseURL(t *testing.T) {
	c, err := NewClient("key", WithBaseURL("http://localhost:9999"))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if c.baseURL != "http://localhost:9999" {
		t.Errorf("baseURL = %q, want %q", c.baseURL, "http://localhost:9999")
	}
}

func TestNewClient_WithHTTPClient(t *testing.T) {
	custom := &http.Client{}
	c, err := NewClient("key", WithHTTPClient(custom))
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if c.httpClient != custom {
		t.Error("httpClient should be the custom client")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test -run TestNewClient -v
```

Expected: FAIL — `NewClient` not defined.

- [ ] **Step 3: Write implementation**

Create `nhkschool.go`:

```go
package nhkschool

import (
	"context"
	"errors"
	"io"
	"net/http"
)

const defaultBaseURL = "https://api.nhk.or.jp/school/v2"

// Client is an NHK for School API client.
type Client struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) {
		cl.httpClient = c
	}
}

// WithBaseURL overrides the default API base URL.
func WithBaseURL(url string) Option {
	return func(cl *Client) {
		cl.baseURL = url
	}
}

// NewClient creates a new NHK for School API client.
// apiKey is required and must not be empty.
func NewClient(apiKey string, opts ...Option) (*Client, error) {
	if apiKey == "" {
		return nil, errors.New("nhkschool: apiKey must not be empty")
	}
	c := &Client{
		apiKey:     apiKey,
		baseURL:    defaultBaseURL,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c, nil
}

// do executes an HTTP GET request and returns the response body.
// It returns an *APIError for non-2xx status codes.
func (c *Client) do(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("nhkschool: creating request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nhkschool: executing request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("nhkschool: reading response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &APIError{
			StatusCode: resp.StatusCode,
			Body:       string(body),
		}
	}

	return body, nil
}
```

- [ ] **Step 4: Fix the missing import**

The `do` method uses `fmt`. Update the import block in `nhkschool.go`:

```go
import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)
```

- [ ] **Step 5: Run tests to verify they pass**

```bash
go test -run TestNewClient -v
```

Expected: PASS (all 4 tests).

- [ ] **Step 6: Commit**

```bash
git add nhkschool.go nhkschool_test.go
git commit -m "feat: add Client with NewClient and functional options"
```

---

### Task 5: GetVideo

**Files:**
- Create: `testdata/video.json`
- Create: `video.go`
- Create: `video_test.go`

- [ ] **Step 1: Create test fixture**

Create `testdata/video.json`:

```json
{
  "queryData": {
    "uri": "https://api.nhk.or.jp/school/v2/nfsvideo/id/D0005110412_00000"
  },
  "result": [
    {
      "id": "D0005110412_00000",
      "contentType": "ばんぐみ",
      "name": "とじこめられた空気",
      "about": {
        "nfsSeriesName": "ふしぎエンドレス　理科４年"
      },
      "description": "空気でっぽうの説明",
      "text": null,
      "url": "https://www2.nhk.or.jp/school/watch/bangumi/?das_id=D0005110412_00000",
      "thumbnailUrl": "https://www.nhk.or.jp/das/image/D0005110/D0005110412_00000_S_001.jpg",
      "grades": ["小4"],
      "subjectAreas": ["理科"],
      "subjects": [],
      "curriculumStandard": [
        {
          "curriculumStandardVersion": "8",
          "curriculumStandardNfsId": "小学 理科 4年 A 1 ア ア",
          "curriculumStandardCode": "8260243111100000"
        }
      ],
      "keywords": ["空気", "体積"],
      "duration": "PT0H10M0S",
      "uploadDate": "2018-11-06T04:00:00+09:00",
      "dateModified": "2022-08-17T12:51:02+09:00",
      "datePublished": "2018-11-06T09:30:00+09:00",
      "expires": null,
      "regionsAllowed": null,
      "usageInfo": "streaming",
      "bitrate": "512kbps",
      "height": 360,
      "width": 640,
      "hasPart": [
        {
          "clipNumber": 1,
          "startOffset": 0,
          "endOffset": 24.991,
          "thumbnailUrl": "https://www.nhk.or.jp/das/image/D0005110/D0005110412_00000_C_001.jpg",
          "name": "オープニング"
        }
      ]
    }
  ]
}
```

- [ ] **Step 2: Write the failing test**

Create `video_test.go`:

```go
package nhkschool

import (
	"context"
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
```

Add `"errors"` to the import block of `video_test.go`.

- [ ] **Step 3: Run tests to verify they fail**

```bash
go test -run TestGetVideo -v
```

Expected: FAIL — `GetVideo` not defined.

- [ ] **Step 4: Write implementation**

Create `video.go`:

```go
package nhkschool

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

// videoResponse is the raw JSON envelope for video endpoints.
type videoResponse struct {
	Result []*Video `json:"result"`
}

// listRawResponse is the raw JSON envelope for list endpoints.
type listRawResponse struct {
	TotalCount int      `json:"totalCount"`
	Page       int      `json:"page"`
	PerPage    int      `json:"perPage"`
	Result     []*Video `json:"result"`
}

// GetVideo retrieves a single video by its ID.
func (c *Client) GetVideo(ctx context.Context, id string) (*Video, error) {
	u := fmt.Sprintf("%s/nfsvideo/id/%s?apikey=%s", c.baseURL, url.PathEscape(id), url.QueryEscape(c.apiKey))

	body, err := c.do(ctx, u)
	if err != nil {
		return nil, err
	}

	var resp videoResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("nhkschool: decoding response: %w", err)
	}

	if len(resp.Result) == 0 {
		return nil, fmt.Errorf("nhkschool: no video found for id %q", id)
	}

	return resp.Result[0], nil
}

// ListByCSCode retrieves videos matching a curriculum standard code.
func (c *Client) ListByCSCode(ctx context.Context, cscode string, opts *ListOptions) (*ListResponse, error) {
	q := url.Values{}
	q.Set("apikey", c.apiKey)

	if opts != nil {
		if opts.IncludeLower != nil {
			q.Set("includelower", strconv.Itoa(*opts.IncludeLower))
		}
		if opts.ResultOrder != nil {
			q.Set("resultorder", strconv.Itoa(*opts.ResultOrder))
		}
		if opts.ContentType != nil {
			q.Set("contenttype", strconv.Itoa(*opts.ContentType))
		}
		if opts.Page != nil {
			q.Set("page", strconv.Itoa(*opts.Page))
		}
		if opts.PerPage != nil {
			q.Set("perpage", strconv.Itoa(*opts.PerPage))
		}
	}

	u := fmt.Sprintf("%s/nfsvideos/cscode/%s?%s", c.baseURL, url.PathEscape(cscode), q.Encode())

	body, err := c.do(ctx, u)
	if err != nil {
		return nil, err
	}

	var raw listRawResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("nhkschool: decoding response: %w", err)
	}

	return &ListResponse{
		Videos:     raw.Result,
		TotalCount: raw.TotalCount,
		Page:       raw.Page,
		PerPage:    raw.PerPage,
	}, nil
}

// SearchByKeyword searches videos by keyword and optional filters.
func (c *Client) SearchByKeyword(ctx context.Context, params *KeywordParams) (*ListResponse, error) {
	if params.SubjectAreas != nil && params.Subjects != nil {
		return nil, fmt.Errorf("nhkschool: subjectareas and subjects cannot be specified together")
	}

	q := url.Values{}
	q.Set("apikey", c.apiKey)
	q.Set("keywords", params.Keywords)

	if params.Grades != nil {
		q.Set("grades", strings.Join(params.Grades, " OR "))
	}
	if params.SubjectAreas != nil {
		q.Set("subjectareas", strings.Join(params.SubjectAreas, " OR "))
	}
	if params.Subjects != nil {
		q.Set("subjects", strings.Join(params.Subjects, " OR "))
	}
	if params.ResultOrder != nil {
		q.Set("resultorder", strconv.Itoa(*params.ResultOrder))
	}
	if params.ContentType != nil {
		q.Set("contenttype", strconv.Itoa(*params.ContentType))
	}
	if params.Page != nil {
		q.Set("page", strconv.Itoa(*params.Page))
	}
	if params.PerPage != nil {
		q.Set("perpage", strconv.Itoa(*params.PerPage))
	}

	u := fmt.Sprintf("%s/nfsvideos/keyword?%s", c.baseURL, q.Encode())

	body, err := c.do(ctx, u)
	if err != nil {
		return nil, err
	}

	var raw listRawResponse
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("nhkschool: decoding response: %w", err)
	}

	return &ListResponse{
		Videos:     raw.Result,
		TotalCount: raw.TotalCount,
		Page:       raw.Page,
		PerPage:    raw.PerPage,
	}, nil
}
```

- [ ] **Step 5: Run tests to verify they pass**

```bash
go test -run TestGetVideo -v
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add video.go video_test.go testdata/video.json
git commit -m "feat: add GetVideo with tests"
```

---

### Task 6: ListByCSCode

**Files:**
- Create: `testdata/list.json`
- Modify: `video_test.go`

- [ ] **Step 1: Create test fixture**

Create `testdata/list.json`:

```json
{
  "totalCount": 2,
  "page": 1,
  "perPage": 20,
  "result": [
    {
      "id": "D0005110412_00000",
      "contentType": "ばんぐみ",
      "name": "とじこめられた空気",
      "about": {"nfsSeriesName": "ふしぎエンドレス　理科４年"},
      "description": "空気でっぽうの説明",
      "text": null,
      "url": "https://www2.nhk.or.jp/school/watch/bangumi/?das_id=D0005110412_00000",
      "thumbnailUrl": "https://www.nhk.or.jp/das/image/D0005110/D0005110412_00000_S_001.jpg",
      "grades": ["小4"],
      "subjectAreas": ["理科"],
      "subjects": [],
      "curriculumStandard": [],
      "keywords": ["空気"],
      "duration": "PT0H10M0S",
      "uploadDate": "2018-11-06T04:00:00+09:00",
      "dateModified": "2022-08-17T12:51:02+09:00",
      "datePublished": "2018-11-06T09:30:00+09:00",
      "expires": null,
      "regionsAllowed": null,
      "usageInfo": "streaming",
      "bitrate": "512kbps",
      "height": 360,
      "width": 640,
      "hasPart": []
    },
    {
      "id": "D0005110413_00000",
      "contentType": "ばんぐみ",
      "name": "とじこめられた水",
      "about": {"nfsSeriesName": "ふしぎエンドレス　理科４年"},
      "description": "水でっぽうの説明",
      "text": null,
      "url": "https://www2.nhk.or.jp/school/watch/bangumi/?das_id=D0005110413_00000",
      "thumbnailUrl": "https://www.nhk.or.jp/das/image/D0005110/D0005110413_00000_S_001.jpg",
      "grades": ["小4"],
      "subjectAreas": ["理科"],
      "subjects": [],
      "curriculumStandard": [],
      "keywords": ["水"],
      "duration": "PT0H10M0S",
      "uploadDate": "2018-11-06T04:00:00+09:00",
      "dateModified": "2022-08-17T12:51:02+09:00",
      "datePublished": "2018-11-06T09:30:00+09:00",
      "expires": null,
      "regionsAllowed": null,
      "usageInfo": "streaming",
      "bitrate": "512kbps",
      "height": 360,
      "width": 640,
      "hasPart": []
    }
  ]
}
```

- [ ] **Step 2: Write the failing tests**

Add to `video_test.go`:

```go
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
```

- [ ] **Step 3: Run tests to verify they pass**

The implementation is already in `video.go` from Task 5.

```bash
go test -run TestListByCSCode -v
```

Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add video_test.go testdata/list.json
git commit -m "test: add ListByCSCode tests"
```

---

### Task 7: SearchByKeyword

**Files:**
- Create: `testdata/keyword.json`
- Modify: `video_test.go`

- [ ] **Step 1: Create test fixture**

Create `testdata/keyword.json`:

```json
{
  "totalCount": 1,
  "page": 1,
  "perPage": 20,
  "result": [
    {
      "id": "D0005110412_00000",
      "contentType": "ばんぐみ",
      "name": "とじこめられた空気",
      "about": {"nfsSeriesName": "ふしぎエンドレス　理科４年"},
      "description": "空気でっぽうの説明",
      "text": null,
      "url": "https://www2.nhk.or.jp/school/watch/bangumi/?das_id=D0005110412_00000",
      "thumbnailUrl": "https://www.nhk.or.jp/das/image/D0005110/D0005110412_00000_S_001.jpg",
      "grades": ["小4"],
      "subjectAreas": ["理科"],
      "subjects": [],
      "curriculumStandard": [],
      "keywords": ["空気", "体積"],
      "duration": "PT0H10M0S",
      "uploadDate": "2018-11-06T04:00:00+09:00",
      "dateModified": "2022-08-17T12:51:02+09:00",
      "datePublished": "2018-11-06T09:30:00+09:00",
      "expires": null,
      "regionsAllowed": null,
      "usageInfo": "streaming",
      "bitrate": "512kbps",
      "height": 360,
      "width": 640,
      "hasPart": []
    }
  ]
}
```

- [ ] **Step 2: Write the failing tests**

Add to `video_test.go`:

```go
func TestSearchByKeyword(t *testing.T) {
	fixture, err := os.ReadFile("testdata/keyword.json")
	if err != nil {
		t.Fatalf("reading fixture: %v", err)
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/nfsvideos/keyword" {
			t.Errorf("path = %q", r.URL.Path)
		}
		q := r.URL.Query()
		if got := q.Get("keywords"); got != "空気" {
			t.Errorf("keywords = %q, want %q", got, "空気")
		}
		if got := q.Get("grades"); got != "24" {
			t.Errorf("grades = %q, want %q", got, "24")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer srv.Close()

	client, _ := NewClient("test-key", WithBaseURL(srv.URL))
	resp, err := client.SearchByKeyword(context.Background(), &KeywordParams{
		Keywords: "空気",
		Grades:   []string{"24"},
	})
	if err != nil {
		t.Fatalf("SearchByKeyword() error = %v", err)
	}

	if resp.TotalCount != 1 {
		t.Errorf("TotalCount = %d, want 1", resp.TotalCount)
	}
	if resp.Videos[0].Name != "とじこめられた空気" {
		t.Errorf("Name = %q", resp.Videos[0].Name)
	}
}

func TestSearchByKeyword_ExclusiveParams(t *testing.T) {
	client, _ := NewClient("test-key")
	_, err := client.SearchByKeyword(context.Background(), &KeywordParams{
		Keywords:     "空気",
		SubjectAreas: []string{"6"},
		Subjects:     []string{"411"},
	})
	if err == nil {
		t.Fatal("expected error when both SubjectAreas and Subjects are set")
	}
}
```

- [ ] **Step 3: Run tests to verify they pass**

The implementation is already in `video.go` from Task 5.

```bash
go test -run TestSearchByKeyword -v
```

Expected: PASS.

- [ ] **Step 4: Commit**

```bash
git add video_test.go testdata/keyword.json
git commit -m "test: add SearchByKeyword tests"
```

---

### Task 8: GetAll (TSV)

**Files:**
- Create: `testdata/all.tsv`
- Create: `all.go`
- Create: `all_test.go`

- [ ] **Step 1: Create test fixture**

Create `testdata/all.tsv` (tab-separated — each field below is separated by a tab character):

```tsv
id	contentType	name	nfsSeriesName	description	url	thumbnailUrl	grades	subjectAreas	subjects	curriculumStandardCode	keywords	duration	datePublished	expires	regionsAllowed
D0005110412_00000	ばんぐみ	とじこめられた空気	ふしぎエンドレス　理科４年	空気でっぽうの説明	https://www2.nhk.or.jp/school/watch/bangumi/?das_id=D0005110412_00000	https://www.nhk.or.jp/das/image/D0005110/D0005110412_00000_S_001.jpg	小4	理科		8260243111100000	空気,体積	PT0H10M0S	2018-11-06T09:30:00+09:00
D0005110413_00000	ばんぐみ	とじこめられた水	ふしぎエンドレス　理科４年	水でっぽうの説明	https://www2.nhk.or.jp/school/watch/bangumi/?das_id=D0005110413_00000	https://www.nhk.or.jp/das/image/D0005110/D0005110413_00000_S_001.jpg	小4	理科		8260243111100000	水,体積	PT0H10M0S	2018-11-06T09:30:00+09:00
```

- [ ] **Step 2: Write the failing test**

Create `all_test.go`:

```go
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
```

- [ ] **Step 3: Run test to verify it fails**

```bash
go test -run TestGetAll -v
```

Expected: FAIL — `GetAll` not defined.

- [ ] **Step 4: Write implementation**

Create `all.go`:

```go
package nhkschool

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strings"
)

// GetAll retrieves all video data as a TSV file and parses it into VideoSummary slices.
func (c *Client) GetAll(ctx context.Context) ([]*VideoSummary, error) {
	u := fmt.Sprintf("%s/nfsvideos/all/nhkforschool.tsv?apikey=%s", c.baseURL, c.apiKey)

	body, err := c.do(ctx, u)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bytes.NewReader(body))
	reader.Comma = '\t'
	reader.LazyQuotes = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("nhkschool: parsing TSV: %w", err)
	}

	if len(records) < 2 {
		return nil, nil
	}

	// Skip header row
	var videos []*VideoSummary
	for _, row := range records[1:] {
		if len(row) < 16 {
			continue
		}
		videos = append(videos, &VideoSummary{
			ID:             row[0],
			ContentType:    row[1],
			Name:           row[2],
			NFSSeriesName:  row[3],
			Description:    row[4],
			URL:            row[5],
			ThumbnailURL:   row[6],
			Grades:         splitNonEmpty(row[7]),
			SubjectAreas:   splitNonEmpty(row[8]),
			Subjects:       splitNonEmpty(row[9]),
			CurriculumCode: splitNonEmpty(row[10]),
			Keywords:       splitNonEmpty(row[11]),
			Duration:       row[12],
			DatePublished:  row[13],
			Expires:        row[14],
			RegionsAllowed: row[15],
		})
	}

	return videos, nil
}

// splitNonEmpty splits a comma-separated string, returning nil for empty input.
func splitNonEmpty(s string) []string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
```

- [ ] **Step 5: Run tests to verify they pass**

```bash
go test -run TestGetAll -v
```

Expected: PASS.

- [ ] **Step 6: Commit**

```bash
git add all.go all_test.go testdata/all.tsv
git commit -m "feat: add GetAll with TSV parsing"
```

---

### Task 9: CLI

**Files:**
- Create: `cmd/nhkschool/main.go`

- [ ] **Step 1: Create CLI**

Create `cmd/nhkschool/main.go`:

```go
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	nhkschool "github.com/kimihito/go-nhkschool"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	apiKey := os.Getenv("NHKSCHOOL_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: NHKSCHOOL_API_KEY environment variable is required")
		os.Exit(1)
	}

	client, err := nhkschool.NewClient(apiKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	subcmd := os.Args[1]

	switch subcmd {
	case "video":
		err = runVideo(ctx, client, os.Args[2:])
	case "list":
		err = runList(ctx, client, os.Args[2:])
	case "search":
		err = runSearch(ctx, client, os.Args[2:])
	case "all":
		err = runAll(ctx, client, os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintln(os.Stderr, `Usage: nhkschool <command> [arguments]

Commands:
  video <id>                    Get video details by ID
  list <cscode>                 List videos by curriculum standard code
  search <keywords> [options]   Search videos by keyword
  all [options]                 Get all video data

Environment:
  NHKSCHOOL_API_KEY             API key (required)`)
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

func runVideo(ctx context.Context, client *nhkschool.Client, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: nhkschool video <id>")
	}
	video, err := client.GetVideo(ctx, args[0])
	if err != nil {
		return err
	}
	return printJSON(video)
}

func runList(ctx context.Context, client *nhkschool.Client, args []string) error {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	includelower := fs.Int("includelower", -1, "Search range (0: exact, 1: include lower)")
	resultorder := fs.Int("resultorder", -1, "Sort order (0: popular, 1: updated)")
	contenttype := fs.Int("contenttype", -1, "Content type (0: all, 1: bangumi, 2: clip)")
	page := fs.Int("page", -1, "Page number")
	perpage := fs.Int("perpage", -1, "Results per page")

	if len(args) < 1 {
		return fmt.Errorf("usage: nhkschool list <cscode> [options]")
	}
	cscode := args[0]
	fs.Parse(args[1:])

	opts := &nhkschool.ListOptions{}
	hasOpts := false
	if *includelower >= 0 {
		opts.IncludeLower = includelower
		hasOpts = true
	}
	if *resultorder >= 0 {
		opts.ResultOrder = resultorder
		hasOpts = true
	}
	if *contenttype >= 0 {
		opts.ContentType = contenttype
		hasOpts = true
	}
	if *page >= 0 {
		opts.Page = page
		hasOpts = true
	}
	if *perpage >= 0 {
		opts.PerPage = perpage
		hasOpts = true
	}

	var optsPtr *nhkschool.ListOptions
	if hasOpts {
		optsPtr = opts
	}

	resp, err := client.ListByCSCode(ctx, cscode, optsPtr)
	if err != nil {
		return err
	}
	return printJSON(resp)
}

func runSearch(ctx context.Context, client *nhkschool.Client, args []string) error {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	grades := fs.String("grades", "", "Grades (comma-separated, e.g. 24,25)")
	subjectareas := fs.String("subjectareas", "", "Subject areas (comma-separated)")
	subjects := fs.String("subjects", "", "Subjects (comma-separated)")
	resultorder := fs.Int("resultorder", -1, "Sort order (0: popular, 1: updated)")
	contenttype := fs.Int("contenttype", -1, "Content type (0: all, 1: bangumi, 2: clip)")
	page := fs.Int("page", -1, "Page number")
	perpage := fs.Int("perpage", -1, "Results per page")

	if len(args) < 1 {
		return fmt.Errorf("usage: nhkschool search <keywords> [options]")
	}
	keywords := args[0]
	fs.Parse(args[1:])

	params := &nhkschool.KeywordParams{
		Keywords: keywords,
	}
	if *grades != "" {
		params.Grades = strings.Split(*grades, ",")
	}
	if *subjectareas != "" {
		params.SubjectAreas = strings.Split(*subjectareas, ",")
	}
	if *subjects != "" {
		params.Subjects = strings.Split(*subjects, ",")
	}
	if *resultorder >= 0 {
		params.ResultOrder = resultorder
	}
	if *contenttype >= 0 {
		params.ContentType = contenttype
	}
	if *page >= 0 {
		params.Page = page
	}
	if *perpage >= 0 {
		params.PerPage = perpage
	}

	resp, err := client.SearchByKeyword(ctx, params)
	if err != nil {
		return err
	}
	return printJSON(resp)
}

func runAll(ctx context.Context, client *nhkschool.Client, args []string) error {
	fs := flag.NewFlagSet("all", flag.ExitOnError)
	output := fs.String("output", "", "Output file path (default: stdout)")
	fs.Parse(args)

	videos, err := client.GetAll(ctx)
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(videos, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}

	if *output != "" {
		if err := os.WriteFile(*output, data, 0644); err != nil {
			return fmt.Errorf("writing file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Wrote %d videos to %s\n", len(videos), *output)
		return nil
	}

	fmt.Println(string(data))
	return nil
}
```

- [ ] **Step 2: Verify it compiles**

```bash
go build ./cmd/nhkschool
```

Expected: produces `nhkschool` binary, no errors.

- [ ] **Step 3: Verify help output**

```bash
./nhkschool
```

Expected: usage message printed to stderr, exit code 1.

- [ ] **Step 4: Clean up binary**

```bash
rm -f nhkschool
```

- [ ] **Step 5: Commit**

```bash
git add cmd/nhkschool/main.go
git commit -m "feat: add CLI with video, list, search, all subcommands"
```

---

### Task 10: Run all tests, final verification

**Files:** none (verification only)

- [ ] **Step 1: Run full test suite**

```bash
go test ./... -v
```

Expected: all tests pass.

- [ ] **Step 2: Run go vet**

```bash
go vet ./...
```

Expected: no issues.

- [ ] **Step 3: Verify build**

```bash
go build ./...
```

Expected: no errors.

- [ ] **Step 4: Verify directory rename instruction**

The module is `github.com/kimihito/go-nhkschool`. The directory will need to be renamed:

```bash
cd /Users/kimihito/ghq/github.com/kimihito
mv nhk-for-school-go go-nhkschool
```

Note: This step should be done last, after all commits, since it changes the working directory path. The git remote URL should also be updated if the GitHub repo is renamed.
