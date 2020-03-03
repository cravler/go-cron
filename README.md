# go-cron

Cron is a system daemon used to execute desired tasks (in the background) at designated times.

Binaries can be found at: https://github.com/cravler/go-cron/releases

## Usage

```sh
cron --help
```

`crontab.yaml` example:

```yaml
exp: "@every 1h15m5s"
cmd:
    name: "ls"
    argv:
        - "-lah"
        - "--color=auto"
    dir: "/var/www"
log:
    # Combined output (stdout & stderr)
    file: "/var/log/app-1.log"
---
exp: "5 15 */1 * * ?"
# Skip command run, if previous still running
lock: true
cmd:
    name: "env"
    env:
        - "TEST1=123"
        - "TEST2=456 789"
log:
    # If file empty, we can specify (stdout & stderr)
    # Defaults to "".
    file: ""
    # The maximum size of the log before it is rolled.
    # A positive integer plus a modifier representing 
    # the  unit of measure (k, m, or g).
    # Defaults to 0 (unlimited).
    size: 1m
    # The maximum number of log files that can be present.
    # If rolling the logs creates excess files,
    # the oldest file is removed.
    # Only effective when size is also set.
    # A positive integer.
    # Defaults to 1.
    num: 3
stdout:
    file: "/var/log/app-2-stdout.log"
stderr:
    file: "/var/log/app-2-stderr.log"
    size: 5m
    num: 1
```

### CRON Expression Format

 A cron expression represents a set of times, using 5-6 space-separated fields.

Field name   | Mandatory? | Allowed values      | Allowed special characters
-------------|:----------:|:-------------------:|:--------------------------:
Seconds      | No         | `0-59`              | `* / , -`
Minutes      | Yes        | `0-59`              | `* / , -`
Hours        | Yes        | `0-23`              | `* / , -`
Day of month | Yes        | `1-31`              | `* / , - ?`
Month        | Yes        | `1-12` or `JAN-DEC` | `* / , -`
Day of week  | Yes        | `0-6` or `SUN-SAT`  | `* / , - ?`

`Month` and `Day of week` field values are case insensitive. `SUN`, `Sun`, and `sun` are equally accepted.

#### Special Characters

##### Asterisk `*`

The asterisk indicates that the cron expression will match for all values of the field; e.g., using an asterisk in 
the `month` field would indicate every month.

##### Slash `/`

Slashes are used to describe increments of ranges. For example `3-59/15` in the `minutes` field would indicate 
the 3rd minute of the hour and every 15 minutes thereafter. The form `*/...` is equivalent to the form `first-last/...`,
that is, an increment over the largest possible range of the field. The form `N/...` is accepted as meaning `N-MAX/...`,
that is, starting at N, use the increment until the end of that specific range. It does not wrap around.

##### Comma `,`

Commas are used to separate items of a list. For example, using `MON,WED,FRI` in the `day of week` field would mean 
Mondays, Wednesdays and Fridays.

##### Hyphen `-`

Hyphens are used to define ranges. For example, `9-17` would indicate every hour between 9am and 5pm inclusive.

##### Question mark `?`

Question mark may be used instead of `*` for leaving either `day of month` or `day of week` blank.

#### Predefined schedules

You may use one of several pre-defined schedules in place of a cron expression.

Entry                  | Description                                | Equivalent To
-----------------------|--------------------------------------------|:-------------:
@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | `0 0 1 1 *`
@monthly               | Run once a month, midnight, first of month | `0 0 1 * *`
@weekly                | Run once a week, midnight between Sat/Sun  | `0 0 * * 0`
@daily (or @midnight)  | Run once a day, midnight                   | `0 0 * * *`
@hourly                | Run once an hour, beginning of hour        | `0 * * * *`

##### @every `<duration>`

where `duration` is a string with time units: `s`, `m`, `h`.

Entry          | Description                                | Equivalent To
---------------|--------------------------------------------|:--------------:
@every 5s      | Run every 5 seconds                        | `*/5 * * * * ?`
@every 15m5s   | Run every 15 minutes 5 seconds             | `5 */15 * * * ?`
@every 1h15m5s | Run every 1 hour 15 minutes 5 seconds      | `5 15 */1 * * ?`

#### Time zones

By default, all scheduling is done in the machine's local time zone.
Individual cron schedules may also override the time zone they are to be interpreted in by providing an additional 
space-separated field at the beginning of the cron spec, of the form `CRON_TZ=Europe/London`.

For example:

```yaml
exp: "CRON_TZ=Europe/London 0 6 * * ?" # Runs at 6am in Europe/London
```

Be aware that jobs scheduled during daylight-savings leap-ahead transitions will not be run!
