package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"gopkg.in/yaml.v3"
)

var validate = validator.New()

func LoadConfig(path string) (*Config, error) {
	buf, err := os.ReadFile(path)
	if err != nil {
		return nil, errors.Join(
			errors.New("failed to open the config file"),
			err,
		)
	}

	cfg := Config{}

	ext := GetExtension(path)

	switch ext {
	case ".json":
		err = json.Unmarshal(buf, &cfg)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(buf, &cfg)
	default:
		return nil, errors.New(
			"the config file extension must be 'json', 'yml' or 'yaml'",
		)
	}

	if err != nil {
		return nil, errors.Join(
			errors.New("failed to decode config file"),
			err,
		)
	}

	if err = validate.Struct(&cfg); err != nil {
		return nil, errors.Join(errors.New("invalid config schema"), err)
	}

	return &cfg, nil
}

func ValidatePayload(secret string, payload []byte, signature string) bool {
	hm := hmac.New(sha1.New, []byte(secret))
	hm.Write(payload)

	expectedHash := hex.EncodeToString(hm.Sum(nil))

	return hmac.Equal([]byte(signature[5:]), []byte(expectedHash))
}

func GetExtension(path string) string {
	s := strings.Split(path, ".")
	if len(s) > 0 {
		return "." + s[len(s)-1]
	}

	return ""
}
