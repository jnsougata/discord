package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/textproto"
)

func MultiPartWriter(data map[string]interface{}, files []File) []byte {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	payload, _ := json.MarshalIndent(data, "", "  ")
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="payload_json"`)
	h.Set("Content-Type", `application/json`)
	field, _ := writer.CreatePart(h)
	_, _ = field.Write(payload)
	for i, file := range files {
		ff, _ := writer.CreateFormFile(fmt.Sprintf("file[%v]", i), file.Name)
		_, _ = ff.Write(file.Data)
	}
	_ = writer.Close()
	b := buffer.Bytes()
	fmt.Println(string(b))
	return b
}
