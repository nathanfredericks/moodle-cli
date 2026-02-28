## moodle assignment text

View or update online text submission

### Synopsis

View the current online text submission, or update it with --save or --stdin.

```
moodle assignment text <assignment-id> [flags]
```

### Examples

```
  # View the current online text submission
  moodle assignment text 101

  # Save text directly
  moodle assignment text 101 --save "My submission text"

  # Pipe text from a file
  cat essay.txt | moodle assignment text 101 --stdin
```

### Options

```
  -h, --help          help for text
      --save string   Save this text as the online text submission
      --stdin         Read online text from stdin
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle assignment](moodle_assignment.md)	 - Manage assignments

