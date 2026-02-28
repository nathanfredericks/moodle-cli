## moodle course download

Download module files

### Synopsis

Download all files from a course module (e.g. a resource or folder activity).

```
moodle course download <course-id> <module-id> [flags]
```

### Examples

```
  # Download files from a module
  moodle course download 42 1001

  # Download to a specific directory
  moodle course download 42 1001 -o ./downloads

  # Overwrite existing files
  moodle course download 42 1001 -o ./downloads -F
```

### Options

```
  -F, --force           Overwrite existing files
  -h, --help            help for download
  -o, --output string   Output directory
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle course](moodle_course.md)	 - Manage courses

