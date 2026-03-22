package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func New(baseURL, apiKey string) *Client {
	return &Client{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) do(req *http.Request, out any) error {
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("API error: %s", resp.Status)
	}

	if out != nil {
		if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}
	return nil
}

// Chunk represents a context chunk from the API
type Chunk struct {
	ID             string          `json:"id"`
	Title          string          `json:"title"`
	Content        string          `json:"content"`
	ChunkType      string          `json:"chunk_type"`
	Scope          string          `json:"scope"`
	QueryKey       string          `json:"query_key"`
	InjectAudience *InjectAudience `json:"inject_audience"`
	Score          float64         `json:"score,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	Version        int             `json:"version"`
}

type InjectAudience struct {
	Rules []InjectRule `json:"rules"`
}

type InjectRule struct {
	All            bool     `json:"all,omitempty"`
	AgentTypes     []string `json:"agent_types,omitempty"`
	AgentIDs       []string `json:"agent_ids,omitempty"`
	LifecycleTypes []string `json:"lifecycle_types,omitempty"`
	AgentTags      []string `json:"agent_tags,omitempty"`
	FocusTags      []string `json:"focus_tags,omitempty"`
}

type ChunkVersion struct {
	Version    int       `json:"version"`
	ChangeNote string    `json:"change_note,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type ChunkDetail struct {
	Chunk
	Versions []ChunkVersion `json:"versions"`
}

type SearchResponse struct {
	Results []Chunk `json:"results"`
}

// SearchChunks performs a semantic search
func (c *Client) SearchChunks(query, scope string) ([]Chunk, error) {
	params := url.Values{}
	params.Set("query", query)
	if scope != "" && scope != "all" {
		params.Set("scope", scope)
	}

	req, err := http.NewRequest("GET", c.baseURL+"/v1/context/search?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var result SearchResponse
	if err := c.do(req, &result); err != nil {
		return nil, err
	}
	return result.Results, nil
}

// Health checks API connectivity
func (c *Client) Health() error {
	req, err := http.NewRequest("GET", c.baseURL+"/health", nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}

// GetChunk fetches full chunk detail including versions
func (c *Client) GetChunk(id string) (*ChunkDetail, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/v1/context/"+id, nil)
	if err != nil {
		return nil, err
	}

	var result ChunkDetail
	if err := c.do(req, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// DeleteChunk deletes a chunk by ID
func (c *Client) DeleteChunk(id string) error {
	req, err := http.NewRequest("DELETE", c.baseURL+"/v1/context/"+id, nil)
	if err != nil {
		return err
	}
	return c.do(req, nil)
}
