# v1.1.0 (TBD)

New features:

- When the clock drift is unusually high, retry clock drift measurement
  multiple times and take the average, to avoid alerts because of a one-time
  mismeasurement.
- The NTP server name is now reported as a metric label.

Changes:

- The default port is now 9559, instead of 9100, to match our [port allocation][alloc].
- Update all dependencies to their current versions.

[alloc]: https://github.com/prometheus/prometheus/wiki/Default-port-allocations#exporters-starting-at-9100

# v1.0.0 (2017-01-13)

Initial release.
