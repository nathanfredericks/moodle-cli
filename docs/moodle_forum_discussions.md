## moodle forum discussions

List discussions in a forum

### Synopsis

List discussions in a forum with pagination.

```
moodle forum discussions <forum-id> [flags]
```

### Examples

```
  # List discussions in forum 5
  moodle forum discussions 5

  # Paginate through discussions
  moodle forum discussions 5 --page 1 --per-page 10
```

### Options

```
  -h, --help           help for discussions
      --page int       Page number (0-indexed)
      --per-page int   Discussions per page (default 20)
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle forum](moodle_forum.md)	 - Manage forums

