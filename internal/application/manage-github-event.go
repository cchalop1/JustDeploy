package application

import "cchalop1.com/deploy/internal/api/service"

type GithubPushEvent struct {
	Ref        string `json:"ref"`
	Deleted    bool   `json:"deleted"`
	Repository struct {
		FullName string `json:"full_name"`
	} `json:"repository"`
}

// Caller holds Operation; never dispatch unverified events to this method.
func ManageVerifiedGithubEvent(ds *service.DeployService, event GithubPushEvent) error {
	services, err := ds.DatabaseAdapter.ReadServices()
	if err != nil {
		return err
	}
	for _, s := range services {
		branch := s.Branch
		if branch == "" {
			branch = "main"
		}
		if s.Type == "github_repo" && s.FullName == event.Repository.FullName && event.Ref == "refs/heads/"+branch {
			if err := redeployService(ds, s.Id); err != nil {
				return err
			}
		}
	}
	return nil
}
