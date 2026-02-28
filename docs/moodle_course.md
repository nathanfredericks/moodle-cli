## moodle course

Manage courses

### Synopsis

List, view, and manage Moodle courses.

### Examples

```
  # List enrolled courses
  moodle course list

  # Get course details
  moodle course get 42

  # View course contents
  moodle course content 42

  # Search for a course
  moodle course search "Biology 101"
```

### Options

```
  -h, --help   help for course
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle](moodle.md)	 - CLI for the Moodle LMS
* [moodle course content](moodle_course_content.md)	 - Get course contents
* [moodle course download](moodle_course_download.md)	 - Download module files
* [moodle course get](moodle_course_get.md)	 - Get course details
* [moodle course grades](moodle_course_grades.md)	 - View your course grades
* [moodle course list](moodle_course_list.md)	 - List enrolled courses
* [moodle course module](moodle_course_module.md)	 - Get course module details
* [moodle course search](moodle_course_search.md)	 - Search courses

