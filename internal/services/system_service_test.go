package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSystemService_GetSystemInfo(t *testing.T) {
	service := NewSystemService()

	t.Run("Success - Get system info", func(t *testing.T) {
		info := service.GetSystemInfo()

		assert.NotNil(t, info)
	})
}
