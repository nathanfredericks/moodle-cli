## moodle user whoami

Show current user info

### Synopsis

Display information about the currently authenticated user.

```
moodle user whoami [flags]
```

### Examples

```
  # Show current user info
  moodle user whoami

  # Output as JSON
  moodle user whoami -f json
```

### Options

```
  -h, --help   help for whoami
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle user](moodle_user.md)	 - Manage users

