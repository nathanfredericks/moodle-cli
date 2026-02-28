## moodle assignment

Manage assignments

### Synopsis

List, view, submit, and manage assignments.

### Examples

```
  # List all assignments
  moodle assignment list

  # List assignments for a specific course
  moodle assignment list --course 42

  # Get assignment details
  moodle assignment get 101

  # Upload a file and submit
  moodle assignment upload 101 report.pdf
  moodle assignment submit 101 --accept-statement
```

### Options

```
  -h, --help   help for assignment
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle](moodle.md)	 - CLI for the Moodle LMS
* [moodle assignment download](moodle_assignment_download.md)	 - Download submission files
* [moodle assignment due](moodle_assignment_due.md)	 - Show upcoming due dates
* [moodle assignment get](moodle_assignment_get.md)	 - Get assignment details
* [moodle assignment list](moodle_assignment_list.md)	 - List assignments
* [moodle assignment status](moodle_assignment_status.md)	 - Check submission status
* [moodle assignment submit](moodle_assignment_submit.md)	 - Submit assignment for grading
* [moodle assignment text](moodle_assignment_text.md)	 - View or update online text submission
* [moodle assignment upload](moodle_assignment_upload.md)	 - Upload files to an assignment

