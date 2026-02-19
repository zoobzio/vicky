package handlers

import "github.com/zoobzio/rocco"

// All returns all admin API handlers for registration.
// All handlers require authentication (enforced via WithAuthentication).
func All() []rocco.Endpoint {
	return []rocco.Endpoint{
		// Users
		ListUsers.WithAuthentication(),
		GetUser.WithAuthentication(),
		UpdateUser.WithAuthentication(),
		DeleteUser.WithAuthentication(),

		// Repositories
		ListRepositories.WithAuthentication(),
		GetRepository.WithAuthentication(),
		UpdateRepository.WithAuthentication(),
		DeleteRepository.WithAuthentication(),

		// Jobs
		ListJobs.WithAuthentication(),
		GetJob.WithAuthentication(),
		CancelJob.WithAuthentication(),
		RetryJob.WithAuthentication(),
		GetJobStats.WithAuthentication(),
	}
}
