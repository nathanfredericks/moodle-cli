## moodle assignment upload

Upload files to an assignment

### Synopsis

Upload one or more files to an assignment's submission draft area and save the submission.

```
moodle assignment upload <assignment-id> <file>... [flags]
```

### Examples

```
  # Upload a single file
  moodle assignment upload 101 report.pdf

  # Upload multiple files
  moodle assignment upload 101 report.pdf appendix.pdf data.csv
```

### Options

```
  -h, --help   help for upload
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle assignment](moodle_assignment.md)	 - Manage assignments

