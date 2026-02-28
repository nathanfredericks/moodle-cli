## moodle config set

Set a configuration value

### Synopsis

Set a configuration key to the given value. The change is persisted to disk immediately.

```
moodle config set <key> <value> [flags]
```

### Examples

```
  # Set the default output format to JSON
  moodle config set format json

  # Set the default output format to table
  moodle config set format table
```

### Options

```
  -h, --help   help for set
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle config](moodle_config.md)	 - Manage configuration

