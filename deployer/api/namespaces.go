package api

import (
	"fmt"
	"log/slog"

	"github.com/okteto-community/redeploy-applications/deployer/model"
)

const (
	developmentNamespaceType = "development"
)

// GetNamespaces retrieves all the namespaces
func GetNamespaces(baseURL, token string, logger *slog.Logger) ([]model.Namespace, error) {
	namespacesURL := fmt.Sprintf("https://%s/%s?type=%s", baseURL, namespacesAPIPath, developmentNamespaceType)
	var namespaces []model.Namespace
	if err := sendRequest(namespacesURL, token, &namespaces, logger); err != nil {
		return nil, err
	}
	return namespaces, nil
}
