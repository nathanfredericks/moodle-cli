## moodle user get

Get user by ID

### Synopsis

Get detailed information about a user by their numeric ID.

```
moodle user get <user-id> [flags]
```

### Examples

```
  # Get user details
  moodle user get 7

  # Output as JSON
  moodle user get 7 -f json
```

### Options

```
  -h, --help   help for get
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle user](moodle_user.md)	 - Manage users

