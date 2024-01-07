package main

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
)

// Read the request body and replace the reader so that subsequent ctx.Bind
// calls succeed.
func peekBody(ctx *gin.Context) ([]byte, error) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read request body: %w", err)
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	return body, nil
}
