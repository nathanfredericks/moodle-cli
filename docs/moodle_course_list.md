## moodle course list

List enrolled courses

### Synopsis

List courses the authenticated user is enrolled in. Use --recent to show recently accessed courses instead.

```
moodle course list [flags]
```

### Examples

```
  # List in-progress courses
  moodle course list

  # List all courses
  moodle course list --timeline all

  # List past courses
  moodle course list --timeline past

  # List recently accessed courses
  moodle course list --recent
```

### Options

```
  -h, --help              help for list
      --recent            Show recently accessed courses
      --timeline string   Timeline filter: all, inprogress, past, future, hidden, favourites (default "inprogress")
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle course](moodle_course.md)	 - Manage courses

