## moodle auth login

Log in to a Moodle instance

### Synopsis

Authenticate with a Moodle instance using username and password to obtain an API token.

```
moodle auth login [flags]
```

### Examples

```
  # Log in interactively (prompts for URL, username, and password)
  moodle auth login

  # Log in with URL and username flags
  moodle auth login --url https://moodle.example.com --username jdoe
```

### Options

```
  -h, --help              help for login
  -u, --url string        Moodle instance URL
      --username string   Username
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle auth](moodle_auth.md)	 - Manage authentication

