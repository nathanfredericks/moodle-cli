## moodle

CLI for the Moodle LMS

### Synopsis

A command-line interface for the Moodle Learning Management System API.

### Examples

```
  # List your enrolled courses
  moodle course list

  # Get details about a specific course
  moodle course get 42

  # Search for courses by name
  moodle course search "Introduction to Computing"

  # View your assignments
  moodle assignment list --course 42

  # Output as JSON for scripting
  moodle course list -f json
```

### Options

```
  -f, --format string   Output format: table, json, csv, yaml, plain
  -h, --help            help for moodle
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle assignment](moodle_assignment.md)	 - Manage assignments
* [moodle auth](moodle_auth.md)	 - Manage authentication
* [moodle config](moodle_config.md)	 - Manage configuration
* [moodle course](moodle_course.md)	 - Manage courses
* [moodle forum](moodle_forum.md)	 - Manage forums
* [moodle user](moodle_user.md)	 - Manage users
* [moodle version](moodle_version.md)	 - Print the version

