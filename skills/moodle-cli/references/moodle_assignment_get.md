## moodle assignment get

Get assignment details

### Synopsis

Get detailed information about a specific assignment including grade and feedback.

```
moodle assignment get <assignment-id> [flags]
```

### Examples

```
  # Get assignment details
  moodle assignment get 101

  # Get assignment details as JSON
  moodle assignment get 101 -f json
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

* [moodle assignment](moodle_assignment.md)	 - Manage assignments

