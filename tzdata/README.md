# Bundled timezone data

`ecosystem::datetime::tzdata` embeds a pinned set of 598 compiled TZif files and aliases. It uses the parent `datetime` package's TZif parser, transition lookup, and local-time resolver. Applications can load a zone without depending on the machine's `/usr/share/zoneinfo`:

```goml
use ecosystem::datetime as dt;
use ecosystem::datetime::tzdata;

fn new_york() -> Result[dt::TimeZone, dt::Error] {
    tzdata::load_versioned("America/New_York", tzdata::VERSION)
}
```

`VERSION` is `2026c`. `names()` lists the sorted zone names; `contains()` and `checksum()` inspect the manifest; `zone_bytes()` returns verified raw TZif bytes; `load()` parses them; and `load_versioned()` rejects a mismatched expected version. These functions do not consult system timezone files. A selected `TimeZone` is immutable and retains its rules after an application starts. Unknown names, corrupt encoding or checksums, incompatible versions, and unsupported TZif content return `datetime::Error`.

The snapshot was generated from Ubuntu `tzdata` package `2026c-0ubuntu0.24.04.1`, whose `tzdata.zi` identifies IANA release `2026c`. The source `tzdata.zi` SHA-256 was `bce4e2473fda8a22f5be65be2ef01e2e22cba5e721d8fa8d68df3d5f67c23f13`. The generated [manifest](data/SHA256SUMS) SHA-256 is `840d3c3fe76b8b374ea28f9f7a3e94a117681b72867d18e11ce5fffc6aa5a7c3`; [VERSION](data/VERSION) records the release. Each manifest line records a relative zone name and SHA-256 of its decoded TZif bytes. The IANA timezone database is [public-domain data](https://data.iana.org/time-zones/tz-link.html).

The generator accepts only safe relative names whose file content begins with `TZif`. It resolves filesystem aliases within the source directory, adds IANA aliases declared in `tzdata.zi`, rejects files over 16 MiB, and writes deterministic sorted entries across `data_first.gom` and `data_second.gom`. It excludes `posix/`, leap-second `right/`, and machine-specific `localtime` and `posixrules` entries. To refresh from a trusted installed dataset:

```sh
cd tzdata/tools
go run . /usr/share/zoneinfo 2026c ..
cd ../..
../../goml-dev/stage2/bin/goml fmt
```

Record the new package version, `tzdata.zi` hash, entry count, and manifest hash here after regeneration. Generation requires local zoneinfo; building and testing the committed GoML package do not. Embedded bytes and checksums are public source data, so applications that need authenticated releases must verify their own dependency source. The library has no automatic data update mechanism; a new release requires regenerating and publishing a new package snapshot.
