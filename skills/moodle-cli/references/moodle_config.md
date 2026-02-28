## moodle config

Manage configuration

### Synopsis

Get, set, and manage CLI configuration.

### Examples

```
  # List all configuration settings
  moodle config list

  # Get a specific configuration value
  moodle config get format

  # Set the default output format
  moodle config set format json
```

### Options

```
  -h, --help   help for config
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle](moodle.md)	 - CLI for the Moodle LMS
* [moodle config get](moodle_config_get.md)	 - Get a configuration value
* [moodle config list](moodle_config_list.md)	 - List all configuration settings
* [moodle config set](moodle_config_set.md)	 - Set a configuration value

