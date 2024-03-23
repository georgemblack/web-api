package repo

import (
	"strings"

	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"github.com/georgemblack/web-api/pkg/types"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func docToLike(doc *firestorepb.Document) types.Like {
	return types.Like{
		ID:        id(doc),
		Timestamp: doc.Fields["timestamp"].GetTimestampValue().AsTime(),
		Title:     doc.Fields["title"].GetStringValue(),
		URL:       doc.Fields["url"].GetStringValue(),
	}
}

func likeToDoc(like types.Like) *firestorepb.Document {
	return &firestorepb.Document{
		Fields: map[string]*firestorepb.Value{
			"timestamp": {ValueType: &firestorepb.Value_TimestampValue{TimestampValue: timestamppb.New(like.Timestamp)}},
			"title":     {ValueType: &firestorepb.Value_StringValue{StringValue: like.Title}},
			"url":       {ValueType: &firestorepb.Value_StringValue{StringValue: like.URL}},
		},
	}
}

func docToPost(doc *firestorepb.Document) types.Post {
	// Convert tags from firestore array to string array
	tags := doc.Fields["tags"].GetArrayValue().Values
	tagsStr := make([]string, len(tags))
	for i, v := range tags {
		tagsStr[i] = v.GetStringValue()
	}

	return types.Post{
		ID:                 id(doc),
		Draft:              doc.Fields["draft"].GetBooleanValue(),
		Listed:             doc.Fields["listed"].GetBooleanValue(),
		Title:              doc.Fields["title"].GetStringValue(),
		Slug:               doc.Fields["slug"].GetStringValue(),
		Content:            doc.Fields["content"].GetStringValue(),
		ContentHTML:        doc.Fields["contentHtml"].GetStringValue(),
		ContentHTMLPreview: doc.Fields["contentHtmlPreview"].GetStringValue(),
		Tags:               tagsStr,
		Published:          doc.Fields["published"].GetTimestampValue().AsTime(),
	}
}

func postToDoc(post types.Post) *firestorepb.Document {
	// Convert tags from string array to firestore array
	tags := make([]*firestorepb.Value, len(post.Tags))
	for i, v := range post.Tags {
		tags[i] = &firestorepb.Value{ValueType: &firestorepb.Value_StringValue{StringValue: v}}
	}

	return &firestorepb.Document{
		Fields: map[string]*firestorepb.Value{
			"draft":              {ValueType: &firestorepb.Value_BooleanValue{BooleanValue: post.Draft}},
			"listed":             {ValueType: &firestorepb.Value_BooleanValue{BooleanValue: post.Listed}},
			"title":              {ValueType: &firestorepb.Value_StringValue{StringValue: post.Title}},
			"slug":               {ValueType: &firestorepb.Value_StringValue{StringValue: post.Slug}},
			"content":            {ValueType: &firestorepb.Value_StringValue{StringValue: post.Content}},
			"contentHtml":        {ValueType: &firestorepb.Value_StringValue{StringValue: post.ContentHTML}},
			"contentHtmlPreview": {ValueType: &firestorepb.Value_StringValue{StringValue: post.ContentHTMLPreview}},
			"tags":               {ValueType: &firestorepb.Value_ArrayValue{ArrayValue: &firestorepb.ArrayValue{Values: tags}}},
			"published":          {ValueType: &firestorepb.Value_TimestampValue{TimestampValue: timestamppb.New(post.Published)}},
		},
	}
}

func docToHash(doc *firestorepb.Document) types.HashList {
	hashes := make(map[string]string)
	for k, v := range doc.Fields {
		hashes[k] = v.GetStringValue()
	}

	return types.HashList{
		Hashes: hashes,
	}
}

func hashToDoc(hashList types.HashList) *firestorepb.Document {
	fields := make(map[string]*firestorepb.Value)
	for k, v := range hashList.Hashes {
		fields[k] = &firestorepb.Value{ValueType: &firestorepb.Value_StringValue{StringValue: v}}
	}

	return &firestorepb.Document{
		Fields: fields,
	}
}

func id(doc *firestorepb.Document) string {
	split := strings.Split(doc.Name, "/")
	return split[len(split)-1]
}
