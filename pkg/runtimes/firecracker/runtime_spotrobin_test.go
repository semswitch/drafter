package firecracker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/loopholelabs/drafter/pkg/ipc"
)

func TestSuspendUsesConfiguredLifecycleWithoutAgentRPC(t *testing.T) {
	want := errors.New("guest lifecycle failed")
	provider := &FirecrackerRuntimeProvider[
		struct{}, ipc.AgentServerRemote[struct{}], struct{},
	]{
		SkipAgentRPC: true,
		BeforeSuspend: func(context.Context) error {
			return want
		},
	}
	provider.running = true

	err := provider.Suspend(context.Background(), time.Second, nil)
	if !errors.Is(err, want) {
		t.Fatalf("Suspend error = %v, want %v", err, want)
	}
}
