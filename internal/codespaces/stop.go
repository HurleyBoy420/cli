package codespaces

import (
	"context"
	"errors"
	"fmt"

	"github.com/cli/cli/v2/internal/codespaces/api"
)

func (a *App) StopCodespace(ctx context.Context, codespaceName string) error {
	if codespaceName == "" {
		codespaces, err := a.apiClient.ListCodespaces(ctx, -1)
		if err != nil {
			return fmt.Errorf("failed to list codespaces: %w", err)
		}

		var runningCodespaces []*api.Codespace
		for _, c := range codespaces {
			cs := codespace{c}
			if cs.running() {
				runningCodespaces = append(runningCodespaces, c)
			}
		}
		if len(runningCodespaces) == 0 {
			return errors.New("no running codespaces")
		}

		codespace, err := a.chooseCodespaceFromList(ctx, runningCodespaces)
		if err != nil {
			return fmt.Errorf("failed to choose codespace: %w", err)
		}
		codespaceName = codespace.Name
	} else {
		c, err := a.apiClient.GetCodespace(ctx, codespaceName, false)
		if err != nil {
			return fmt.Errorf("failed to get codespace: %q: %w", codespaceName, err)
		}
		cs := codespace{c}
		if !cs.running() {
			return fmt.Errorf("codespace %q is not running", codespaceName)
		}
	}

	if err := a.apiClient.StopCodespace(ctx, codespaceName); err != nil {
		return fmt.Errorf("failed to stop codespace: %w", err)
	}
	a.logger.Println("Codespace stopped")

	return nil
}
