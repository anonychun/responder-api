package message

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if value != "" {
			return value
		}
	}

	return ""
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}
