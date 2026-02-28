## moodle forum

Manage forums

### Synopsis

List forums, read discussions, and post replies.

### Examples

```
  # List forums in a course
  moodle forum list --course 42

  # List discussions in a forum
  moodle forum discussions 5

  # Read a discussion thread
  moodle forum read 100

  # Create a new discussion
  moodle forum post 5 --subject "Question" --message "Hello"
```

### Options

```
  -h, --help   help for forum
```

### Options inherited from parent commands

```
  -f, --format string   Output format: table, json, csv, yaml, plain
      --no-color        Disable color output
  -v, --verbose         Enable verbose output
```

### SEE ALSO

* [moodle](moodle.md)	 - CLI for the Moodle LMS
* [moodle forum delete](moodle_forum_delete.md)	 - Delete a post
* [moodle forum discussions](moodle_forum_discussions.md)	 - List discussions in a forum
* [moodle forum edit](moodle_forum_edit.md)	 - Edit a post
* [moodle forum list](moodle_forum_list.md)	 - List forums in a course
* [moodle forum post](moodle_forum_post.md)	 - Create a new discussion
* [moodle forum read](moodle_forum_read.md)	 - Read a discussion thread
* [moodle forum reply](moodle_forum_reply.md)	 - Reply to a post

