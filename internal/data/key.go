package data

import (
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidKeyFormat = errors.New("invalid key format")

type Key string

func (k *Key) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidKeyFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")

	if len(parts) != 2 || len(parts[0]) > 2 {
		return ErrInvalidKeyFormat
	}

	tonic, mode := parts[0], parts[1]

	if !strings.ContainsAny(string(tonic[0]), "ABCDEFG") {
		return ErrInvalidKeyFormat
	}

	if len(tonic) == 2 && !strings.ContainsAny(string(tonic[1]), "b#") {
		return ErrInvalidKeyFormat
	}

	validModes := map[string]bool{
		"major": true,
		"minor": true,
	}

	if !validModes[mode] {
		return ErrInvalidKeyFormat
	}

	*k = Key(unquotedJSONValue)

	return nil
}
