package i18n

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type translatorTestSuite struct {
	suite.Suite
	translator *TranslatorImpl
}

func (s *translatorTestSuite) SetupSuite() {
	settings := Settings{
		SupportedLanguages: []string{"en", "sl"},
	}

	translator, err := NewTranslator(settings)
	s.Require().NoError(err)

	s.translator = translator
}

func (s *translatorTestSuite) TestTranslate() {
	tests := []struct {
		name                string
		messageId           string
		language            string
		args                map[string]interface{}
		expectedTranslation string
		wantErr             bool
	}{
		{
			name:      "ConnectorFinishing in english",
			messageId: "ConnectorFinishing",
			language:  "en",
			args: map[string]interface{}{
				"Id": 1,
			},
			expectedTranslation: "Finishing charging on connector 1.",
			wantErr:             false,
		},
		{
			name:      "ConnectorAvailable in English",
			messageId: "ConnectorAvailable",
			language:  "en",
			args: map[string]interface{}{
				"Id": 2,
			},
			expectedTranslation: "Connector 2 is available. Plug in your vehicle to start charging!",
			wantErr:             false,
		},
		{
			name:      "WelcomeMessage in Slovene",
			messageId: "WelcomeMessage",
			language:  "sl",
			args: map[string]interface{}{
				"StationName": "Polje",
			},
			expectedTranslation: "Pozdravljeni na polnilni postaji Polje. Priključite svoje vozilo in začnite polniti!",
			wantErr:             false,
		},
		{
			name:      "ConnectorCharging in English",
			messageId: "ConnectorCharging",
			language:  "en",
			args: map[string]interface{}{
				"Id": "3",
			},
			expectedTranslation: "Started charging on connector 3.",
			wantErr:             false,
		},
		{
			name:      "Slovak is unsupported",
			messageId: "ConnectorFaulted",
			language:  "sk", // Slovak not supported
			args: map[string]interface{}{
				"Id": 2,
			},
			expectedTranslation: "",
			wantErr:             true,
		},
		{
			name:                "AuthenticateWithRFID in Slovene",
			messageId:           "AuthenticateWithRFID",
			language:            "sl",
			args:                map[string]interface{}{},
			expectedTranslation: "Za nadaljevanje postopka se identificirajte s svojo RFID kartico.",
			wantErr:             false,
		},
		{
			name:      "AuthenticateWithApp in English",
			messageId: "AuthenticateWithApp",
			language:  "en",
			args: map[string]interface{}{
				"AppName": "EVCC",
			},
			expectedTranslation: "To start charging, open the EVCC app and scan the QR code.",
			wantErr:             false,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			translate, err := s.translator.Translate(tt.messageId, tt.language, tt.args)
			if tt.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
			s.Equal(tt.expectedTranslation, translate)
		})
	}
}

func (s *translatorTestSuite) Test_localize() {
	tests := []struct {
		name                string
		messageId           string
		language            string
		args                map[string]interface{}
		expectedTranslation string
		wantErr             bool
	}{
		{
			name:      "Translate in English",
			messageId: "ConnectorFinishing",
			language:  "en",
			args: map[string]interface{}{
				"Id": 1,
			},
			expectedTranslation: "Finishing charging on connector 1.",
			wantErr:             false,
		},
		{
			name:      "Translate in Slovene",
			messageId: "ConnectorFinishing",
			language:  "sl",
			args: map[string]interface{}{
				"Id": 1,
			},
			expectedTranslation: "Finishing charging on connector 1.",
			wantErr:             false,
		},
		{
			name:      "Unsupported language",
			messageId: "ConnectorFinishing",
			language:  "sk",
			args: map[string]interface{}{
				"Id": 1,
			},
			expectedTranslation: "",
			wantErr:             true,
		},
		{
			name:      "Empty key",
			messageId: "",
			language:  "en",
			args: map[string]interface{}{
				"Id": 1,
			},
			expectedTranslation: "",
			wantErr:             true,
		},
		{
			name:      "Empty language",
			messageId: "ConnectorFinishing",
			language:  "",
			args: map[string]interface{}{
				"Id": 1,
			},
			expectedTranslation: "",
			wantErr:             true,
		},
		{
			name:      "Unknown key",
			messageId: "UnknownKey",
			language:  "en",
			args: map[string]interface{}{
				"Id": 1,
			},
			expectedTranslation: "",
			wantErr:             true,
		},
		{
			name:                "Invalid args",
			messageId:           "ConnectorFinishing",
			language:            "en",
			args:                map[string]interface{}{},
			expectedTranslation: "",
			wantErr:             true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			translate, err := s.translator.localize(tt.messageId, tt.language, tt.args)
			if tt.wantErr {
				s.Error(err)
				return
			}

			s.NoError(err)
			s.Equal(tt.expectedTranslation, translate)
		})
	}

}

func TestTranslator(t *testing.T) {
	suite.Run(t, new(translatorTestSuite))
}
