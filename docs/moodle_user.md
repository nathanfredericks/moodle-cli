## moodle user

Manage users

### Synopsis

List, view, and manage Moodle users.

### Examples

```
  # Show current user info
  moodle user whoami

  # List users in a course
  moodle user list --course 42

  # Get user details by ID
  moodle user get 7
```

### Options

```
  -h, --help   help for user
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle](moodle.md)	 - CLI for the Moodle LMS
* [moodle user get](moodle_user_get.md)	 - Get user by ID
* [moodle user list](moodle_user_list.md)	 - List users enrolled in a course
* [moodle user whoami](moodle_user_whoami.md)	 - Show current user info

