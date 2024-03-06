package repo

import (
	"strings"

	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"github.com/georgemblack/web-api/pkg/types"
)

func toLike(doc *firestorepb.Document) types.Like {
	return types.Like{
		ID:        id(doc),
		Timestamp: doc.Fields["timestamp"].GetTimestampValue().AsTime(),
		Title:     doc.Fields["title"].GetStringValue(),
		URL:       doc.Fields["url"].GetStringValue(),
	}
}

func toPost(doc *firestorepb.Document) types.Post {
	// Convert tags from firestore array to string array
	tags := doc.Fields["tags"].GetArrayValue().Values
	tagsStr := make([]string, len(tags))
	for i, v := range tags {
		tagsStr[i] = v.GetStringValue()
	}

	return types.Post{
		ID:        id(doc),
		Draft:     doc.Fields["draft"].GetBooleanValue(),
		Listed:    doc.Fields["listed"].GetBooleanValue(),
		Title:     doc.Fields["title"].GetStringValue(),
		Slug:      doc.Fields["slug"].GetStringValue(),
		Content:   doc.Fields["content"].GetStringValue(),
		Tags:      tagsStr,
		Published: doc.Fields["published"].GetTimestampValue().AsTime(),
	}
}

func id(doc *firestorepb.Document) string {
	split := strings.Split(doc.Name, "/")
	return split[len(split)-1]
}
