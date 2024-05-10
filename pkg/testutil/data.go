package testutil

import (
	"math/rand"
	"time"

	"github.com/georgemblack/web-api/pkg/types"
	"github.com/google/uuid"
)

// NewLike generates a like with random test data.
func NewLike() types.Like {
	return types.Like{
		ID:        uuid.New().String(),
		Timestamp: time.Now(),
		Title:     "Like",
		URL:       "https://google.com",
	}
}

// NewLikes generates multiple likes with random test data.
func NewLikes() []types.Like {
	num := rand.Intn(10) + 1
	res := make([]types.Like, num)

	for i := range res {
		res[i] = NewLike()
	}

	return res
}
