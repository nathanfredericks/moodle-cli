## moodle assignment submit

Submit assignment for grading

### Synopsis

Submit an assignment for grading. The assignment must have content saved first.

```
moodle assignment submit <assignment-id> [flags]
```

### Examples

```
  # Submit an assignment for grading
  moodle assignment submit 101

  # Submit and accept the submission statement
  moodle assignment submit 101 --accept-statement
```

### Options

```
      --accept-statement   Accept the submission statement
  -h, --help               help for submit
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle assignment](moodle_assignment.md)	 - Manage assignments

