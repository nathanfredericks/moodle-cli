## moodle assignment due

Show upcoming due dates

### Synopsis

Show upcoming assignment and activity due dates from the calendar timeline.

```
moodle assignment due [flags]
```

### Examples

```
  # Show items due in the next 7 days (default)
  moodle assignment due

  # Show overdue items
  moodle assignment due --period overdue

  # Show items due in the next 30 days
  moodle assignment due --period 30days

  # Filter by course
  moodle assignment due --course 42

  # Output as JSON
  moodle assignment due -f json
```

### Options

```
      --course int      Filter by course ID
  -h, --help            help for due
      --period string   Time period: overdue, 7days, 30days, 3months, 6months (default "7days")
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle assignment](moodle_assignment.md)	 - Manage assignments

