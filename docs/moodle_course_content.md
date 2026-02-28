## moodle course content

Get course contents

### Synopsis

Display the sections and activities of a course.

```
moodle course content <course-id> [flags]
```

### Examples

```
  # View course sections and activities
  moodle course content 42

  # Output as JSON for processing
  moodle course content 42 -f json
```

### Options

```
  -h, --help   help for content
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle course](moodle_course.md)	 - Manage courses

