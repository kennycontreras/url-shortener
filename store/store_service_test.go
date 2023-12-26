package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	userUUID = "90a1d2b1-f480-4f90-9aba-795a19273e58"
)

var testStoreService = &StorageService{}

func init() {
	testStoreService = InitializeStore()
}

func TestStoreInit(t *testing.T) {
	assert.True(t, testStoreService.redisClient != nil)
}

func TestInsertionAndRetrieval(t *testing.T) {
	initialUrl := "https://www.mydealz.de/deals"
	shortUrl := "f480"

	SaveUrlMapping(shortUrl, initialUrl, userUUID)
	retrievedUrl := RetrieveInitialUrl(shortUrl)
	assert.Equal(t, initialUrl, retrievedUrl)
}
