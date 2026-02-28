## moodle auth status

Show authentication status

### Synopsis

Display the current authentication status by querying the Moodle instance.

```
moodle auth status [flags]
```

### Examples

```
  # Show current authentication status
  moodle auth status

  # Output status as JSON
  moodle auth status -f json
```

### Options

```
  -h, --help   help for status
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle auth](moodle_auth.md)	 - Manage authentication

