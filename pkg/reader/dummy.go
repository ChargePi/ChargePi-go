package reader

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/xBlaz3kx/ChargePi-go/pkg/models/settings"
)

type Dummy struct {
	logger   *log.Logger
	listener chan string
	settings settings.DummyReader
	tagIndex int
}

func NewDummy(settings *settings.DummyReader) (*Dummy, error) {
	return &Dummy{
		settings: *settings,
		listener: make(chan string),
		tagIndex: 0,
	}, nil
}

func (d *Dummy) ListenForTags(ctx context.Context) {
	d.logger.Info("Listening for tags")

	tag := d.settings.TagIds[d.tagIndex]

	d.listener <- tag
	d.logger.Infof("Tag %s read", tag)
	d.tagIndex++
}

func (d *Dummy) Cleanup() {
	d.logger.Info("Cleaning up reader")
}

func (d *Dummy) Reset() {
	d.logger.Info("Resetting reader")
}

func (d *Dummy) GetTagChannel() <-chan string {
	return d.listener
}

func (d *Dummy) GetType() string {
	return TypeDummy
}
