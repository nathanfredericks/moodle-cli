## moodle auth token

Print the authentication token

### Synopsis

Print the stored authentication token.

```
moodle auth token [flags]
```

### Examples

```
  # Print the stored token
  moodle auth token

  # Use the token in a script
  curl -H "Authorization: $(moodle auth token)" https://moodle.example.com/api
```

### Options

```
  -h, --help   help for token
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle auth](moodle_auth.md)	 - Manage authentication

