## moodle forum list

List forums in a course

### Synopsis

List all forums in a course.

```
moodle forum list [flags]
```

### Examples

```
  # List forums in course 42
  moodle forum list --course 42

  # Output as JSON
  moodle forum list --course 42 -f json
```

### Options

```
      --course int   Course ID (required)
  -h, --help         help for list
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle forum](moodle_forum.md)	 - Manage forums

