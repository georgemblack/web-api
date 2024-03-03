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

type HashList struct {
	Hashes map[string]string `json:"hashes"`
}

func (f *FirestoreService) GetHashList() (HashList, error) {
	ctx := context.Background()
	req := firestorepb.GetDocumentRequest{
		Name: fmt.Sprintf("projects/%s/databases/%s/documents/web-metadata/hashes", f.config.GCloudProjectID, f.config.FirestoreDatabasename),
	}
	doc, err := f.client.GetDocument(ctx, &req)
	if err != nil {
		return HashList{}, types.WrapErr(err, "failed to get hash list")
	}

	// Convert from firestore doc to hash list
	hashes := make(map[string]string)
	for k, v := range doc.Fields {
		hashes[k] = v.GetStringValue()
	}

	return HashList{
		Hashes: hashes,
	}, nil
}

func (f *FirestoreService) Close() {
	f.client.Close()
}
