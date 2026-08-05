package main

import (
	"testing"

	"github.com/loopholelabs/silo/pkg/storage/protocol/packets"
)

func TestMigrationOptions(t *testing.T) {
	const concurrency = 17

	options := migrationOptions(concurrency)
	if options.Concurrency != concurrency {
		t.Fatalf("concurrency = %d, want %d", options.Concurrency, concurrency)
	}
	if !options.Compression {
		t.Fatal("compression is disabled")
	}
	if options.CompressionType == 0 {
		t.Fatal("compression type is the invalid zero value")
	}
	if options.CompressionType != packets.CompressionTypeZeroes {
		t.Fatalf("compression type = %d, want %d", options.CompressionType, packets.CompressionTypeZeroes)
	}

	if _, err := packets.EncodeWriteAtComp(options.CompressionType, 0, make([]byte, 4096)); err != nil {
		t.Fatalf("compression configuration rejected: %v", err)
	}
}
