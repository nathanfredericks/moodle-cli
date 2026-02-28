## moodle forum edit

Edit a post

### Synopsis

Edit the subject and/or message of an existing post.

```
moodle forum edit <post-id> [flags]
```

### Examples

```
  # Edit a post's message
  moodle forum edit 200 --message "Updated content"

  # Edit both subject and message
  moodle forum edit 200 --subject "New Title" --message "New content"
```

### Options

```
  -h, --help             help for edit
      --message string   New message
      --subject string   New subject
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle forum](moodle_forum.md)	 - Manage forums

