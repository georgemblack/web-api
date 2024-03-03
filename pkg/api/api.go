package api

import (
	"fmt"

	"github.com/georgemblack/web-api/pkg/conf"
	"github.com/georgemblack/web-api/pkg/repo"
	"github.com/georgemblack/web-api/pkg/types"
)

func Run() error {
	config, err := conf.LoadConfig()
	if err != nil {
		return types.WrapErr(err, "failed to load config")
	}

	firestore, err := repo.NewFirestoreService(config)
	if err != nil {
		return types.WrapErr(err, "failed to create firestore service")
	}

	list, err := firestore.GetHashList()
	if err != nil {
		return types.WrapErr(err, "failed to get hash list")
	}

	fmt.Println(list)
	return nil
}
