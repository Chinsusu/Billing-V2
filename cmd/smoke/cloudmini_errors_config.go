package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

func cloudminiErrorEvidenceConfigFromEnv() (cloudminiErrorEvidenceConfig, error) {
	appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if appEnv == "" {
		return cloudminiErrorEvidenceConfig{}, fmt.Errorf("APP_ENV is required")
	}
	switch appEnv {
	case "local", "dev", "staging", "sandbox":
	case "prod", "production":
		return cloudminiErrorEvidenceConfig{}, fmt.Errorf("refusing to run cloudmini error evidence with APP_ENV=%s", appEnv)
	default:
		return cloudminiErrorEvidenceConfig{}, fmt.Errorf("APP_ENV must be local, dev, staging, or sandbox")
	}
	if os.Getenv("BILLING_CLOUDMINI_ERROR_EVIDENCE_APPROVED") != "yes" {
		return cloudminiErrorEvidenceConfig{}, fmt.Errorf("BILLING_CLOUDMINI_ERROR_EVIDENCE_APPROVED=yes is required")
	}
	for _, key := range []string{
		"CLOUDMINI_SOURCE_ACCOUNT_OWNER",
		"CLOUDMINI_ENGINEERING_OWNER",
		"CLOUDMINI_OPS_OWNER",
		"CLOUDMINI_SECURITY_OWNER",
		"CLOUDMINI_CLEANUP_OWNER",
		"CLOUDMINI_REVIEWER_SIGNOFF",
		"CLOUDMINI_PILOT_STOP_CONDITION",
		"CLOUDMINI_PILOT_READONLY_EVIDENCE_REF",
	} {
		if err := requireCloudminiEvidenceFilled(key); err != nil {
			return cloudminiErrorEvidenceConfig{}, err
		}
	}
	config := cloudminiErrorEvidenceConfig{
		AppEnv:                    appEnv,
		BaseURL:                   strings.TrimSpace(os.Getenv("CLOUDMINI_V3_BASE_URL")),
		APIToken:                  strings.TrimSpace(os.Getenv("CLOUDMINI_V3_API_TOKEN")),
		IncludeCreate:             os.Getenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_INVALID_CREATE") == "yes",
		IncludePermissionDenied:   os.Getenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_PERMISSION_DENIED") == "yes",
		IncludeOutOfCapacity:      os.Getenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_OUT_OF_CAPACITY") == "yes",
		IncludeRateLimited:        os.Getenv("CLOUDMINI_ERROR_EVIDENCE_ALLOW_RATE_LIMITED") == "yes",
		PermissionKeyManagementOK: strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MANAGEMENT_APPROVED")),
		PermissionKeyMaxCreate:    strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MAX_CREATE")),
		OutOfCapacityApproved:     strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_APPROVED")),
		OutOfCapacityMaxAttempts:  strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_MAX_RESERVATIONS")),
		OutOfCapacityKind:         strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_KIND")),
		RateLimitApproved:         strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_APPROVED")),
		RateLimitMaxRequests:      strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_MAX_REQUESTS")),
		RateLimitFixturePath:      strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_FIXTURE_PATH")),
	}
	if ttlRaw := strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_TTL_SECONDS")); ttlRaw != "" {
		ttlSeconds, err := strconv.Atoi(ttlRaw)
		if err != nil {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_TTL_SECONDS must be an integer")
		}
		config.OutOfCapacityTTLSeconds = ttlSeconds
	}
	if config.BaseURL == "" {
		return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_V3_BASE_URL is required")
	}
	if config.APIToken == "" {
		return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_V3_API_TOKEN is required")
	}
	if _, err := resolveCloudminiErrorEvidenceURL(config.BaseURL, "/api/v3/capabilities"); err != nil {
		return cloudminiErrorEvidenceConfig{}, err
	}
	if config.IncludeCreate {
		if strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_MUTATING_ROUTE_APPROVED")) != "yes" {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_MUTATING_ROUTE_APPROVED=yes is required for malformed create validation evidence")
		}
		if strings.TrimSpace(os.Getenv("CLOUDMINI_ERROR_EVIDENCE_MAX_CREATE_ATTEMPTS")) != "1" {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_MAX_CREATE_ATTEMPTS must be 1")
		}
	}
	if config.IncludePermissionDenied {
		if config.PermissionKeyManagementOK != "yes" {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MANAGEMENT_APPROVED=yes is required for permission-denied evidence")
		}
		if config.PermissionKeyMaxCreate != "1" {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_PERMISSION_KEY_MAX_CREATE must be 1")
		}
	}
	if config.IncludeOutOfCapacity {
		if config.OutOfCapacityApproved != "yes" {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_APPROVED=yes is required for out-of-capacity evidence")
		}
		if config.OutOfCapacityMaxAttempts != "1" {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_MAX_RESERVATIONS must be 1")
		}
		if config.OutOfCapacityKind != "ipv4_dc" && config.OutOfCapacityKind != "residential" {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_KIND must be ipv4_dc or residential")
		}
		if config.OutOfCapacityTTLSeconds <= 0 || config.OutOfCapacityTTLSeconds > 60 {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_OUT_OF_CAPACITY_TTL_SECONDS must be between 1 and 60")
		}
	}
	if config.IncludeRateLimited {
		if config.RateLimitApproved != "yes" {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_APPROVED=yes is required for rate-limit evidence")
		}
		if config.RateLimitMaxRequests != "1" {
			return cloudminiErrorEvidenceConfig{}, fmt.Errorf("CLOUDMINI_ERROR_EVIDENCE_RATE_LIMIT_MAX_REQUESTS must be 1")
		}
		if err := validateCloudminiRateLimitFixturePath(config.RateLimitFixturePath); err != nil {
			return cloudminiErrorEvidenceConfig{}, err
		}
	}
	return config, nil
}
