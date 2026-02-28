package cmdutil

import (
	"io"
	"os"

	"github.com/nathanfredericks/moodle-cli/internal/api"
	"github.com/nathanfredericks/moodle-cli/internal/auth"
	"github.com/nathanfredericks/moodle-cli/internal/config"
	"github.com/nathanfredericks/moodle-cli/internal/output"
)

// IOStreams bundles the standard I/O streams.
type IOStreams struct {
	In     io.Reader
	Out    io.Writer
	ErrOut io.Writer
}

// DefaultIOStreams returns IOStreams using os.Stdin/Stdout/Stderr.
func DefaultIOStreams() IOStreams {
	return IOStreams{
		In:     os.Stdin,
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}
}

// Factory provides shared dependencies for all commands.
type Factory struct {
	Config  config.ConfigManager
	Auth    auth.CredentialStore
	Client  func() (api.MoodleClient, error) // lazy init
	Output  func(data any, opts output.FormatOptions) error
	IO      IOStreams

	// Runtime flags (not persisted to disk)
	NoColor bool
	Verbose bool
}

// NewFactory creates a Factory with standard wiring.
func NewFactory() (*Factory, error) {
	cfg, err := config.NewConfigManager()
	if err != nil {
		return nil, err
	}

	creds := auth.NewFileCredentialStore(cfg.ConfigDir())

	f := &Factory{
		Config: cfg,
		Auth:   creds,
		IO:     DefaultIOStreams(),
		Output: func(data any, opts output.FormatOptions) error {
			return output.Print(data, opts)
		},
	}

	// Lazy client initialization
	f.Client = func() (api.MoodleClient, error) {
		site, err := f.Config.Site()
		if err != nil {
			return nil, err
		}

		return api.NewClientFromConfig(site.URL, creds)
	}

	return f, nil
}
