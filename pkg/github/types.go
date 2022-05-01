package github

import (
	ghapi "github.com/google/go-github/v44/github"
)

// CommentType allows the ability to distinguish different comment types.
type CommentType int

const (
	// CommentTypeIssueComment represents a GitHub issue comment.
	CommentTypeIssueComment CommentType = iota
	// CommentTypePullRequestReviewComment represents a GitHub pull request
	// review comment.
	CommentTypePullRequestReviewComment
	// CommentTypeUnknown is the indication that the comment type is unknown
	// and the output should likely not be trusted.
	CommentTypeUnknown
)

// CommentPayload represents the payload from GitHub for an issue comment.
type CommentPayload struct {
	Action      string       `json:"action"`
	Comment     Comment      `json:"comment"`
	Issue       *Issue       `json:"issue,omitempty"`
	PullRequest *PullRequest `json:"pull_request,omitempty"`
	Repository  Repository   `json:"repository"`
}

// Repository represents a GitHub repo.
type Repository struct {
	Owner RepositoryOwner `json:"owner"`
}

// RepositoryOwner represents the repo owner.
type RepositoryOwner struct {
	Login string `json:"login"`
}

// Comment is the comment on a GitHub issue from the payload.
type Comment struct {
	Body                string `json:"body"`
	ID                  int    `json:"id"`
	PullRequestReviewID *int   `json:"pull_request_review_id,omitempty"`
}

// Issue represents a GitHub issue.
type Issue struct {
	URL string `json:"url"`
}

// PullRequest represents a GitHub pull request.
type PullRequest struct {
	URL string `json:"url"`
}

// App is a representation of the GitHub App.
type App struct {
	client *ghapi.Client
	appID  int
}
