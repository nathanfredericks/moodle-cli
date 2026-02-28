## moodle forum reply

Reply to a post

### Synopsis

Reply to an existing post in a discussion.

```
moodle forum reply <post-id> [flags]
```

### Examples

```
  # Reply to a post
  moodle forum reply 200 --message "Thanks, that helped!"

  # Reply with a custom subject
  moodle forum reply 200 --subject "Re: Lab 3" --message "Here is my approach..."
```

### Options

```
  -h, --help             help for reply
      --message string   Reply message (required)
      --subject string   Reply subject (optional, defaults to original subject)
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle forum](moodle_forum.md)	 - Manage forums

