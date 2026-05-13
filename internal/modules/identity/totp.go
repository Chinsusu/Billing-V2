package identity

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	ErrTwoFactorCodeMissing = errors.New("two factor code missing")
	ErrTwoFactorCodeInvalid = errors.New("two factor code invalid")
)

const (
	totpDigits       = 6
	totpPeriod       = 30 * time.Second
	totpSecretLength = 20
)

func GenerateTOTPSecret() (string, error) {
	raw := make([]byte, totpSecretLength)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate totp secret: %w", err)
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(raw), nil
}

func VerifyTOTPCode(secret string, code string, now time.Time) (bool, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return false, ErrTwoFactorCodeMissing
	}
	if len(code) != totpDigits {
		return false, ErrTwoFactorCodeInvalid
	}
	for _, current := range code {
		if current < '0' || current > '9' {
			return false, ErrTwoFactorCodeInvalid
		}
	}
	for offset := int64(-1); offset <= 1; offset++ {
		expected, err := TOTPCodeAt(secret, now.Add(time.Duration(offset)*totpPeriod))
		if err != nil {
			return false, err
		}
		if hmac.Equal([]byte(expected), []byte(code)) {
			return true, nil
		}
	}
	return false, nil
}

func TOTPCodeAt(secret string, at time.Time) (string, error) {
	secret = strings.ToUpper(strings.TrimSpace(secret))
	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		return "", ErrTwoFactorCodeInvalid
	}
	counter := uint64(at.Unix() / int64(totpPeriod.Seconds()))
	var message [8]byte
	binary.BigEndian.PutUint64(message[:], counter)

	mac := hmac.New(sha1.New, key)
	_, _ = mac.Write(message[:])
	sum := mac.Sum(nil)
	offset := sum[len(sum)-1] & 0x0f
	binaryCode := (uint32(sum[offset])&0x7f)<<24 |
		(uint32(sum[offset+1])&0xff)<<16 |
		(uint32(sum[offset+2])&0xff)<<8 |
		(uint32(sum[offset+3]) & 0xff)
	value := binaryCode % uint32(math.Pow10(totpDigits))
	return leftPadTOTP(strconv.FormatUint(uint64(value), 10)), nil
}

func leftPadTOTP(value string) string {
	for len(value) < totpDigits {
		value = "0" + value
	}
	return value
}

func totpProvisionURI(email string, secret string) string {
	label := url.QueryEscape("Billing:" + strings.TrimSpace(email))
	query := url.Values{}
	query.Set("secret", secret)
	query.Set("issuer", "Billing")
	query.Set("algorithm", "SHA1")
	query.Set("digits", strconv.Itoa(totpDigits))
	query.Set("period", strconv.Itoa(int(totpPeriod.Seconds())))
	return "otpauth://totp/" + label + "?" + query.Encode()
}
