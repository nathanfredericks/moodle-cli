## moodle forum post

Create a new discussion

### Synopsis

Create a new discussion thread in a forum.

```
moodle forum post <forum-id> [flags]
```

### Examples

```
  # Create a new discussion
  moodle forum post 5 --subject "Help with Lab 3" --message "I'm stuck on step 2..."
```

### Options

```
  -h, --help             help for post
      --message string   Discussion message (required)
      --subject string   Discussion subject (required)
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle forum](moodle_forum.md)	 - Manage forums

