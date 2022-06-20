# v2.2.0 (TBD)

Changes:

- Add `ntp_build_info` metric.

# v2.1.0 (2022-06-17)

Changes:

- New metrics were added: `ntp_rtt_seconds`, `ntp_reference_timestamp_seconds`, `ntp_root_delay_seconds`, `ntp_root_dispersion_seconds`, `ntp_root_distance_seconds`, `ntp_precision_seconds`, `ntp_leap`
- The `ntp_stratum` metric now has the label `server`, and is reported separately for each server.
- Go got updated to version 1.18.

# v2.0.2 (2022-05-04)

Changes:

- Update all dependencies to their current versions.

# v2.0.1 (2021-09-24)

Changes:

- Update all dependencies to their current versions.

# v2.0.0 (2020-08-04)

**Backwards-incompatible changes:**

- With `-ntp.source http`, the query parameter containing the NTP server has
  been renamed from `source` to `target`. With this change, you can use the
  [multi-target exporter pattern](https://prometheus.io/docs/guides/multi-target-exporter/).

Changes:

- Update all dependencies to their current versions.

# v1.1.3 (2020-05-28)

Changes:

- Update all dependencies to their current versions.

# v1.1.2 (2020-04-08)

Changes:

- Update all dependencies to their current versions.

# v1.1.1 (2020-02-10)

Changes:

- Update all dependencies to their current versions.

# v1.1.0 (2019-11-19)

New features:

- When the clock drift is unusually high, retry clock drift measurement
  multiple times and take the average, to avoid alerts because of a one-time
  mismeasurement.
- The NTP server name is now reported as a metric label.
- The option `-ntp.source` has been added. With `-ntp.source http`, the NTP
  server is not defined through command-line options, but through query
  parameters on the HTTP GET request.

Changes:

- The default port is now 9559, instead of 9100, to match our [port allocation][alloc].
- Update all dependencies to their current versions.

[alloc]: https://github.com/prometheus/prometheus/wiki/Default-port-allocations#exporters-starting-at-9100

# v1.0.0 (2017-01-13)

Initial release.
