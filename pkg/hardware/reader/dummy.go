package reader

import (
	"context"
	"go.uber.org/zap"
)

const TypeDummy = "dummy"

type DummySettings struct {
	TagIds []string `json:"tagIds,omitempty" yaml:"tagIds,omitempty" mapstructure:"tagIds,omitempty"`
}

type Dummy struct {
	logger   *zap.Logger
	listener chan string
	settings DummySettings
	tagIndex int
}

func NewDummy(logger *zap.Logger, settings *DummySettings) (*Dummy, error) {
	return &Dummy{
		settings: *settings,
		listener: make(chan string),
		logger: logger,
		tagIndex: 0,
	}, nil
}

func (d *Dummy) ListenForTags(ctx context.Context) {
	d.logger.Info("Listening for tags")

	tag := d.settings.TagIds[d.tagIndex]

	d.listener <- tag
	d.logger.Sugar().Infof("Tag %s read", tag)
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
