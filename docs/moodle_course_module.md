## moodle course module

Get course module details

### Synopsis

Display detailed information about a specific course module.

```
moodle course module <module-id> [flags]
```

### Examples

```
  # Get module details
  moodle course module 1001

  # Get module details as JSON
  moodle course module 1001 -f json
```

### Options

```
  -h, --help   help for module
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle course](moodle_course.md)	 - Manage courses

