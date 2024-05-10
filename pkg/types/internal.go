package types

import (
	"fmt"
	"time"
)

type HashList struct {
	Hashes map[string]string `json:"hashes"`
}

type Post struct {
	ID                 string    `json:"id"`
	Draft              bool      `json:"draft"`
	Listed             bool      `json:"listed"`
	Title              string    `json:"title"`
	Slug               string    `json:"slug"`
	Content            string    `json:"content"`
	ContentHTML        string    `json:"contentHtml"`
	ContentHTMLPreview string    `json:"contentHtmlPreview"`
	Tags               []string  `json:"tags"`
	Published          time.Time `json:"published"`
}

type Like struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Title     string    `json:"title"`
	URL       string    `json:"url"`
}

// WrapErr wraps an error and returns a new one
func WrapErr(err error, message string) error {
	return fmt.Errorf("%s; %w", message, err)
}
