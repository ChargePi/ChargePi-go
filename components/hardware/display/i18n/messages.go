package i18n

func TranslateConnectorAvailableMessage(lang string, connectorId int) ([]string, error) {
	data := make(map[string]interface{})
	data["Id"] = connectorId

	firstPart, err := Localize(lang, "ConnectorTemplate", data, nil)
	secondPart, err := Localize(lang, "ConnectorAvailable", nil, nil)
	if err != nil {
		return nil, err
	}

	return []string{firstPart, secondPart}, nil
}

func TranslateConnectorFinishingMessage(lang string, connectorId int) ([]string, error) {
	data := make(map[string]interface{})
	data["Id"] = connectorId

	firstPart, err := Localize(lang, "ConnectorFinishing", nil, nil)
	secondPart, err := Localize(lang, "ConnectorStopTemplate", data, nil)
	if err != nil {
		return nil, err
	}

	return []string{firstPart, secondPart}, nil
}

func TranslateConnectorFaultedMessage(lang string, connectorId int) ([]string, error) {
	data := make(map[string]interface{})
	data["Id"] = connectorId

	firstPart, err := Localize(lang, "ConnectorTemplate", data, nil)
	secondPart, err := Localize(lang, "ConnectorFaulted", nil, nil)
	if err != nil {
		return nil, err
	}

	return []string{firstPart, secondPart}, nil
}

func TranslateConnectorChargingMessage(lang string, connectorId int) ([]string, error) {
	data := make(map[string]interface{})
	data["Id"] = connectorId

	firstPart, err := Localize(lang, "ConnectorCharging", nil, nil)
	secondPart, err := Localize(lang, "ConnectorStopTemplate", data, nil)
	if err != nil {
		return nil, err
	}

	return []string{firstPart, secondPart}, nil
}

func TranslateWelcomeMessage(lang string) ([]string, error) {
	firstPart, err := Localize(lang, "WelcomeMessage", nil, nil)
	secondPart, err := Localize(lang, "WelcomeMessage2", nil, nil)
	if err != nil {
		return nil, err
	}

	return []string{firstPart, secondPart}, nil
}
