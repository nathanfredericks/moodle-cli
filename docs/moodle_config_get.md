## moodle config get

Get a configuration value

### Synopsis

Get the value of a specific configuration key. Returns an error if the key is not set.

```
moodle config get <key> [flags]
```

### Examples

```
  # Get the default output format
  moodle config get format
```

### Options

```
  -h, --help   help for get
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle config](moodle_config.md)	 - Manage configuration

