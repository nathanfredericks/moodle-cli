## moodle assignment download

Download submission files

### Synopsis

Download all files from your submission for an assignment.

```
moodle assignment download <assignment-id> [flags]
```

### Examples

```
  # Download submission files
  moodle assignment download 101

  # Download to a specific directory
  moodle assignment download 101 -o ./submissions

  # Overwrite existing files
  moodle assignment download 101 -o ./submissions -F
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

* [moodle assignment](moodle_assignment.md)	 - Manage assignments

