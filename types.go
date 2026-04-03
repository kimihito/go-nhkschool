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
	// IncludeLower controls whether lower-level curriculum codes are included.
	// 0 = exact match only, 1 = include lower codes.
	IncludeLower *int
	// ResultOrder specifies the sort order of results.
	// 0 = popular (default), 1 = updated.
	ResultOrder *int
	// ContentType filters the content type.
	// 0 = all, 1 = bangumi, 2 = clip.
	ContentType *int
	// Page is the 1-based page number.
	Page *int
	// PerPage is the number of results per page.
	PerPage *int
}

// KeywordParams are the parameters for SearchByKeyword.
//
// Grades and Subjects use NHK's numeric codes.
// For example, "24" means 小学4年 (4th grade elementary school).
type KeywordParams struct {
	// Keywords is the search query string. Must not be empty.
	Keywords string
	// Grades filters by grade using NHK numeric codes (e.g. "24" for 小学4年).
	// Multiple grades are OR-combined.
	Grades []string
	// SubjectAreas filters by subject area. Cannot be used together with Subjects.
	SubjectAreas []string
	// Subjects filters by subject. Cannot be used together with SubjectAreas.
	// Uses NHK numeric codes.
	Subjects []string
	// ResultOrder specifies the sort order of results.
	// 0 = popular (default), 1 = updated.
	ResultOrder *int
	// ContentType filters the content type.
	// 0 = all, 1 = bangumi, 2 = clip.
	ContentType *int
	// Page is the 1-based page number.
	Page *int
	// PerPage is the number of results per page.
	PerPage *int
}
