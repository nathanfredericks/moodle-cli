## moodle forum read

Read a discussion thread

### Synopsis

Read all posts in a discussion, displayed as a threaded conversation.

```
moodle forum read <discussion-id> [flags]
```

### Examples

```
  # Read a discussion thread
  moodle forum read 100

  # Output posts as JSON
  moodle forum read 100 -f json
```

### Options

```
  -h, --help   help for read
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle forum](moodle_forum.md)	 - Manage forums

