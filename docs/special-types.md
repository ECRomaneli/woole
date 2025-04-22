# Special Types

## URL Patterns

All options that requires URLs **MUST** follow one of the patterns below:

- `protocol`://`host`:`port`;
- `protocol`://`host` (default port `80`);
- `host`:`port`, (default protocol `HTTP`);
- `host`, (default protocol `HTTP` and port `80`);
- `port`, only digits (default protocol `HTTP` and host `localhost`).

## Duration Format

The **Duration Format** allows you to specify time intervals using a combination of numeric values and time unit qualifiers. This format is used in options like `-reconnect-interval` and other timeout-related configurations.

### Supported Qualifiers
- `d` - Days
- `h` - Hours
- `m` or `min` - Minutes
- `s` - Seconds
- `ms` - Milliseconds
- `ns` - Nanoseconds

### Examples
| Input String               | Description                             | Equivalent Duration         |
|----------------------------|-----------------------------------------|-----------------------------|
| `1d`                       | 1 day                                   | 24 hours                    |
| `2h 30m`                   | 2 hours and 30 minutes                  | 2 hours, 30 minutes         |
| `45s`                      | 45 seconds                              | 45 seconds                  |
| `100ms`                    | 100 milliseconds                        | 100 milliseconds            |
| `1h 15min 10s`             | 1 hour, 15 minutes, and 10 seconds      | 1 hour, 15 minutes, 10 secs |
| `0`                        | Zero duration                           | 0                           |

## Size Format

The **Size Format** allows you to specify sizes in bytes using a combination of numeric values and unit qualifiers. This format is used in options like `-tunnel-request-size` and `-tunnel-response-size`.

### Supported Qualifiers
- `b` - Bytes
- `kb` - Kilobytes (1 KB = 1024 bytes)
- `mb` - Megabytes (1 MB = 1024 KB)
- `gb` - Gigabytes (1 GB = 1024 MB)
- `tb` - Terabytes (1 TB = 1024 GB)

### Examples
| Input String     | Description  | Equivalent Size                     |
|------------------|--------------|-------------------------------------|
| `1024b`          | 1024 bytes   | 1024 bytes                          |
| `1kb`            | 1 kilobyte   | 1024 bytes                          |
| `2mb`            | 2 megabytes  | 2 * 1024 * 1024 bytes               |
| `1gb`            | 1 gigabyte   | 1 * 1024 * 1024 * 1024 bytes        |
| `1tb`            | 1 terabyte   | 1 * 1024 * 1024 * 1024 * 1024 bytes |