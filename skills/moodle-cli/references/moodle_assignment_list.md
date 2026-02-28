## moodle assignment list

List assignments

### Synopsis

List assignments across enrolled courses.

```
moodle assignment list [flags]
```

### Examples

```
  # List all assignments
  moodle assignment list

  # List assignments for a specific course
  moodle assignment list --course 42

  # Output as JSON
  moodle assignment list -f json
```

### Options

```
      --course int   Filter by course ID
  -h, --help         help for list
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle assignment](moodle_assignment.md)	 - Manage assignments

