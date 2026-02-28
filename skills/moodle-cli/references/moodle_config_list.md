## moodle config list

List all configuration settings

### Synopsis

Display all configuration keys and their current values in a table.

```
moodle config list [flags]
```

### Examples

```
  # List all settings
  moodle config list

  # List settings as JSON
  moodle config list -f json
```

### Options

```
  -h, --help   help for list
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle config](moodle_config.md)	 - Manage configuration

