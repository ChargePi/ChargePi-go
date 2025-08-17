package diagnostics

import "context"

type httpUploader struct {
}

func newHttpUploader() *httpUploader {
	return &httpUploader{}
}

func (*httpUploader) Upload(context.Context, []byte, string) error {
	// todo implement
	return nil
}
