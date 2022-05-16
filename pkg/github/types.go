package github

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
// For more info or to see other fields that might be available, take a look
// at the webhook events and payloads for:
//   issue_comment: https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#issue_comment
//   pull_request_review_comment: https://docs.github.com/en/developers/webhooks-and-events/webhooks/webhook-events-and-payloads#pull_request_review_comment
type CommentPayload struct {
	Action      string       `json:"action"`
	Comment     Comment      `json:"comment"`
	Issue       *Issue       `json:"issue,omitempty"`
	PullRequest *PullRequest `json:"pull_request,omitempty"`
	Repository  Repository   `json:"repository"`
	Sender      Sender       `json:"sender"`
}

// Sender represents the sender of the action from the GitHub API.
type Sender struct {
	Login string `json:"login"`
}

// Repository represents a GitHub repo.
type Repository struct {
	FullName string          `json:"full_name"`
	Name     string          `json:"name"`
	Owner    RepositoryOwner `json:"owner"`
}

// RepositoryOwner represents the repo owner.
type RepositoryOwner struct {
	Login string `json:"login"`
}

// CommentUser represents the user of the comment.
type CommentUser struct {
	Login string `json:"login"`
}

// Comment is the comment on a GitHub issue from the payload.
type Comment struct {
	Body                string      `json:"body"`
	ID                  int64       `json:"id"`
	PullRequestReviewID *int64      `json:"pull_request_review_id,omitempty"`
	CommentUser         CommentUser `json:"user"`
}

// Issue represents a GitHub issue.
type Issue struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

// PullRequest represents a GitHub pull request.
type PullRequest struct {
	URL string `json:"url"`
}
