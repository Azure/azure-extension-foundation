package settings

import "azure-extension-foundation/internal/settings"

// GetExtensionSettings reads the settings for the provided sequenceNumber and assigns the settings to the
// respective structure reference
func GetExtensionSettings(sequenceNumber int, publicSettings, protectedSettings interface{}) error {
	return settings.GetExtensionSettings(sequenceNumber, publicSettings, protectedSettings)
}
