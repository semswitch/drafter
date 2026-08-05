# SemSwitch Drafter

This is an unofficial SemSwitch-maintained fork of Loophole Labs Drafter. The initial maintenance baseline is Drafter v0.7.4, with the original module path and upstream history preserved.

This baseline pins SemSwitch Silo `v0.2.20-semswitch.1` and corrects `drafter-peer` to use `CompressionTypeZeroes` when compression is enabled. It preserves the exact Firecracker and jailer artifacts paired with Drafter v0.7.4.

The stack has passed Drafter's official repeated-migration integrity test and bidirectional Azure host-to-host Firecracker live migration. SemSwitch intends to keep modifications narrow and compatibility-focused.
