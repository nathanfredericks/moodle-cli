## moodle course search

Search courses

### Synopsis

Search for courses by name across the Moodle instance.

```
moodle course search <query> [flags]
```

### Examples

```
  # Search for courses by name
  moodle course search "Biology"

  # Search and output as JSON
  moodle course search "CS101" -f json
```

### Options

```
  -h, --help   help for search
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle course](moodle_course.md)	 - Manage courses

