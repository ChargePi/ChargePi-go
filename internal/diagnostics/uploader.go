package diagnostics

import "errors"

var ErrUnsupportedUploadProtocol = errors.New("unsupported upload protocol")

type uploaderType string

const (
	uploaderTypeHttp uploaderType = "http"
	uploaderTypeFtp  uploaderType = "ftp"
)

type uploaderFactory struct {
}

func (u *uploaderFactory) UploaderFromType(t uploaderType) (uploader, error) {
	switch t {

	case uploaderTypeHttp:
		return newHttpUploader(), nil
	case uploaderTypeFtp:
		return newFtpUploader(FtpUploaderSettings{}), nil
	}

	return nil, ErrUnsupportedUploadProtocol
}
