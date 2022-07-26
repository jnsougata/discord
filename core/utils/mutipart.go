package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jnsougata/disgo/core/file"
	"log"
	"mime/multipart"
	"net/textproto"
	"strings"
)

var quoteEscaper = strings.NewReplacer("\\", "\\\\", `"`, "\\\"")

func escapeQuotes(s string) string {
	return quoteEscaper.Replace(s)
}

func MultiPartWriter(data map[string]interface{}, files []file.File) ([]byte, string) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	payload, _ := json.MarshalIndent(data, "", "  ")
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="payload_json"`)
	h.Set("Content-Type", `application/json`)
	field, _ := writer.CreatePart(h)
	_, _ = field.Write(payload)
	for i, f := range files {
		if f.Content != nil {
			ff, _ := writer.CreateFormFile(fmt.Sprintf(`files[%v]`, i), escapeQuotes(f.Name))
			_, _ = ff.Write(f.Content)
		} else {
			log.Println(fmt.Sprintf(`File content for %v is empty`, f.Name))
		}
	}
	_ = writer.Close()
	return buffer.Bytes(), writer.Boundary()
}
