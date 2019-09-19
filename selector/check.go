package selector

// CheckLabels checks a label set with the default validation rules.
func CheckLabels(labels Labels) error {
	var err error
	for key, value := range labels {
		if err = CheckKey(key); err != nil {
			return err
		}
		if err = CheckValue(value); err != nil {
			return err
		}
	}
	return nil
}

// CheckKey validates a key with the default validation options.
func CheckKey(key string) error {
	return new(Parser).CheckKey(key)
}

// CheckValue returns if the value is valid with the default validation options.
func CheckValue(value string) error {
	return new(Parser).CheckValue(value)
}

// CheckDNS returns if a given string is a valid dns name with a given set of options.
func CheckDNS(dnsName string) error {
	return new(Parser).CheckDNS(dnsName)
}
