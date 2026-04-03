# go-nhkschool Design Spec

## Overview

NHK for School API v2 の Go クライアントライブラリ + CLI ツール。
NHK for School で公開されている教育動画に関するデータを取得するための4つのAPIエンドポイント全てをカバーする。

## Project Setup

- **Module**: `github.com/kimihito/go-nhkschool`
- **Package name**: `nhkschool`
- **Go version**: 1.24
- **Dependencies**: 標準ライブラリのみ（テスト含む）
- **Directory rename**: `nhk-for-school-go` → `go-nhkschool`

## API Endpoints

NHK for School API v2 は以下の4つのエンドポイントを提供する。

| API | Method | Endpoint | Format | Description |
|-----|--------|----------|--------|-------------|
| コンテンツリスト | GET | `/v2/nfsvideos/cscode/{cscode}` | JSON | 学習指導要領コードで動画リスト取得 |
| コンテンツ詳細 | GET | `/v2/nfsvideo/id/{id}` | JSON | 動画IDで詳細情報取得 |
| キーワード検索 | GET | `/v2/nfsvideos/keyword` | JSON | キーワード・学年・教科で検索 |
| 全件データ取得 | GET | `/v2/nfsvideos/all/nhkforschool.tsv` | TSV | 全動画データ一括取得 |

共通仕様:
- 認証: クエリパラメータ `apikey` によるAPIキー認証
- 利用制限: 全API合計で1ユーザーあたり3000回/日
- Base URL: `https://api.nhk.or.jp/school/v2`

## Architecture

### Directory Structure

```
go-nhkschool/
  nhkschool.go      // Client, NewClient, Option
  video.go          // GetVideo, ListByCSCode, SearchByKeyword
  all.go            // GetAll (TSV parse)
  error.go          // APIError
  types.go          // Video, VideoSummary, ListResponse, etc.
  testdata/          // JSON fixtures for tests
  cmd/nhkschool/
    main.go          // CLI entry point
```

フラットパッケージ構成。APIが4エンドポイントかつ全て動画関連のため、サブパッケージ分割は不要。

### Client

```go
type Client struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
}

func NewClient(apiKey string, opts ...Option) (*Client, error)

type Option func(*Client)
func WithHTTPClient(c *http.Client) Option
func WithBaseURL(url string) Option
```

- `apiKey` は必須。空文字の場合は error を返す
- `baseURL` のデフォルトは `https://api.nhk.or.jp/school/v2`
- Functional Options パターンでテスト時の差し替えを可能にする

### API Methods

```go
// コンテンツ詳細
func (c *Client) GetVideo(ctx context.Context, id string) (*Video, error)

// コンテンツリスト
func (c *Client) ListByCSCode(ctx context.Context, cscode string, opts *ListOptions) (*ListResponse, error)

// キーワード検索
func (c *Client) SearchByKeyword(ctx context.Context, params *KeywordParams) (*ListResponse, error)

// 全件データ取得
func (c *Client) GetAll(ctx context.Context) ([]*VideoSummary, error)
```

全メソッドの第一引数は `context.Context`。

### Request Parameters

```go
type ListOptions struct {
    IncludeLower *int // 0: 指定コードのみ, 1: 下位コード含む
    ResultOrder  *int // 0: よく見られている順, 1: 更新順
    ContentType  *int // 0: すべて, 1: ばんぐみ, 2: クリップ
    Page         *int
    PerPage      *int
}

type KeywordParams struct {
    Keywords     string   // 必須。"夏 AND 星座" や "虫 OR 体"
    Grades       []string // nil = 未指定
    SubjectAreas []string // nil = 未指定（Subjects と排他）
    Subjects     []string // nil = 未指定（SubjectAreas と排他）
    ResultOrder  *int
    ContentType  *int
    Page         *int
    PerPage      *int
}
```

Optional の扱い:
- `[]string` 型: `nil` で未指定を表現
- `int` 型: `0` が有効な値のためポインタ型にし、`nil` で未指定を表現
- `SubjectAreas` と `Subjects` の排他制御: リクエスト時に `if` でバリデーション

### Response Types

```go
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

type About struct {
    NFSSeriesName string `json:"nfsSeriesName"`
}

type CurriculumStandard struct {
    Version string `json:"curriculumStandardVersion"`
    NfsID   string `json:"curriculumStandardNfsId"`
    Code    string `json:"curriculumStandardCode"`
}

type Part struct {
    ClipNumber   int     `json:"clipNumber"`
    StartOffset  float64 `json:"startOffset"`
    EndOffset    float64 `json:"endOffset"`
    ThumbnailURL string  `json:"thumbnailUrl"`
    Name         string  `json:"name"`
}

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

type ListResponse struct {
    Videos     []*Video
    TotalCount int
    Page       int
    PerPage    int
}
```

APIレスポンスで `null` になりうるフィールド（`text`, `expires`, `regionsAllowed`）はポインタ型。
`Video` はAPIの全フィールドを忠実に表現する（クライアントの責務）。
`VideoSummary` は全件データ取得API（TSV形式）で返るフィールドのみ。

### Error Handling

```go
type APIError struct {
    StatusCode int
    Body       string
}

func (e *APIError) Error() string {
    return fmt.Sprintf("nhkschool: HTTP %d: %s", e.StatusCode, e.Body)
}
```

- HTTP 2xx 以外 → `*APIError` を返す
- ネットワークエラー、JSONパース失敗 → `fmt.Errorf` でラップして返す
- 呼び出し側は `errors.As` で `*APIError` を判別可能

## CLI

### Entry Point

`cmd/nhkschool/main.go`

### Subcommands

```bash
# コンテンツ詳細
nhkschool video D0005110412_00000

# コンテンツリスト
nhkschool list 8260243111100000

# キーワード検索
nhkschool search "夏 AND 星座" --grades 24,25

# 全件データ取得
nhkschool all --output nhkforschool.tsv
```

- APIキー: 環境変数 `NHKSCHOOL_API_KEY` から取得
- 出力: デフォルトJSON
- CLI フレームワーク: 不使用。`flag` パッケージ + 自前サブコマンド分岐

## Testing

- `net/http/httptest` でモックサーバーを立てる
- テスト用JSONフィクスチャは `testdata/` に配置
- 各APIメソッドごとに正常系・エラー系のテスト
- 外部テストライブラリは使用しない（標準ライブラリのみ）
