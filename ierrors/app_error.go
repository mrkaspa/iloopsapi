package ierrors

import (
	"bytes"
	"fmt"
)

type AppError struct {
	Code      int               `json:"code"`
	ErrorS    string            `json:"error"`
	MapErrors map[string]string `json:"map_errors"`
}

func (a AppError) Error() string {
	if a.ErrorS != "" {
		return a.ErrorS
	}
	var buffer bytes.Buffer
	buffer.WriteString("Validation errors:\n")
	for key, value := range a.MapErrors {
		buffer.WriteString(fmt.Sprintf("Field validation for '%s' failed on the field '%s'", key, value))
	}
	return buffer.String()
}
