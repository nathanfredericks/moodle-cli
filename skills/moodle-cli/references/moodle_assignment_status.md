## moodle assignment status

Check submission status

### Synopsis

View submission status for the current user on an assignment.

```
moodle assignment status <assignment-id> [flags]
```

### Examples

```
  # Check submission status
  moodle assignment status 101

  # Output status as JSON
  moodle assignment status 101 -f json
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

* [moodle assignment](moodle_assignment.md)	 - Manage assignments

