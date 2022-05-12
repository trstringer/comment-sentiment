package github

import (
	"context"
	"fmt"

	ghapi "github.com/google/go-github/v44/github"
	"github.com/rs/zerolog/log"
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
		log.Debug().Msg("Updating comment with type issue")
		return c.updateIssueComment(client, newComment)
	case CommentTypePullRequestReviewComment:
		log.Debug().Msg("Updating comment with type pull request review")
		return c.updatePullRequestReviewComment(client, newComment)
	default:
		log.Error().Msg("Unknown comment type")
		return fmt.Errorf("unable to update comment due to unknown type")
	}
}

func (c CommentPayload) updateIssueComment(client *ghapi.Client, newComment string) error {
	_, _, err := client.Issues.EditComment(
		context.Background(),
		c.Repository.Owner.Login,
		c.Repository.Name,
		c.Comment.ID,
		&ghapi.IssueComment{Body: &newComment},
	)
	if err != nil {
		return fmt.Errorf("error updating issue comment: %w", err)
	}

	return nil
}

func (c CommentPayload) updatePullRequestReviewComment(client *ghapi.Client, newComment string) error {
	return nil
}
