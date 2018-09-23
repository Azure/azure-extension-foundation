package sequence

import "github.com/Azure/azure-extension-foundation/internal/sequence"

// GetMostRecentSequenceNumber return the extension and environment most recent sequence number
func GetMostRecentSequenceNumber() (int, int, error) {
	extensionMrseq, err := GetExtensionMostRecentSequenceNumber()
	if err != nil {
		return -1, -1, err
	}

	environmentMrseq, err := GetEnvironmentMostRecentSequenceNumber()
	if err != nil {
		return -1, -1, err
	}

	return extensionMrseq, environmentMrseq, nil
}

// ShouldBeProcessed returns true when the extension most recent sequence number is below the environment most
// recent sequence number else returns false
func ShouldBeProcessed(extensionMrseq int, environmentMrseq int) bool {
	return extensionMrseq < environmentMrseq
}

// GetEnvironmentMostRecentSequenceNumber returns the environment most recent sequence number
func GetEnvironmentMostRecentSequenceNumber() (int, error) {
	return sequence.GetEnvironmentMostRecentSequenceNumber()
}

// GetExtensionMostRecentSequenceNumber returns the extension most recent sequence number
func GetExtensionMostRecentSequenceNumber() (int, error) {
	return sequence.GetExtensionSequenceNumber()
}

// SetExtensionMostRecentSequenceNumber sets the extension most recent sequence number
func SetExtensionMostRecentSequenceNumber(sequenceNumber int) error {
	return sequence.SetExtensionMostRecentSequenceNumber(sequenceNumber)
}
