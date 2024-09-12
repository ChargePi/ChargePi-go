package util

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateRandomTag() string {
	tagId := uuid.New().String()
	tagId = strings.ReplaceAll(tagId, "-", "")
	// With OCPP 1.6, the maximum length for a tag ID is 20 characters
	return tagId[:20]
}
