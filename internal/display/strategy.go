package display

import (
	"sync"

	"github.com/ChargePi/ChargePi-go/internal/display/i18n"
	"github.com/ChargePi/ChargePi-go/pkg/hardware/display"
	message "github.com/lorenzodonini/ocpp-go/ocpp2.0.1/display"
)

type Strategy interface {
	DisplayMessage(manager *DisplayManager, message message.MessageInfo) error
}

// translateAndDisplayStrategy translates the message content before displaying it.
type translateAndDisplayStrategy struct {
	translator i18n.Translator
}

func newTranslateAndDisplayStrategy(translator i18n.Translator) Strategy {
	return translateAndDisplayStrategy{
		translator: translator,
	}
}

func (t translateAndDisplayStrategy) DisplayMessage(manager *DisplayManager, message message.MessageInfo) error {
	// Translate the message content
	translatedContent, err := t.translator.Translate(message.Message.Language, message.Message.Content, nil)
	if err != nil {
		return err
	}

	// Update the message content with the translated version
	message.Message.Content = translatedContent

	return displayConcurrently(manager.displays, message)
}

// displayOnlyStrategy displays a message without any translation, letting the display handle the message as is.
type displayOnlyStrategy struct {
}

func (d displayOnlyStrategy) DisplayMessage(manager *DisplayManager, message message.MessageInfo) error {
	return displayConcurrently(manager.displays, message)
}

// displayConcurrently displays a message concurrently on all registered displays.
func displayConcurrently(displays map[string]display.Display, message message.MessageInfo) error {
	var wg sync.WaitGroup
	wg.Add(len(displays))

	errChan := make(chan error, len(displays))

	// Display the message concurrently on all displays
	for _, display := range displays {

		go func() {
			err := display.DisplayMessage(message)
			if err != nil {
				errChan <- err
			}
		}()
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
