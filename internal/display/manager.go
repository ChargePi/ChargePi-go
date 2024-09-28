package display

import (
	"context"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/ChargePi/ChargePi-go/pkg/hardware/display"
	"github.com/ChargePi/ChargePi-go/pkg/util"
	"github.com/go-co-op/gocron"
	message "github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
	log "github.com/sirupsen/logrus"
)

type Manager interface {
	DisplayMessage(message message.MessageInfo) error
	RemoveMessage(messageId string) error
	GetCurrentMessage() (*message.MessageInfo, error)
	SetStrategy(strategy Strategy) error
	// GetMessagesInQueue() ([]display.GetDisplayMessagesResponse, error)
	Cleanup(ctx context.Context) error
}

type DisplayManager struct {
	mu              sync.Mutex
	scheduler       *gocron.Scheduler
	logger          log.FieldLogger
	displayStrategy Strategy
	displays        map[string]display.Display
	//messageQueuePerDisplay map[string]*messagePriorityQueue

	// Message priority queue per display???
	messageQueue *messagePriorityQueue
}

func NewDisplayManager() (*DisplayManager, error) {
	logger := log.StandardLogger()
	scheduler := gocron.NewScheduler(time.UTC)

	return &DisplayManager{
		scheduler:    scheduler,
		logger:       logger,
		messageQueue: newMessageQueue(),
		displays:     make(map[string]display.Display),
	}, nil
}

func (d *DisplayManager) DisplayMessage(message message.MessageInfo) error {
	// Schedule the display of the message a bit later than the start time to prevent the message from being displayed
	if message.StartDateTime != nil && message.StartDateTime.After(time.Now()) {
		_, err := d.scheduler.At(*message.StartDateTime).Tag("displayMessage", strconv.Itoa(message.ID)).Do(d.DisplayMessage, message)
		if err != nil {
			d.logger.WithError(err).Errorf("Error scheduling ClearMessage")
			return err
		}

		return nil
	}

	// Trigger display strategy
	d.logger.Debugf("Displaying message %v", message)
	err := d.displayStrategy.DisplayMessage(d, message)
	if err != nil {
		return err
	}

	// Schedule the clearing of the message
	if message.EndDateTime != nil {
		_, err := d.scheduler.At(*message.EndDateTime).Tag("clearMessage", strconv.Itoa(message.ID)).Do(d.RemoveMessage, message.ID)
		if err != nil {
			d.logger.WithError(err).Errorf("Error scheduling ClearMessage")
			return err
		}
	}

	return nil
}

func (d *DisplayManager) SetStrategy(strategy Strategy) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if util.IsNilInterfaceOrPointer(strategy) {
		return errors.New("strategy cannot be nil")
	}

	d.displayStrategy = strategy
	return nil
}

// RemoveMessage removes a message from the display
func (d *DisplayManager) RemoveMessage(messageId string) error {
	// Remove the message
	d.mu.Lock()
	defer d.mu.Unlock()
	// todo
	return nil
}

// GetCurrentMessage returns the current message being displayed
func (d *DisplayManager) GetCurrentMessage() (*message.MessageInfo, error) {
	// todo
	return nil, nil
}

// Cleanup cleans up the display manager
func (d *DisplayManager) Cleanup(ctx context.Context) error {
	d.logger.Info("Cleaning up display manager")
	d.logger.Debug("Removing jobs from scheduler")
	err := d.clearScheduler()

	errChan := make(chan error, len(d.displays))

	d.logger.Debug("Cleaning up displays")

	err2 := d.cleanupDisplays(ctx, err, errChan)
	if err2 != nil {
		return err2
	}

	return nil
}

// cleanupDisplays cleans up all displays
func (d *DisplayManager) cleanupDisplays(ctx context.Context, err error, errChan chan error) error {
	for _, dis := range d.displays {
		err = dis.Cleanup(ctx)
		if err != nil {
			d.logger.WithError(err).Error("Error cleaning up display")
			errChan <- err
		}
	}

	close(errChan)

	// Wait for all displays to be cleaned up
	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}

// clearScheduler removes all jobs from the scheduler
func (d *DisplayManager) clearScheduler() error {
	jobs, err := d.scheduler.FindJobsByTag("displayMessage")
	if err == nil {
		for _, job := range jobs {
			err := d.scheduler.RemoveByID(job)
			if err != nil {
				d.logger.Error("Error removing job while cleaning up")
			}
		}
	}

	jobs, err = d.scheduler.FindJobsByTag("clearMessage")
	if err == nil {
		for _, job := range jobs {
			err := d.scheduler.RemoveByID(job)
			if err != nil {
				d.logger.Error("Error removing job while cleaning up")
			}
		}
	}

	d.scheduler.Stop()
	return err
}
