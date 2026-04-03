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
	if params == nil {
		return nil, fmt.Errorf("nhkschool: params must not be nil")
	}
	if params.Keywords == "" {
		return nil, fmt.Errorf("nhkschool: keywords must not be empty")
	}
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
