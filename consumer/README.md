# datetime consumer

An independent module using `ecosystem::datetime = "0.1.0"` through the local verification registry. Its executable resolves an ambiguous New York appointment with explicit fold policies, validates nanosecond Serde round trips, and applies a month-end clamp policy. Its tests also load the pinned `2026c` bundled timezone data through the versioned package API and check an alias across the registry boundary.

The `--oracle <timezone-directory>` mode is used by the library's differential tests. It accepts tab-separated calendar, UTC, local-time, RFC3339, and zone-load requests on standard input. It is test tooling, not an untrusted-input network protocol.

Run `just ecosystem-test datetime` from the repository root to create the versioned registry snapshot and verify both modules.
