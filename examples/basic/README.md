# datetime example

An example using the library in this repository. Its executable resolves an ambiguous New York appointment with explicit fold policies, validates nanosecond Serde round trips, and applies a month-end clamp policy. Its tests also load the pinned `2026c` bundled timezone data through the public API and check an alias. Independent downstream verification repeats these tests across the registry boundary.

The `--oracle <timezone-directory>` mode is used by the library's differential tests. It accepts tab-separated calendar, UTC, local-time, RFC3339, and zone-load requests on standard input. It is test tooling, not an untrusted-input network protocol.

Run `(cd ../../../verification && just ecosystem-test datetime)` from this example directory to create the versioned registry snapshot and verify the library, example, and independent downstream snapshot.

This example shares the library root manifest and development dependencies. Run `goml verify --example basic` to build and test it as an independent downstream module.
