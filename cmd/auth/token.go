package auth

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/auth"
	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
)

func newTokenCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "token",
		Short: "Print the authentication token",
		Long:  "Print the stored authentication token.",
		Example: `  # Print the stored token
  moodle auth token

  # Use the token in a script
  curl -H "Authorization: $(moodle auth token)" https://moodle.example.com/api`,
		RunE: func(cmd *cobra.Command, args []string) error {
			token, err := f.Auth.Get(auth.TokenKey)
			if err != nil {
				return fmt.Errorf("no token found; run 'moodle auth login'")
			}

			fmt.Fprintln(f.IO.Out, token)
			return nil
		},
	}

	return cmd
}
