package github

import (
	"fmt"

	ghapi "github.com/google/go-github/v44/github"
)

// CommentType returns the type of comment that it is.
func (c CommentPayload) CommentType() (CommentType, error) {
	if c.Issue != nil {
		return CommentTypeIssueComment, nil
	} else if c.PullRequest != nil {
		return CommentTypePullRequestReviewComment, nil
	}

	return CommentTypeUnknown, fmt.Errorf("unable to determine comment type")
}

// UpdateComment updates the comment payload text.
func (c CommentPayload) UpdateComment(client *ghapi.Client, newComment string) error {
	commentType, err := c.CommentType()
	if err != nil {
		return fmt.Errorf("error trying to update comment: %w", err)
	}

	switch commentType {
	case CommentTypeIssueComment:
		return c.updateIssueComment(newComment)
	case CommentTypePullRequestReviewComment:
		return c.updatePullRequestReviewComment(newComment)
	default:
		return fmt.Errorf("unable to update comment due to unknown type")
	}
}

func (c CommentPayload) updateIssueComment(newComment string) error {
	return nil
}

func (c CommentPayload) updatePullRequestReviewComment(newComment string) error {
	return nil
}
