//go:build !integration

package common

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestSuspendKeepsRequiredSnapshotAndSkipsDuplicateMsync(t *testing.T) {
	var suspends atomic.Int32
	var msyncs atomic.Int32
	var after atomic.Int32
	state := NewVMStateMgr(
		context.Background(),
		func(context.Context, time.Duration) error { suspends.Add(1); return nil },
		time.Second,
		func(context.Context) error { msyncs.Add(1); return nil },
		func() error { return nil },
		func() { after.Add(1) },
	)

	if err := state.SuspendAndMsync(); err != nil {
		t.Fatal(err)
	}
	if suspends.Load() != 1 || after.Load() != 1 {
		t.Fatalf("suspend/after calls = %d/%d, want 1/1", suspends.Load(), after.Load())
	}
	if msyncs.Load() != 0 {
		t.Fatalf("duplicate msync calls = %d, want 0", msyncs.Load())
	}
}

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
