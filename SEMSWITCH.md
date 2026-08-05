# SemSwitch Drafter

This is an unofficial SemSwitch-maintained fork of Loophole Labs Drafter. The initial maintenance baseline is Drafter v0.7.4, with the original module path and upstream history preserved.

This baseline pins SemSwitch Silo `v0.2.20-semswitch.1` and corrects `drafter-peer` to use `CompressionTypeZeroes` when compression is enabled. It preserves the exact Firecracker and jailer artifacts paired with Drafter v0.7.4.

The stack has passed Drafter's official repeated-migration integrity test and bidirectional Azure host-to-host Firecracker live migration. SemSwitch intends to keep modifications narrow and compatibility-focused.

## Consumption

Go `replace` directives apply only when their declaring module is the main module. Cloning and building this repository directly therefore uses the SemSwitch Silo replacement, but a downstream module importing Drafter does not inherit it. Downstream consumers must explicitly replace both Drafter and Silo.

After this release is published, the intended downstream configuration is:

```go
require (
	github.com/loopholelabs/drafter v0.7.4
	github.com/loopholelabs/silo v0.2.20
)

replace github.com/loopholelabs/drafter =>
	github.com/semswitch/drafter v0.7.4-semswitch.1

replace github.com/loopholelabs/silo =>
	github.com/semswitch/silo v0.2.20-semswitch.1
```
