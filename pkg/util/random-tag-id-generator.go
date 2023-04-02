package util

import (
	"strings"

	"github.com/google/uuid"
)

func GenerateRandomTag() string {
	tagId := uuid.New().String()
	tagId = strings.ReplaceAll(tagId, "-", "")
	return tagId[:20]
}
