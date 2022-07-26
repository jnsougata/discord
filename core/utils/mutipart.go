package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/textproto"
	"strings"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func MultiPartWriter(data map[string]interface{}, files []File) ([]byte, string) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	payload, _ := json.MarshalIndent(data, "", "  ")
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="payload_json"`)
	h.Set("Content-Type", `application/json`)
	field, _ := writer.CreatePart(h)
	_, _ = field.Write(payload)
	for i, file := range files {
		ff, _ := writer.CreateFormFile(fmt.Sprintf("file[%v]", i), escapeQuotes(file.Name))
		_, _ = ff.Write(file.Data)
	}
	_ = writer.Close()
	return buffer.Bytes(), writer.Boundary()
}
