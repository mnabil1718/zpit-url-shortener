package qr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateQR_ReturnBytes(t *testing.T) {
	b, err := GenerateQR("https://google.com")
	assert.NoError(t, err)
	assert.NotEmpty(t, b)
}
