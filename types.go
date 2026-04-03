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
