package repo

import (
	"context"
	"fmt"

	firestore "cloud.google.com/go/firestore/apiv1"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/types"
)

type FirestoreService struct {
	client *firestore.Client
	config conf.Config
}

func NewFirestoreService(config conf.Config) (FirestoreService, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx)
	if err != nil {
		return FirestoreService{}, types.WrapErr(err, "failed to create firestore client")
	}

	return FirestoreService{
		client: client,
		config: config,
	}, nil
}

func (f *FirestoreService) GetLike(id string) (types.Like, error) {
	ctx := context.Background()
	req := firestorepb.GetDocumentRequest{
		Name: fmt.Sprintf("projects/%s/databases/%s/documents/web-likes/%s", f.config.GCloudProjectID, f.config.FirestoreDatabasename, id),
	}
	doc, err := f.client.GetDocument(ctx, &req)
	if err != nil {
		return types.Like{}, types.WrapErr(err, "failed to get like")
	}

	return types.Like{
		ID:        id,
		Timestamp: doc.Fields["timestamp"].GetTimestampValue().AsTime(),
		Title:     doc.Fields["title"].GetStringValue(),
		URL:       doc.Fields["url"].GetStringValue(),
	}, nil
}

func (f *FirestoreService) GetLikes() ([]types.Like, error) {
	ctx := context.Background()
	req := firestorepb.ListDocumentsRequest{
		Parent:       fmt.Sprintf("projects/%s/databases/%s/documents", f.config.GCloudProjectID, f.config.FirestoreDatabasename),
		CollectionId: "web-likes",
	}
	iter := f.client.ListDocuments(ctx, &req)

	likes := make([]types.Like, 0)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		likes = append(likes, toLike(doc))
	}

	return likes, nil
}

func toLike(doc *firestorepb.Document) types.Like {
	return types.Like{
		ID:        doc.Name,
		Timestamp: doc.Fields["timestamp"].GetTimestampValue().AsTime(),
		Title:     doc.Fields["title"].GetStringValue(),
		URL:       doc.Fields["url"].GetStringValue(),
	}
}

func (f *FirestoreService) GetPost(id string) (types.Post, error) {
	ctx := context.Background()
	req := firestorepb.GetDocumentRequest{
		Name: fmt.Sprintf("projects/%s/databases/%s/documents/web-posts/%s", f.config.GCloudProjectID, f.config.FirestoreDatabasename, id),
	}
	doc, err := f.client.GetDocument(ctx, &req)
	if err != nil {
		return types.Post{}, types.WrapErr(err, "failed to get post")
	}

	// Convert tags from firestore array to string array
	tags := doc.Fields["tags"].GetArrayValue().Values
	tagsStr := make([]string, len(tags))
	for i, v := range tags {
		tagsStr[i] = v.GetStringValue()
	}

	return toPost(doc), nil
}

func (f *FirestoreService) GetPosts() ([]types.Post, error) {
	ctx := context.Background()
	req := firestorepb.ListDocumentsRequest{
		Parent:       fmt.Sprintf("projects/%s/databases/%s/documents", f.config.GCloudProjectID, f.config.FirestoreDatabasename),
		CollectionId: "web-posts",
	}
	iter := f.client.ListDocuments(ctx, &req)

	posts := make([]types.Post, 0)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		posts = append(posts, toPost(doc))
	}

	return posts, nil
}

func toPost(doc *firestorepb.Document) types.Post {
	// Convert tags from firestore array to string array
	tags := doc.Fields["tags"].GetArrayValue().Values
	tagsStr := make([]string, len(tags))
	for i, v := range tags {
		tagsStr[i] = v.GetStringValue()
	}

	return types.Post{
		ID:        doc.Name,
		Draft:     doc.Fields["draft"].GetBooleanValue(),
		Listed:    doc.Fields["listed"].GetBooleanValue(),
		Title:     doc.Fields["title"].GetStringValue(),
		Slug:      doc.Fields["slug"].GetStringValue(),
		Content:   doc.Fields["content"].GetStringValue(),
		Tags:      tagsStr,
		Published: doc.Fields["published"].GetTimestampValue().AsTime(),
	}
}

func (f *FirestoreService) GetHashList() (types.HashList, error) {
	ctx := context.Background()
	req := firestorepb.GetDocumentRequest{
		Name: fmt.Sprintf("projects/%s/databases/%s/documents/web-metadata/hashes", f.config.GCloudProjectID, f.config.FirestoreDatabasename),
	}
	doc, err := f.client.GetDocument(ctx, &req)
	if err != nil {
		return types.HashList{}, types.WrapErr(err, "failed to get hash list")
	}

	// Convert from firestore doc to hash list
	hashes := make(map[string]string)
	for k, v := range doc.Fields {
		hashes[k] = v.GetStringValue()
	}

	return types.HashList{
		Hashes: hashes,
	}, nil
}

func (f *FirestoreService) Close() {
	f.client.Close()
}
