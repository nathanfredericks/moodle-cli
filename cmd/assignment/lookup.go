package assignment

import (
	"context"
	"fmt"

	"github.com/nathanfredericks/moodle-cli/internal/api"
)

func lookupAssignment(ctx context.Context, client api.MoodleClient, assignID int) (*assignmentItem, string, error) {
	var listResult assignmentListResponse
	if err := client.Call(ctx, "mod_assign_get_assignments", map[string]any{}, &listResult); err != nil {
		return nil, "", fmt.Errorf("failed to get assignments: %w", err)
	}

	for _, c := range listResult.Courses {
		for i, a := range c.Assignments {
			if a.ID == assignID {
				return &c.Assignments[i], c.FullName, nil
			}
		}
	}

	return nil, "", fmt.Errorf("assignment %d not found", assignID)
}
