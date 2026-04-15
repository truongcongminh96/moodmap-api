package httpx

import "strings"

func NormalizeQueryValue(value string) string {
	return strings.TrimSpace(value)
}
