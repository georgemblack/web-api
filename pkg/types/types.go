package types

import (
	"fmt"
	"time"
)

type HashList struct {
	Hashes map[string]string `json:"hashes"`
}

type Post struct {
	ID                 string
	Draft              bool
	Listed             bool
	Title              string
	Slug               string
	Content            string
	ContentHTML        string
	ContentHTMLPreview string
	Tags               []string
	Published          time.Time
}

type Like struct {
	ID        string
	Timestamp time.Time
	Title     string
	URL       string
}

// WrapErr wraps an error and returns a new one
func WrapErr(err error, message string) error {
	return fmt.Errorf("%s; %w", message, err)
}
