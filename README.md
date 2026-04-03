# go-nhkschool

[![Go Reference](https://pkg.go.dev/badge/github.com/kimihito/go-nhkschool.svg)](https://pkg.go.dev/github.com/kimihito/go-nhkschool)
[![Go Report Card](https://goreportcard.com/badge/github.com/kimihito/go-nhkschool)](https://goreportcard.com/report/github.com/kimihito/go-nhkschool)

NHK for School API v2 の Go クライアントライブラリ + CLI ツール。NHK for School で公開されている教育動画に関するデータを取得するための4つの API エンドポイントを全てカバーしています。

## Features

- **4つの API エンドポイントを完全サポート**
  - コンテンツ詳細取得 (`GetVideo`)
  - 学習指導要領コードによるリスト取得 (`ListByCSCode`)
  - キーワード検索 (`SearchByKeyword`)
  - 全件 TSV データ取得 (`GetAll`)
- **ゼロ依存** — 標準ライブラリのみ
- **コンテキストサポート** — 全メソッドが `context.Context` を受け取る
- **Functional Options** — `WithHTTPClient`, `WithBaseURL` で柔軟にカスタマイズ可能
- **CLI** — `video`, `list`, `search`, `all` の4サブコマンド
- **エラーハンドリング** — `*APIError` で HTTP ステータスコードとレスポンスボディを取得可能

## Installation

### Library

```bash
go get github.com/kimihito/go-nhkschool
```

### CLI

```bash
go install github.com/kimihito/go-nhkschool/cmd/nhkschool@latest
```

## Quick Start

### Library

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/kimihito/go-nhkschool"
)

func main() {
    client, err := nhkschool.NewClient("YOUR_API_KEY")
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // 動画詳細取得
    video, err := client.GetVideo(ctx, "D0005110412_00000")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(video.Name)

    // 学習指導要領コードで検索
    resp, err := client.ListByCSCode(ctx, "8260243111100000", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%d件ヒット\n", resp.TotalCount)

    // キーワード検索
    results, err := client.SearchByKeyword(ctx, &nhkschool.KeywordParams{
        Keywords: "空気",
        Grades:   []string{"24"}, // 小学4年
    })
    if err != nil {
        log.Fatal(err)
    }
    for _, v := range results.Videos {
        fmt.Println(v.Name)
    }
}
```

### CLI

環境変数 `NHKSCHOOL_API_KEY` に API キーを設定してください。

```bash
export NHKSCHOOL_API_KEY="YOUR_API_KEY"

# 動画詳細取得
nhkschool video D0005110412_00000

# 学習指導要領コードで検索
nhkschool list 8260243111100000

# キーワード検索
nhkschool search "空気" --grades 24

# 全件データ取得
nhkschool all
nhkschool all --output all.json
```

## API Methods

| Method | Description |
|--------|-------------|
| `GetVideo(ctx, id)` | 動画IDで詳細情報を取得 |
| `ListByCSCode(ctx, cscode, opts)` | 学習指導要領コードで動画リストを取得 |
| `SearchByKeyword(ctx, params)` | キーワード・学年・教科で検索 |
| `GetAll(ctx)` | 全動画データを TSV で一括取得 |

### Options

```go
// GetVideo: オプションなし
video, _ := client.GetVideo(ctx, "D0005110412_00000")

// ListByCSCode: オプション
resp, _ := client.ListByCSCode(ctx, "8260243111100000", &nhkschool.ListOptions{
    IncludeLower: ptr(1), // 下位コードを含む
    PerPage:      ptr(50),
})

// SearchByKeyword: フィルター
results, _ := client.SearchByKeyword(ctx, &nhkschool.KeywordParams{
    Keywords:     "空気",
    Grades:       []string{"24", "25"},
    SubjectAreas: []string{"6"}, // 理科
})
```

## Error Handling

```go
_, err := client.GetVideo(ctx, "INVALID")
if err != nil {
    var apiErr *nhkschool.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("HTTP %d: %s\n", apiErr.StatusCode, apiErr.Body)
    }
}
```

## CLI Subcommands

```
Usage: nhkschool <command> [arguments]

Commands:
  video <id>                    Get video details by ID
  list <cscode>                 List videos by curriculum standard code
  search <keywords> [options]   Search videos by keyword
  all [options]                 Get all video data

Flags:
  --includelower int   Search range (0: exact, 1: include lower)
  --resultorder int    Sort order (0: popular, 1: updated)
  --contenttype int    Content type (0: all, 1: bangumi, 2: clip)
  --page int           Page number
  --perpage int        Results per page
  --grades string      Grades (comma-separated, e.g. 24,25)
  --subjectareas string Subject areas (comma-separated)
  --subjects string    Subjects (comma-separated)
  --output string      Output file path (for 'all' command)

Environment:
  NHKSCHOOL_API_KEY     API key (required)
```

## API Rate Limit

NHK for School API v2 には **3000回/日** の利用制限があります（全エンドポイント合計）。

## Getting an API Key

API キーは [NHK for School API](https://school-api-portal.nhk.or.jp/) から取得してください。

## Requirements

- Go 1.24+

## License

MIT
