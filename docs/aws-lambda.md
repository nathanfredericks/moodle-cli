# Deploying Moodle MCP to AWS Lambda

The production MCP server runs as an ARM64 Lambda container behind a response-streaming Function URL. AWS SAM builds the image, uploads it to the private `moodle-mcp` ECR repository, and deploys the `moodle-mcp` CloudFormation stack in `ca-central-1`.

The Function URL uses AWS `NONE` authorization so ordinary MCP clients can connect. The application itself requires `Authorization: Bearer <MCP_API_KEY>` on `/mcp`. The `/healthz` endpoint is intentionally public and returns no configuration data.

## One-time bootstrap

Use the scoped deployment profile, never the default root credentials:

```bash
aws cloudformation deploy \
  --profile steamspace-geocaching-deployment \
  --region ca-central-1 \
  --stack-name moodle-mcp-bootstrap \
  --template-file infrastructure/bootstrap.yaml \
  --capabilities CAPABILITY_NAMED_IAM
```

This creates the private ECR repository, `MoodleMCPGitHubDeploymentBroker` user, and `MoodleMCPGitHubDeploymentRole`. The broker can only assume that deployment role.

Create the `moodle-mcp` Secrets Manager secret as a JSON `SecretString`:

```json
{
  "MOODLE_URL": "https://moodle.acadiau.ca",
  "MOODLE_TOKEN": "...",
  "MCP_API_KEY": "..."
}
```

Generate `MCP_API_KEY` with a cryptographically secure random generator. Do not pass any secret value as a CloudFormation or SAM parameter and do not commit it to the repository.

Create one access key for the broker and store its two values as the GitHub repository secrets `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`. The local copy of the key response should be destroyed immediately after both secrets are installed. The workflow uses those credentials only to assume `MoodleMCPGitHubDeploymentRole`.

## Continuous deployment

Pull requests run tests, validate both templates, and build the Lambda image without deploying. Successful pushes to `main` deploy automatically. The workflow can also be started manually from GitHub Actions.

The deploy job:

1. Assumes the dedicated deployment role.
2. Builds the Lambda image with AWS SAM.
3. Pushes the image to the private `moodle-mcp` ECR repository.
4. Updates the `moodle-mcp` stack with rollback enabled.

Production deployment jobs are serialized. ECR retains the ten most recent images.

## Local verification

Run the unit tests and validate the infrastructure:

```bash
go test ./... -count=1
make sam-validate
make sam-build
```

The original container remains available for local MCP use:

```bash
docker build -t moodle-cli .
docker run --rm -p 8080:8080 \
  -e MOODLE_URL \
  -e MOODLE_TOKEN \
  -e MCP_API_KEY \
  moodle-cli
```

The Lambda image reads its credentials from Secrets Manager and is built with `make build-lambda`.

## Endpoints and smoke tests

Read the deployed endpoints from CloudFormation:

```bash
aws cloudformation describe-stacks \
  --profile steamspace-geocaching-deployment \
  --region ca-central-1 \
  --stack-name moodle-mcp \
  --query 'Stacks[0].Outputs'
```

Verify the health endpoint, confirm an unauthenticated MCP request returns `401`, then initialize an MCP client with the secret's bearer key and run a read-only tool such as `user_whoami`.

Lambda streaming requests have a 15-minute ceiling. MCP clients should reconnect when a stream closes.

## Secret rotation

Update the existing `moodle-mcp` secret in place. Warm Lambda environments cache the secret loaded at startup, so redeploy the function after rotation to recycle execution environments. Never log or echo the secret JSON.

## Rollback

Re-run the workflow from a known-good commit. CloudFormation updates the function to that commit's image, and ECR retains the recent image history needed for rollback.
