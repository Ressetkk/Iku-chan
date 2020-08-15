package nhapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	NHentaiApiUrl       = "https://nhentai.net"
	NHentaiThumbnailUrl = "https://t.nhentai.net"
	NHentaiImagesUrl    = "https://i.nhentai.net"
)

type Image struct {
	Type   string `json:"t"`
	Width  int    `json:"w"`
	Height int    `json:"h"`
}

type Images struct {
	Pages     []Image `json:"pages"`
	Cover     Image   `json:"cover"`
	Thumbnail Image   `json:"thumbnail"`
}

type Tag struct {
	ID    int    `json:"id"`
	Type  string `json:"string"`
	Name  string `json:"name"`
	URL   string `json:"url"`
	Count int    `json:"count"`
}

type Result struct {
	ID              interface{}       `json:"id"` //nhentai API is so poorly written and results sometimes have this field as string
	MediaID         string            `json:"media_id"`
	Title           map[string]string `json:"title"`
	Scanlator       string            `json:"scanlator,omitempty"`
	UploadTimestamp int               `json:"uploaded"`
	NumOfPages      int               `json:"num_pages"`
	NumOfFavorites  int               `json:"num_favorites"`
	Images          Images            `json:"images"`
	Tags            []Tag             `json:"tags"`
}

type SearchResult struct {
	Results    []Result `json:"result"`
	NumOfPages int      `json:"num_pages,omitempty"`
	NumPerPage int      `json:"per_page,omitempty"`
}

// TODO (@Ressetkk): add logging
type Client struct {
	url    string
	client *http.Client
}

type Options struct {
	Timeout time.Duration
	Url     string
}

func New(o Options) *Client {
	c := &http.Client{Timeout: o.Timeout}
	return &Client{
		client: c,
		url:    o.Url,
	}
}

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

func (c Client) Search(query, sort string, page int) (*SearchResult, error) {
	if sort == "" {
		sort = "date"
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

// TODO (@Ressetkk): Implement random
func (c Client) Random() int {
	return 0
}

func (r Result) GetThumbnailUrl() string {
	return fmt.Sprintf("%v/galleries/%v/cover.jpg", NHentaiThumbnailUrl, r.MediaID)
}
