package helper

import (
	"strings"

	"github.com/kumparan/go-utils"
	"github.com/luckyAkbar/jatis-redis-training/internal/config"
	"github.com/mattheath/base62"
)

func GenerateToken(uniqueID int64) string {
	sb := strings.Builder{}

	encodedID := base62.EncodeInt64(uniqueID)
	sb.WriteString(encodedID)
	sb.WriteString("___")

	randString := utils.GenerateRandomAlphanumeric(config.DefaultTokenLength)
	sb.WriteString(randString)

	return sb.String()
}
