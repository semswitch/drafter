//go:build !integration

package common

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSuspendAndMsyncStopsWhenBeforeSuspendFails(t *testing.T) {
	holdErr := errors.New("route hold failed")
	suspendCalled := false
	msyncCalled := false
	afterSuspendCalled := false
	state := NewVMStateMgr(
		context.Background(),
		func(context.Context, time.Duration) error {
			suspendCalled = true
			return nil
		},
		time.Second,
		func(context.Context) error {
			msyncCalled = true
			return nil
		},
		func() error { return holdErr },
		func() { afterSuspendCalled = true },
	)

	if err := state.SuspendAndMsync(); !errors.Is(err, holdErr) {
		t.Fatalf("SuspendAndMsync error = %v, want %v", err, holdErr)
	}
	if suspendCalled {
		t.Fatal("suspend function called after OnBeforeSuspend failure")
	}
	if msyncCalled {
		t.Fatal("msync called after OnBeforeSuspend failure")
	}
	if afterSuspendCalled {
		t.Fatal("OnAfterSuspend called after OnBeforeSuspend failure")
	}
}
