package diagnostics

import (
	"bytes"
	"context"
	"net/url"
	"time"

	"github.com/jlaffaye/ftp"
)

type ftpUploader struct {
	settings *FtpUploaderSettings
}

type FtpUploaderSettings struct {
	Timeout *time.Duration
}

func newFtpUploader(settings FtpUploaderSettings) *ftpUploader {
	if settings.Timeout == nil {
		defaultTimeout := 10 * time.Second
		settings.Timeout = &defaultTimeout
	}

	return &ftpUploader{
		settings: &settings,
	}
}

func (f *ftpUploader) Upload(ctx context.Context, stream []byte, location string) error {
	parse, err := url.Parse(location)
	if err != nil {
		return err
	}

	// Set a default timeout if none is provided.
	if f.settings.Timeout == nil {
		defaultTimeout := 10 * time.Second
		f.settings.Timeout = &defaultTimeout
	}

	// Connect to the FTP server
	connection, err := ftp.Dial(
		parse.Hostname(),
		ftp.DialWithContext(ctx),
		ftp.DialWithTimeout(*f.settings.Timeout),
	)
	if err != nil {
		return err
	}

	pass, hasPassword := parse.User.Password()
	if !hasPassword {
		pass = ""
	}

	// Login to the FTP server.
	err = connection.Login(parse.User.Username(), pass)
	if err != nil {
		return err
	}

	// Upload the local file to the remote FTP server.
	return connection.Stor(parse.Path, bytes.NewReader(stream))
}
