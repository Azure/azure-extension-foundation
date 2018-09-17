package status

import (
	"azure-extension-foundation/internal/settings"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type statusReport []statusItem

type statusItem struct {
	Version      float64 `json:"version"`
	TimestampUTC string  `json:"timestampUTC"`
	Status       status  `json:"status"`
}

type status struct {
	Operation        string           `json:"operation"`
	Status           string           `json:"status"`
	FormattedMessage formattedMessage `json:"formattedMessage"`
}
type formattedMessage struct {
	Lang    string `json:"lang"`
	Message string `json:"message"`
}

// ReportStatus saves operation status to the status file for the extension
// handler with the optional given message, if the given cmd requires reporting
// status.
//
// If an error occurs reporting the status, it will be logged and returned.
func ReportStatus(sequenceNumber int, t string, operation, message string) error {
	s := newStatus(t, operation, message)
	hEnv, err := settings.GetEnvironment()
	if err != nil {
		return fmt.Errorf("unable to get handler environment settings : %v", err)
	}

	if err := s.Save(hEnv.HandlerEnvironment.StatusFolder, sequenceNumber); err != nil {
		//ctx.Log("event", "failed to save handler status", "error", err)
		return fmt.Errorf("failed to save handler operation status : %s", err)
	}
	return nil
}

// Save persists the status message to the specified status folder using the
// sequence number. The operation consists of writing to a temporary file in the
// same folder and moving it to the final destination for atomicity.
func (r statusReport) Save(statusFolder string, seqNum int) error {
	fn := fmt.Sprintf("%d.status", seqNum)
	path := filepath.Join(statusFolder, fn)
	tmpFile, err := ioutil.TempFile(statusFolder, fn)
	if err != nil {
		return fmt.Errorf("status: failed to create temporary file: %v", err)
	}
	tmpFile.Close()

	b, err := r.marshal()
	if err != nil {
		return fmt.Errorf("status: failed to marshal into json: %v", err)
	}
	if err := ioutil.WriteFile(tmpFile.Name(), b, 0644); err != nil {
		return fmt.Errorf("status: failed to path=%s error=%v", tmpFile.Name(), err)
	}

	if err := os.Rename(tmpFile.Name(), path); err != nil {
		return fmt.Errorf("status: failed to move to path=%s error=%v", path, err)
	}
	return nil
}

func newStatus(t string, operation, message string) statusReport {
	return []statusItem{
		{
			Version:      1.0,
			TimestampUTC: time.Now().UTC().Format(time.RFC3339),
			Status: status{
				Operation: operation,
				Status:    t,
				FormattedMessage: formattedMessage{
					Lang:    "en",
					Message: message},
			},
		},
	}
}

func (r statusReport) marshal() ([]byte, error) {
	return json.MarshalIndent(r, "", "\t")
}
