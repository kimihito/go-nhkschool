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
