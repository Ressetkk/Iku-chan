package nhapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

const (
	NHentaiApiUrl       = "https://nhentai.net"
	NHentaiThumbnailUrl = "https://t.nhentai.net"
	NHentaiImagesUrl    = "https://i.nhentai.net"
)

type URL string

// Images struct stores information about one image in the Result
type Image struct {
	Type   string `json:"t"`
	Width  int    `json:"w"`
	Height int    `json:"h"`
}

// Images struct stores information about images in the Result
type Images struct {
	Pages     []Image `json:"pages"`
	Cover     Image   `json:"cover"`
	Thumbnail Image   `json:"thumbnail"`
}

// Tag struct represents one of Result tags
type Tag struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	URL   URL    `json:"url"`
	Count int    `json:"count"`
}

// Result struct describes one result from API
type Result struct {
	ID              interface{}       `json:"id"` //nhentai API is so poorly written and results sometimes have this field as string
	MediaID         string            `json:"media_id"`
	Title           map[string]string `json:"title"`
	Scanlator       string            `json:"scanlator,omitempty"`
	UploadTimestamp int64             `json:"upload_date"`
	NumOfPages      int               `json:"num_pages"`
	NumOfFavorites  int               `json:"num_favorites"`
	Images          Images            `json:"images"`
	Tags            []Tag             `json:"tags"`
}

// SearchResult is a struct that describes complete search result.
type SearchResult struct {
	Results    []Result `json:"result"`
	NumOfPages int      `json:"num_pages,omitempty"`
	NumPerPage int      `json:"per_page,omitempty"`
}

// TODO (@Ressetkk): add logging
// Client struct defines NHentai API Client
type Client struct {
	url    string
	client *http.Client
}

// Options struct define Client options
type Options struct {
	Timeout time.Duration
	Url     string
}

// New returns new NHentai Client
func New(o Options) *Client {
	c := &http.Client{Timeout: o.Timeout}
	return &Client{
		client: c,
		url:    o.Url,
	}
}

// Get returns Result from NHentai.
func (c Client) Get(id int) (*Result, error) {
	uri := fmt.Sprintf("%v/api/gallery/%v", c.url, id)
	r, err := c.client.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("nhentai api fetch error: %w", err)
	}
	response := new(Result)
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("response decoding error: %w", err)
	}
	defer r.Body.Close()
	return response, nil
}

// Search returns list of Result with the search query results.
// Supports pagination.
func (c Client) Search(query, sort string, page int) (*SearchResult, error) {
	if sort == "" {
		sort = "popular"
	}
	if page == 0 {
		page = 1
	}

	uri := fmt.Sprintf("%v/api/galleries/search?%v",
		c.url, fmt.Sprintf("query=%v&sort=%v&page=%v", url.QueryEscape(query), sort, page))
	r, err := c.client.Get(uri)
	if err != nil {
		return nil, fmt.Errorf("nhentai api fetch error: %w", err)
	}
	response := new(SearchResult)
	err = json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		return nil, fmt.Errorf("response decoding error: %w", err)
	}
	defer r.Body.Close()
	return response, nil
}

// Random returns random id from NHentai.
func (c Client) Random() (int, error) {
	resp, err := c.client.Head(c.url + "/random")
	if err != nil {
		return -1, fmt.Errorf("nhentai api fetch error: %w", err)
	}
	r, err := regexp.Compile("(/)|(g)")
	if err != nil {
		return -1, err
	}
	strId := r.ReplaceAllString(resp.Request.URL.Path, "")
	id, err := strconv.Atoi(strId)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (r Result) GetCoverThumbnail() string {
	return fmt.Sprintf("%v/galleries/%v/cover.%v", NHentaiThumbnailUrl, r.MediaID, r.Images.Cover.FileType())
}

func (u URL) Full() string {
	return NHentaiApiUrl + string(u)
}

func (i Image) FileType() string {
	switch i.Type {
	case "j":
		return "jpg"
	case "p":
		return "png"
	default:
		return "jpg"
	}
}
