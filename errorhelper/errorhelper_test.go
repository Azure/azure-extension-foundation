// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT license.

package errorhelper

import (
	"fmt"
	"strings"
	"testing"
)

func ErrorWithStack(stackDepth int, err error) error {
	if stackDepth == 0 {
		if err != nil {
			return AddStackToError(err)
		} else {
			return NewErrorWithStack("Reached the bottom of the error call stack")
		}
	} else {
		return ErrorWithStack(stackDepth-1, err)
	}
}

func TestErrorsWrap(t *testing.T) {
	err1 := ErrorWithStack(5, fmt.Errorf("misc error"))
	//err2 := ErrorWithWrap(5, fmt.Errorf("misc error"))
	err2 := ErrorWithStack(5, nil)

	t.Logf("the error stack trace is: %+v", err1)
	//t.Logf("the error stack trace is %s, %+v", err2, err2)
	t.Logf("the error stack trace is: %+v", err2)

	if !strings.Contains(err1.Error(), "ErrorWithStack") {
		t.Fatal("error doesn't contain stack trace")
	}

	if !strings.Contains(err2.Error(), "ErrorWithStack") {
		t.Fatal("error doesn't contain stack trace")
	}

	if !strings.Contains(err1.Error(), "misc error") {
		t.Fatal("error doesn't base error")
	}

	if !strings.Contains(err2.Error(), "Reached the bottom of the error call stack") {
		t.Fatal("error doesn't base error")
	}
}
