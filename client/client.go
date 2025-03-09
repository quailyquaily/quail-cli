package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/tabwriter"
	"time"
)

type Client struct {
	AccessToken string
	APIBase     string
}

type CreateOrUpdateListPostPayload struct {
	Slug             string     `json:"slug"`
	CoverImageURL    string     `json:"cover_image_url"`
	Title            string     `json:"title"`
	Summary          string     `json:"summary"`
	Content          string     `json:"content"`
	Datetime         *time.Time `json:"datetime,omitempty"`
	FirstPublishedAt *time.Time `json:"first_published_at,omitempty"`
	Tags             string     `json:"tags"`
	Theme            string     `json:"theme"`
}

func New(accessToken, apiBase string) *Client {
	return &Client{
		AccessToken: accessToken,
		APIBase:     apiBase,
	}
}

func (c *Client) GetList(listID string) (any, error) {
	url := fmt.Sprintf("%s/lists/%s", c.APIBase, listID)
	return c.sendRequest("GET", url, nil)
}

func (c *Client) GetMe() (*UserResponse, error) {
	resp, err := c.sendRequest("GET", fmt.Sprintf("%s/users/me", c.APIBase), nil)
	if err != nil {
		return nil, err
	}
	ur := &UserResponse{}
	if err := json.Unmarshal(resp, ur); err != nil {
		return nil, err
	}
	return ur, nil
}

func (c *Client) sendRequest(method, url string, payload any) ([]byte, error) {
	var body []byte
	var err error

	if payload != nil {
		body, err = json.Marshal(payload)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.AccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		buf, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(buf))
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func PrettyPrintJSON(data any) {
	prettyJSON, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("Error generating pretty JSON: %v", err)
	}
	fmt.Println(string(prettyJSON))
}

func PrettyPrintUser(data *UserResponse) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "User:")
	fmt.Fprintf(w, "ID:\t%d\n", data.Data.ID)
	fmt.Fprintf(w, "Name:\t%s\n", data.Data.Name)
	fmt.Fprintf(w, "Email:\t%s\n", data.Data.Email)
	fmt.Fprintf(w, "Avatar Image URL:\t%s\n", data.Data.AvatarImageURL)
	fmt.Fprintf(w, "Bio:\t%s\n", data.Data.Bio)
	fmt.Fprintf(w, "Tagline:\t%s\n", data.Data.Tagline)
	fmt.Fprintf(w, "Created At:\t%s\n", data.Data.CreatedAt)
	fmt.Fprintf(w, "Social IDs:\n")
	for _, socialID := range data.Data.SocialIDs {
		fmt.Fprintf(w, "\t%s: %s\n", socialID.Name, socialID.Value)
	}
	w.Flush()
}

func PrettyPrintPost(data *PostResponse) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	fmt.Fprintln(w, "Post:")
	fmt.Fprintf(w, "ID:\t%d\n", data.Data.ID)
	fmt.Fprintf(w, "Slug:\t%s\n", data.Data.Slug)
	fmt.Fprintf(w, "Cover Image URL:\t%s\n", data.Data.CoverImageURL)
	fmt.Fprintf(w, "Title:\t%s\n", data.Data.Title)
	fmt.Fprintf(w, "Summary:\t%s\n", data.Data.Summary)
	fmt.Fprintf(w, "PublishedAt:\t%s\n", data.Data.PublishedAt)
	fmt.Fprintf(w, "First Published At:\t%s\n", data.Data.FirstPublishedAt)
	fmt.Fprintf(w, "Tags:\t%s\n", data.Data.Tags)
	fmt.Fprintf(w, "Theme:\t%s\n", data.Data.Theme)
	w.Flush()
}
