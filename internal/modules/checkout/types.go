package checkout

import "strings"

const IdempotencyKeyHeader = "Idempotency-Key"

func trim(value string) string {
	return strings.TrimSpace(value)
}
