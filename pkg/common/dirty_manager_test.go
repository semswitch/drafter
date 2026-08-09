//go:build !integration
// +build !integration

package common

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDirtyManager(t *testing.T) {

	devName := "test"

	devices := map[string]*DeviceStatus{
		devName: {
			MinCycles:      5,
			MaxCycles:      10,
			CycleThrottle:  100 * time.Millisecond,
			MaxDirtyBlocks: 10,
		},
	}

	var suspendCalled sync.WaitGroup
	var authTransferCalled sync.WaitGroup

	authTransferCalled.Add(1)
	authTransfer := func() error {
		authTransferCalled.Done()
		return nil
	}

	suspendCalled.Add(1)
	suspendFunc := func(ctx context.Context, timeout time.Duration) error {
		suspendCalled.Done()
		return nil
	}

	suspendTimeout := 100 * time.Millisecond

	msyncFunc := func(context.Context) error { return nil }

	onBeforeSuspend := func() error { return nil }

	onAfterSuspend := func() {}

	vm := NewVMStateMgr(context.TODO(), suspendFunc, suspendTimeout, msyncFunc, onBeforeSuspend, onAfterSuspend)

	dm := NewDirtyManager(vm, devices, authTransfer)

	// Try something simple...

	noMoreDirtyBlocks := 8

	count := 0
	for {
		blocks := []uint{0} // Dummy, under the threshold.
		if count >= noMoreDirtyBlocks {
			blocks = []uint{} // No dirty blocks...
		}

		_, err := dm.PreGetDirty(devName)
		assert.NoError(t, err)
		more, err := dm.PostGetDirty(devName, blocks)
		assert.NoError(t, err)

		if count >= noMoreDirtyBlocks {
			assert.False(t, more)
			break
		}

		assert.True(t, more)
		more, err = dm.PostMigrateDirty(devName, blocks)
		assert.NoError(t, err)
		if count > 5 {
			assert.False(t, more)
		} else {
			assert.True(t, more)
		}
		count++

		// It should have run some things since it's past min
		if count > 5 {
			suspendCalled.Wait()
		}

		// Auth is tranfered on NEXT loop to allow dirtyList to be sent
		if count > 6 {
			// Make sure auth was transferred
			authTransferCalled.Wait()
		}
	}

}

func TestDirtyManagerIdleWaitWakesOnSuspension(t *testing.T) {
	state := NewVMStateMgr(
		context.Background(),
		func(context.Context, time.Duration) error { return nil },
		time.Second,
		func(context.Context) error { return nil },
		func() error { return nil },
		func() {},
	)
	manager := NewDirtyManager(state, map[string]*DeviceStatus{}, func() error { return nil })
	returned := make(chan error, 1)
	go func() { returned <- manager.WaitWhenIdle(DeviceMemoryName) }()

	select {
	case <-returned:
		t.Fatal("idle wait returned before suspension")
	case <-time.After(25 * time.Millisecond):
	}
	if err := state.SuspendAndMsync(); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-returned:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("idle wait did not wake after suspension")
	}
}

func TestDirtyManagerStopsCyclingWhenDirtySetDoesNotConverge(t *testing.T) {
	device := &DeviceStatus{
		MinCycles: 1, MaxCycles: 3, MaxDirtyBlocks: 1, StopOnNonConvergence: true,
	}
	suspended := false
	state := NewVMStateMgr(
		context.Background(),
		func(context.Context, time.Duration) error { suspended = true; return nil },
		time.Second,
		func(context.Context) error { return nil },
		func() error { return nil },
		func() {},
	)
	manager := NewDirtyManager(state, map[string]*DeviceStatus{DeviceMemoryName: device}, func() error { return nil })

	for _, count := range []int{4, 4} {
		blocks := make([]uint, count)
		if _, err := manager.PreGetDirty(DeviceMemoryName); err != nil {
			t.Fatal(err)
		}
		if _, err := manager.PostGetDirty(DeviceMemoryName, blocks); err != nil {
			t.Fatal(err)
		}
		if _, err := manager.PostMigrateDirty(DeviceMemoryName, blocks); err != nil {
			t.Fatal(err)
		}
	}

	if !device.Ready || !suspended {
		t.Fatalf("non-converging device ready=%t suspended=%t", device.Ready, suspended)
	}
}
