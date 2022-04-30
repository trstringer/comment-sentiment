package github

// CommentType allows the ability to distinguish different comment types.
type CommentType string

const (
	// CommentTypeIssueComment represents a GitHub issue comment.
	CommentTypeIssueComment CommentType = "issue_comment"
	// CommentTypePullRequestReviewComment represents a GitHub pull request
	// review comment.
	CommentTypePullRequestReviewComment CommentType = "pull_request_review_comment"
)

// CommentPayload represents the payload from GitHub for an issue comment.
type CommentPayload struct {
	Action  string  `json:"action"`
	Comment Comment `json:"comment"`
}

// Comment is the comment on a GitHub issue from the payload.
type Comment struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
}
