package main

import (
	"fmt"
	"strings"
)

const cloudminiErrorFixtureHeader = "X-Cloudmini-Error-Fixture"

func validateCloudminiErrorFixturePath(envKey string, path string, requiredTerms ...string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return fmt.Errorf("%s is required for fixture evidence", envKey)
	}
	if !strings.HasPrefix(path, "/api/v3/") {
		return fmt.Errorf("%s must start with /api/v3/", envKey)
	}
	if strings.ContainsAny(path, "?#") {
		return fmt.Errorf("%s must not contain query or fragment", envKey)
	}
	for _, term := range requiredTerms {
		if !strings.Contains(path, term) {
			return fmt.Errorf("%s must be a side-effect-free fixture path", envKey)
		}
	}
	if path == "/api/v3/proxies" || strings.HasPrefix(path, "/api/v3/capacity/") {
		return fmt.Errorf("%s must be a side-effect-free fixture path", envKey)
	}
	return nil
}
