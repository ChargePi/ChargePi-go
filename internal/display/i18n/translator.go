package i18n

import (
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	log "github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

//go:embed translations/active.*.yaml
var locale embed.FS

type Translator interface {
	Translate(key string, lang string, args map[string]interface{}) (string, error)
}

type Settings struct {
	SupportedLanguages []string `json:"supportedLanguages,omitempty" yaml:"supportedLanguages,omitempty" mapstructure:"supportedLanguages,omitempty"`
}

type TranslatorImpl struct {
	bundle             *i18n.Bundle
	defaultMessages    map[string]i18n.Message
	mu                 sync.Mutex
	matcher            language.Matcher
	supportedLanguages []language.Tag
	logger             *log.Logger
}

func NewTranslator(settings Settings) (*TranslatorImpl, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	_, err := bundle.LoadMessageFileFS(locale, "translations/active.en.yaml")
	if err != nil {
		return nil, err
	}

	translator := &TranslatorImpl{
		bundle:          bundle,
		defaultMessages: map[string]i18n.Message{},
		logger:          log.StandardLogger(),
	}

	// Add defaults
	translator.setupDefaults()

	// Validate & load supported languages
	var supportedLanguages []language.Tag
	for _, lang := range settings.SupportedLanguages {
		tag, err := language.Parse(lang)
		if err != nil {
			return nil, fmt.Errorf("invalid language")
		}

		path := filepath.Join("translations", "active."+tag.String()+".yaml")
		err = translator.loadTranslation(path)
		if err != nil {
			return nil, err
		}

		supportedLanguages = append(supportedLanguages, tag)
	}

	translator.matcher = language.NewMatcher(supportedLanguages)

	return translator, nil
}

func (t *TranslatorImpl) setupDefaults() {
	// Add defaults
	t.addDefaultMessage(i18n.Message{
		ID: "ConnectorAvailable",
	})

	t.addDefaultMessage(i18n.Message{
		ID: "ConnectorCharging",
	})

	t.addDefaultMessage(i18n.Message{
		ID: "ConnectorFaulted",
	})

	t.addDefaultMessage(i18n.Message{
		ID: "ConnectorFinishing",
	})

	t.addDefaultMessage(i18n.Message{
		ID: "WelcomeMessage",
	})

	t.addDefaultMessage(i18n.Message{
		ID: "AuthenticateWithRFID",
	})

	t.addDefaultMessage(i18n.Message{
		ID: "AuthenticateWithApp",
	})
}

func (t *TranslatorImpl) addDefaultMessage(message i18n.Message) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.defaultMessages[message.ID] = message
}

func (t *TranslatorImpl) Translate(messageId string, lang string, args map[string]interface{}) (string, error) {
	// todo define plurals based on the messageId
	switch messageId {
	case "test":

	}

	return t.localize(messageId, lang, args)
}

func (t *TranslatorImpl) localize(messageId string, lang string, args map[string]interface{}) (string, error) {
	tag, _ := language.MatchStrings(t.matcher, lang)
	locale := i18n.NewLocalizer(t.bundle, tag.String())

	t.mu.Lock()
	defaultMessage, ok := t.defaultMessages[messageId]
	if !ok {
		return "", errors.New("message not found")
	}
	t.mu.Unlock()

	msg, err := locale.Localize(&i18n.LocalizeConfig{
		DefaultMessage: &defaultMessage,
		TemplateData:   args["data"],
		PluralCount:    args["plural"],
	})
	if err != nil {
		return "", err
	}

	return msg, nil
}

func (t *TranslatorImpl) loadTranslation(path string) error {
	// active.en.yaml -> en
	strs := strings.Split(filepath.Base(path), ".")
	if len(strs) < 2 {
		return errors.New("invalid translation file name")
	}

	// The language is second to last
	lang := strs[len(strs)-2]
	log.Debugf("loading translation: %s", lang)

	// Load the translation file
	_, err := t.bundle.LoadMessageFile(path)
	return err
}
