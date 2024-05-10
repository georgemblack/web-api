package repo

import (
	"context"
	"fmt"
	"time"

	firestore "cloud.google.com/go/firestore/apiv1"
	"cloud.google.com/go/firestore/apiv1/firestorepb"
	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/types"
	"github.com/google/uuid"
)

type Firestore struct {
	client *firestore.Client
	config conf.Config
}

func NewFirestoreService(config conf.Config) (*Firestore, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx)
	if err != nil {
		return &Firestore{}, types.WrapErr(err, "failed to create firestore client")
	}

	return &Firestore{
		client: client,
		config: config,
	}, nil
}

func (f *Firestore) GetLike(id string) (types.Like, error) {
	ctx := context.Background()
	req := firestorepb.GetDocumentRequest{
		Name: fmt.Sprintf("projects/%s/databases/%s/documents/web-likes/%s", f.config.GCloudProjectID, f.config.FirestoreDatabasename, id),
	}
	doc, err := f.client.GetDocument(ctx, &req)
	if err != nil {
		return types.Like{}, types.WrapErr(err, "failed to get like")
	}

	return docToLike(doc), nil
}

func (f *Firestore) GetLikes() ([]types.Like, error) {
	ctx := context.Background()
	req := firestorepb.ListDocumentsRequest{
		Parent:       fmt.Sprintf("projects/%s/databases/%s/documents", f.config.GCloudProjectID, f.config.FirestoreDatabasename),
		CollectionId: "web-likes",
		OrderBy:      "timestamp desc",
	}
	iter := f.client.ListDocuments(ctx, &req)

	likes := make([]types.Like, 0)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		likes = append(likes, docToLike(doc))
	}

	return likes, nil
}

func (f *Firestore) AddLike(like types.Like) (string, error) {
	ctx := context.Background()
	id := uuid.New().String()
	req := firestorepb.CreateDocumentRequest{
		Parent:       fmt.Sprintf("projects/%s/databases/%s/documents", f.config.GCloudProjectID, f.config.FirestoreDatabasename),
		CollectionId: "web-likes",
		DocumentId:   id,
		Document:     likeToDoc(like),
	}
	_, err := f.client.CreateDocument(ctx, &req)
	if err != nil {
		return "", types.WrapErr(err, "failed to create like")
	}

	return id, nil
}

func (f *Firestore) DeleteLike(id string) error {
	ctx := context.Background()
	req := firestorepb.DeleteDocumentRequest{
		Name: fmt.Sprintf("projects/%s/databases/%s/documents/web-likes/%s", f.config.GCloudProjectID, f.config.FirestoreDatabasename, id),
	}
	err := f.client.DeleteDocument(ctx, &req)
	if err != nil {
		return types.WrapErr(err, "failed to delete like")
	}

	return nil
}

func (f *Firestore) GetPost(id string) (types.Post, error) {
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

	return docToPost(doc), nil
}

type PostFilters struct {
	Listed    *bool
	Published *bool
}

func (f *Firestore) GetPosts(filters PostFilters) ([]types.Post, error) {
	ctx := context.Background()
	req := firestorepb.ListDocumentsRequest{
		Parent:       fmt.Sprintf("projects/%s/databases/%s/documents", f.config.GCloudProjectID, f.config.FirestoreDatabasename),
		CollectionId: "web-posts",
		OrderBy:      "published desc",
	}
	iter := f.client.ListDocuments(ctx, &req)

	posts := make([]types.Post, 0)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		post := docToPost(doc)

		// 'Listed' filter verifies the post is marked as listed
		if filters.Listed != nil {
			if *filters.Listed && !post.Listed {
				continue
			}
			if !*filters.Listed && post.Listed {
				continue
			}
		}

		// 'Published' filter checks:
		//	1. Whether a post is a draft
		//	2. Whether the post's publsihed date is in the future
		if filters.Published != nil {
			if *filters.Published && (post.Draft || post.Published.After(time.Now())) {
				continue
			}
			if !*filters.Published && (!post.Draft && post.Published.Before(time.Now())) {
				continue
			}
		}
		posts = append(posts, docToPost(doc))
	}

	return posts, nil
}

func (f *Firestore) AddPost(post types.Post) (string, error) {
	ctx := context.Background()
	id := uuid.New().String()
	req := firestorepb.CreateDocumentRequest{
		Parent:       fmt.Sprintf("projects/%s/databases/%s/documents", f.config.GCloudProjectID, f.config.FirestoreDatabasename),
		CollectionId: "web-posts",
		DocumentId:   id,
		Document:     postToDoc(post),
	}
	for i, v := range post.Tags {
		req.Document.Fields["tags"].GetArrayValue().Values[i] = &firestorepb.Value{ValueType: &firestorepb.Value_StringValue{StringValue: v}}
	}
	_, err := f.client.CreateDocument(ctx, &req)
	if err != nil {
		return "", types.WrapErr(err, "failed to create post")
	}

	return id, nil
}

func (f *Firestore) UpdatePost(post types.Post) error {
	ctx := context.Background()

	// Convert tags from string array to firestore array
	tags := make([]*firestorepb.Value, len(post.Tags))
	for i, v := range post.Tags {
		tags[i] = &firestorepb.Value{ValueType: &firestorepb.Value_StringValue{StringValue: v}}
	}

	req := firestorepb.UpdateDocumentRequest{
		Document: &firestorepb.Document{
			Name:   fmt.Sprintf("projects/%s/databases/%s/documents/web-posts/%s", f.config.GCloudProjectID, f.config.FirestoreDatabasename, post.ID),
			Fields: postToDoc(post).Fields,
		},
	}
	_, err := f.client.UpdateDocument(ctx, &req)
	if err != nil {
		return types.WrapErr(err, "failed to update post")
	}

	return nil
}

func (f *Firestore) DeletePost(id string) error {
	ctx := context.Background()
	req := firestorepb.DeleteDocumentRequest{
		Name: fmt.Sprintf("projects/%s/databases/%s/documents/web-posts/%s", f.config.GCloudProjectID, f.config.FirestoreDatabasename, id),
	}
	err := f.client.DeleteDocument(ctx, &req)
	if err != nil {
		return types.WrapErr(err, "failed to delete like")
	}

	return nil
}

func (f *Firestore) GetHashList() (types.HashList, error) {
	ctx := context.Background()
	req := firestorepb.GetDocumentRequest{
		Name: fmt.Sprintf("projects/%s/databases/%s/documents/web-metadata/hashes", f.config.GCloudProjectID, f.config.FirestoreDatabasename),
	}
	doc, err := f.client.GetDocument(ctx, &req)
	if err != nil {
		return types.HashList{}, types.WrapErr(err, "failed to get hash list")
	}

	return docToHash(doc), nil
}

func (f *Firestore) UpdateHashList(hashList types.HashList) error {
	ctx := context.Background()
	req := firestorepb.UpdateDocumentRequest{
		Document: &firestorepb.Document{
			Name:   fmt.Sprintf("projects/%s/databases/%s/documents/web-metadata/hashes", f.config.GCloudProjectID, f.config.FirestoreDatabasename),
			Fields: hashToDoc(hashList).Fields,
		},
	}
	_, err := f.client.UpdateDocument(ctx, &req)
	if err != nil {
		return types.WrapErr(err, "failed to update hash list")
	}

	return nil
}

func (f *Firestore) Close() {
	f.client.Close()
}
