package nhkschool

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"strings"
)

// GetAll retrieves all video data from the bulk TSV endpoint and parses it into
// a slice of VideoSummary.
//
// This endpoint returns all videos in a single TSV response, which may be large.
// Returns nil (not an error) if the API returns an empty dataset.
// Returns *APIError for non-2xx HTTP responses.
func (c *Client) GetAll(ctx context.Context) ([]*VideoSummary, error) {
	u := fmt.Sprintf("%s/nfsvideos/all/nhkforschool.tsv?apikey=%s", c.baseURL, c.apiKey)

	body, err := c.do(ctx, u)
	if err != nil {
		return nil, err
	}

	reader := csv.NewReader(bytes.NewReader(body))
	reader.Comma = '\t'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

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

// splitNonEmpty splits a comma-separated string into a slice of trimmed strings.
// Returns nil if the input is empty or whitespace-only.
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
