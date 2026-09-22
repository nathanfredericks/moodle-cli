package main

import (
	"context"
	"fmt"
	"log"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/nathanfredericks/moodle-cli/internal/lambdaconfig"
	"github.com/nathanfredericks/moodle-cli/internal/mcpserver"
)

// Version is set at build time.
var Version = "dev"

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	secretID := os.Getenv("MOODLE_SECRET_ID")
	if secretID == "" {
		return fmt.Errorf("MOODLE_SECRET_ID environment variable is required")
	}

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("load AWS configuration: %w", err)
	}
	runtimeCfg, err := lambdaconfig.Load(ctx, secretsmanager.NewFromConfig(awsCfg), secretID)
	if err != nil {
		return err
	}
	if err := runtimeCfg.ApplyEnvironment(); err != nil {
		return err
	}

	return mcpserver.Serve(mcpserver.Options{
		Addr:    ":8080",
		APIKey:  runtimeCfg.MCPAPIKey,
		Version: Version,
	})
}
