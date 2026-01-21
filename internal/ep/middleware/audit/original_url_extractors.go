package audit

import (
	"bytes"
	"encoding/json"

	"github.com/ElfAstAhe/url-shortener2/internal/ep/dto"
)

type OriginalURLExtractorFunc func(data []byte) (string, error)

func OriginalURLExtractGetRoot(data []byte) (string, error) {
	return string(data), nil
}

func OriginalURLExtractPostRoot(data []byte) (string, error) {
	return string(data), nil
}

func OriginalURLExtractPostApiShorten(data []byte) (string, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	var res = new(dto.ShortenCreateRequest)
	err := dec.Decode(res)
	if err != nil {
		return "", err
	}

	return res.URL, nil
}
