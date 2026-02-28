package auth

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nathanfredericks/moodle-cli/internal/auth"
	"github.com/nathanfredericks/moodle-cli/internal/cmdutil"
)

func newLogoutCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Log out of the Moodle instance",
		Long:  "Remove the stored authentication token and site configuration.",
		Example: `  # Log out of the current Moodle instance
  moodle auth logout`,
		RunE: func(cmd *cobra.Command, args []string) error {
			site, err := f.Config.Site()
			if err != nil {
				return fmt.Errorf("not logged in")
			}

			// Delete the token
			if err := f.Auth.Delete(auth.TokenKey); err != nil {
				return fmt.Errorf("failed to remove token: %w", err)
			}

			// Delete the site config
			if err := f.Config.DeleteSite(); err != nil {
				return fmt.Errorf("failed to remove site: %w", err)
			}

			fmt.Fprintf(f.IO.Out, "Logged out of %s\n", site.URL)
			return nil
		},
	}

	return cmd
}
