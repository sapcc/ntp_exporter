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
