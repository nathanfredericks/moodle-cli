## moodle user list

List users enrolled in a course

### Synopsis

List all users enrolled in a specific course.

```
moodle user list [flags]
```

### Examples

```
  # List users in course 42
  moodle user list --course 42

  # Output as JSON
  moodle user list --course 42 -f json
```

### Options

```
      --course int   Course ID to list enrolled users for
  -h, --help         help for list
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle user](moodle_user.md)	 - Manage users

