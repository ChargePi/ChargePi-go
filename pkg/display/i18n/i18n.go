package i18n

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

var (
	bundle             *i18n.Bundle
	defaultMessages    map[string]i18n.Message
	matcher            language.Matcher
	supportedLanguages = []language.Tag{
		language.English, // The first language is used as fallback.
	}
)

func init() {
	once := sync.Once{}
	once.Do(func() {

		defaultMessages = make(map[string]i18n.Message)

		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

		loadTranslations()

		// Default messages
		addDefaultMessage(i18n.Message{
			ID:    "ConnectorTemplate",
			Other: "Connector {{.Id}}",
		})
		addDefaultMessage(i18n.Message{
			ID:    "ConnectorAvailable",
			Other: "available.",
		})
		addDefaultMessage(i18n.Message{
			ID:    "ConnectorFinishing",
			Other: "Stopped charging",
		})
		addDefaultMessage(i18n.Message{
			ID:    "ConnectorCharging",
			Other: "Started charging",
		})
		addDefaultMessage(i18n.Message{
			ID:    "ConnectorStopTemplate",
			Other: "at {{.Id}}.",
		})
		addDefaultMessage(i18n.Message{
			ID:    "ConnectorFaulted",
			Other: "has faulted.",
		})
		addDefaultMessage(i18n.Message{
			ID:    "WelcomeMessage",
			Other: "Welcome to",
		})
		addDefaultMessage(i18n.Message{
			ID:    "WelcomeMessage2",
			Other: "ChargePi!",
		})
	})
}

// addDefaultMessage exposes the API for plugins
func addDefaultMessage(message i18n.Message) {
	defaultMessages[message.ID] = message
}

// loadTranslations loads all available translations from the translations folder into the bundle.
func loadTranslations() {
	log.Debug("Loading translations..")

	err := filepath.Walk("./hardware/display/i18n/translations", func(path string, info os.FileInfo, err error) error {
		// Load all active.*.yaml translations into the bundle
		if info != nil && !info.IsDir() && strings.Contains(info.Name(), "active.") {
			return loadTranslation(path, info)
		}

		return nil
	})

	if err != nil {
		log.Errorf("Error loading translations: %v", err)
	}

	// Create a matcher based on imported translation files.
	matcher = language.NewMatcher(supportedLanguages)
}

func loadTranslation(path string, info os.FileInfo) error {
	// active.en.yaml -> en
	strs := strings.Split(info.Name(), ".")
	if len(strs) < 2 {
		return fmt.Errorf("invalid file name")
	}

	// The language is second to last
	lang := strs[len(strs)-2]
	log.Debugf("loading translation: %s", lang)

	err := addLanguageSupport(lang)
	if err != nil {
		return err
	}

	// Load the translation file
	_, err = bundle.LoadMessageFile(path)
	return err
}

func addLanguageSupport(lang string) error {
	tag, err := language.Parse(lang)
	if err != nil {
		return err
	}

	supportedLanguages = append(supportedLanguages, tag)
	return nil
}

// Localize translates the message based on the language of the chat.
func Localize(lang string, messageId string, data map[string]interface{}, plural interface{}) (string, error) {
	tag, _ := language.MatchStrings(matcher, lang)
	locale := i18n.NewLocalizer(bundle, tag.String())
	defaultMessage, ok := defaultMessages[messageId]
	if !ok {
		return "", fmt.Errorf("default message not found")
	}

	msg, err := locale.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &defaultMessage,
		TemplateData:   data,
		PluralCount:    plural,
	})
	if err != nil {
		return "", err
	}

	return msg, nil
}
