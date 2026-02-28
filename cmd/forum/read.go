package forum

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
	"github.com/nathanfredericks/moodle-cli/internal/output"
	"github.com/nathanfredericks/moodle-cli/internal/text"
)

type postAuthor struct {
	ID       int    `json:"id"`
	Fullname string `json:"fullname"`
}

type post struct {
	ID           int        `json:"id"`
	DiscussionID int        `json:"discussionid"`
	ParentID     int        `json:"parentid"`
	Subject      string     `json:"subject"`
	Message      string     `json:"message"`
	TimeCreated  int64      `json:"timecreated"`
	HasParent    bool       `json:"hasparent"`
	Author       postAuthor `json:"author"`
}

type postsResponse struct {
	Posts []post `json:"posts"`
}

type threadNode struct {
	post     post
	children []*threadNode
}

func newReadCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read <discussion-id>",
		Short: "Read a discussion thread",
		Long:  "Read all posts in a discussion, displayed as a threaded conversation.",
		Example: `  # Read a discussion thread
  moodle forum read 100

  # Output posts as JSON
  moodle forum read 100 -f json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			discussionID, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("invalid discussion ID: %s", args[0])
			}

			client, err := f.Client()
			if err != nil {
				return err
			}

			var result postsResponse
			params := map[string]any{
				"discussionid": discussionID,
				"sortby":       "created",
				"sortdirection": "ASC",
			}
			if err := client.Call(cmd.Context(), "mod_forum_get_discussion_posts", params, &result); err != nil {
				return fmt.Errorf("failed to read discussion: %w", err)
			}

			if len(result.Posts) == 0 {
				fmt.Fprintln(f.IO.Out, "No posts found in this discussion.")
				return nil
			}

			formatStr, _ := cmd.Flags().GetString("format")
			opts := output.FormatOptions{
				Format: output.ParseFormat(formatStr),
				Writer: f.IO.Out,
			}

			if opts.Format == output.FormatJSON || opts.Format == output.FormatYAML {
				return f.Output(result.Posts, opts)
			}

			// Build thread tree
			roots := buildThreadTree(result.Posts)

			// Render threaded view
			w := f.IO.Out
			fmt.Fprintf(w, "Discussion #%d  |  %d posts\n\n", discussionID, len(result.Posts))

			for _, root := range roots {
				renderThread(w, root, 0)
			}

			return nil
		},
	}

	return cmd
}

func buildThreadTree(posts []post) []*threadNode {
	nodeMap := make(map[int]*threadNode, len(posts))
	var roots []*threadNode

	for i := range posts {
		nodeMap[posts[i].ID] = &threadNode{post: posts[i]}
	}

	for i := range posts {
		node := nodeMap[posts[i].ID]
		if posts[i].ParentID == 0 || !posts[i].HasParent {
			roots = append(roots, node)
		} else if parent, ok := nodeMap[posts[i].ParentID]; ok {
			parent.children = append(parent.children, node)
		} else {
			roots = append(roots, node)
		}
	}

	return roots
}

func renderThread(w io.Writer, node *threadNode, depth int) {
	indent := strings.Repeat("  ", depth)
	connector := ""
	if depth > 0 {
		connector = "|-- "
	}

	// Author and timestamp line
	fmt.Fprintf(w, "%s%s%s  |  %s  |  Post #%d\n",
		indent, connector,
		node.post.Author.Fullname,
		formatTimestamp(node.post.TimeCreated),
		node.post.ID,
	)

	// Message content
	stripped := text.StripHTML(node.post.Message)
	contentIndent := indent
	if depth > 0 {
		contentIndent = indent + "|   "
	}
	for _, line := range strings.Split(stripped, "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			fmt.Fprintf(w, "%s%s\n", contentIndent, line)
		}
	}
	fmt.Fprintln(w)

	for _, child := range node.children {
		renderThread(w, child, depth+1)
	}
}
