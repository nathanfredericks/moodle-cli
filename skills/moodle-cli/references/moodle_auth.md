## moodle auth

Manage authentication

### Synopsis

Log in, log out, and manage authentication tokens for Moodle instances.

### Examples

```
  # Log in to a Moodle instance
  moodle auth login --url https://moodle.example.com

  # Check authentication status
  moodle auth status

  # Print the stored token
  moodle auth token
```

### Options

```
  -h, --help   help for auth
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle](moodle.md)	 - CLI for the Moodle LMS
* [moodle auth login](moodle_auth_login.md)	 - Log in to a Moodle instance
* [moodle auth logout](moodle_auth_logout.md)	 - Log out of the Moodle instance
* [moodle auth status](moodle_auth_status.md)	 - Show authentication status
* [moodle auth token](moodle_auth_token.md)	 - Print the authentication token

