package main

import (
	"fmt"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"time"

	"github.com/okteto-community/redeploy-applications/deployer/api"
	"github.com/okteto-community/redeploy-applications/deployer/git"
)

const redeployAppCommandTemplate = "okteto pipeline deploy -n \"%s\" --name \"%s\" --repository \"%s\" --branch \"%s\" --reuse-params --wait=false"

func main() {
	token := os.Getenv("OKTETO_TOKEN")
	oktetoURL := os.Getenv("OKTETO_URL")
	targetRepo := os.Getenv("TARGET_REPOSITORY")
	targetBranch := os.Getenv("TARGET_BRANCH")

	logLevel := &slog.LevelVar{} // INFO
	opts := &slog.HandlerOptions{
		Level: logLevel,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))

	if token == "" || oktetoURL == "" || targetRepo == "" {
		logger.Error("OKTETO_TOKEN, OKTETO_URL and TARGET_REPOSITORY environment variables are required")
		os.Exit(1)
	}

	u, err := url.Parse(oktetoURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Invalid OKTETO_URL %s", err))
		os.Exit(1)
	}

	nsList, err := api.GetNamespaces(u.Host, token, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("There was an error requesting the namespaces: %s", err))
		os.Exit(1)
	}

	logger.Info(fmt.Sprintf("Looking for dev environments with repository %q", targetRepo))

	// We check for applications that were last updated more than 24 hours ago
	updateThreshold := time.Now().Add(-time.Hour * 24)
	for _, ns := range nsList {
		logger.Info(fmt.Sprintf("Processing namespace '%s'", ns.Name))

		applications, err := api.GetApplicationsWithinNamespace(u.Host, token, ns.Name, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("There was an error requesting the applications within namespace '%s': %s", ns.Name, err))
			logger.Info("-----------------------------------------------")
			continue
		}

		for _, app := range applications {
			if app.Repository == "" {
				logger.Info(fmt.Sprintf("Skipping application '%s' within namespace '%s' as does not have a repository", app.Name, ns.Name))
				continue
			}

			if !git.AreSameRepository(app.Repository, targetRepo) {
				logger.Info(fmt.Sprintf("Skipping application '%s' within namespace '%s' as repository doesn't match", app.Name, ns.Name))
				continue
			}

			if targetBranch != "" && app.Branch != targetBranch {
				logger.Info(fmt.Sprintf("Skipping application '%s' within namespace '%s' as deployed branch doesn't match", app.Name, ns.Name))
				continue
			}

			if app.LastUpdated.After(updateThreshold) {
				logger.Info(fmt.Sprintf("Skipping application '%s' within namespace '%s' as it was updated recently", app.Name, ns.Name))
				continue
			}

			logger.Info(fmt.Sprintf("Redeploying application '%s' within namespace '%s'", app.Name, ns.Name))

			out, err := redeployApp(ns.Name, app.Name, app.Repository, app.Branch)
			if err != nil {
				logger.Error(fmt.Sprintf("There was an error redeploying the application '%s' within namespace '%s': %s", app.Name, ns.Name, err))
			} else {
				logger.Info(out)
			}
		}
		logger.Info("-----------------------------------------------")
	}
}

// redeployApp executes the Okteto CLI command to redeploy an application
func redeployApp(ns, appName, repo, branch string) (string, error) {
	cmd := exec.Command("bash", "-c", fmt.Sprintf(redeployAppCommandTemplate, ns, appName, repo, branch))

	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
}
