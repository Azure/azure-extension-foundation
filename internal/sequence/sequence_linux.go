// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package sequence

import (
	"fmt"
	"github.com/Azure/azure-extension-foundation/internal/settings"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var mostRecentSequenceFileName = "mrseq"

const chmod = os.FileMode(0600)

// GetEnvironmentMostRecentSequenceNumber returns the environment most recent sequence number
func GetEnvironmentMostRecentSequenceNumber() (int, error) {
	hEnv, err := settings.GetEnvironment()
	if err != nil {
		return -1, errorhelper.AddStackToError(fmt.Errorf("unable to parse handler environment : %s", err))
	}
	return findEnvironmentMostRecentSequenceNumber(hEnv.HandlerEnvironment.ConfigFolder)
}

// GetExtensionMostRecentSequenceNumber returns the extension most recent sequence number
func GetExtensionSequenceNumber() (int, error) {
	return findExtensionMostRecentSequenceNumber()
}

// SetExtensionMostRecentSequenceNumber sets the extension most recent sequence number by writing the sequence
// number to the respective extension "mrseq" file
func SetExtensionMostRecentSequenceNumber(sequenceNumber int) error {
	return setExtensionMostRecentSequenceNumber(sequenceNumber)
}

// findEnvironmentMostRecentSequenceNumber finds the most recent environment mrseq by looking up at the
// highest *.settings file in the handler config folder
func findEnvironmentMostRecentSequenceNumber(configFolder string) (int, error) {
	g, err := filepath.Glob(configFolder + "/*.settings")
	if err != nil {
		return 0, err
	}

	sequence := make([]int, len(g))
	for _, v := range g {
		f := filepath.Base(v)
		i, err := strconv.Atoi(strings.Replace(f, ".settings", "", 1))
		if err != nil {
			return 0, errorhelper.AddStackToError(fmt.Errorf("can't parse int from filename: %s", f))
		}
		sequence = append(sequence, i)
	}

	if len(sequence) == 0 {
		return 0, errorhelper.AddStackToError(fmt.Errorf("can't find out seqnum from %s, not enough files", configFolder))
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sequence)))
	return sequence[0], nil
}

// findExtensionMostRecentSequenceNumber find the most recent extension mrseq by reading the extension "mrseq" file
func findExtensionMostRecentSequenceNumber() (int, error) {
	mrseqStr, err := ioutil.ReadFile(mostRecentSequenceFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return -1, nil
		}
		return -1, errorhelper.AddStackToError(fmt.Errorf("failed to read mrseq file : %s", err))
	}

	mrseq, err := strconv.Atoi(string(mrseqStr))
	return mrseq, nil
}

// setExtensionMostRecentSequenceNumber sets the extension mrseq by writing the current mrseq in the extension
// "mrseq" file
func setExtensionMostRecentSequenceNumber(sequenceNumber int) error {
	b := []byte(fmt.Sprintf("%v", sequenceNumber))
	err := ioutil.WriteFile(mostRecentSequenceFileName, b, chmod)
	if err != nil {
		return err
	}
	return nil
}
