package shortener

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	userId = "90a1d2b1-f480-4f90-9aba-795a19273e58"
)

func TestShortLinkGenerator(t *testing.T) {
	originalUrl := "https://www.reddit.com/"
	shortenedUrl := GenerateShortLink(originalUrl, userId)
	assert.Equal(t, shortenedUrl, "dg7dvfd3")
}
