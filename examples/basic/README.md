# datetime example

An example using the library in this repository. Its executable resolves an ambiguous New York appointment with explicit fold policies, validates nanosecond Serde round trips, and applies a month-end clamp policy. Its tests also load the pinned `2026c` bundled timezone data through the public API and check an alias. The ecosystem verifier repeats these tests using an isolated registry snapshot.

The `--oracle <timezone-directory>` mode is used by the library's differential tests. It accepts tab-separated calendar, UTC, local-time, RFC3339, and zone-load requests on standard input. It is test tooling, not an untrusted-input network protocol.

Run `(cd ../../../workflows && just ecosystem-test datetime)` from this example directory to create the registry snapshot and verify the library, example, and cached builds.

This example shares the library root manifest and development dependencies. Run `goml test` to build and test the library and its examples.
