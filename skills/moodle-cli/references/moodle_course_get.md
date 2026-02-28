## moodle course get

Get course details

### Synopsis

Get detailed information about a specific course by ID.

```
moodle course get <course-id> [flags]
```

### Examples

```
  # Get course details
  moodle course get 42

  # Get course details as JSON
  moodle course get 42 -f json
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

* [moodle course](moodle_course.md)	 - Manage courses

