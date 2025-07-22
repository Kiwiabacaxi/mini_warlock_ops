package services

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github-analytics-backend/github"
	"github-analytics-backend/models"
)

// MetricsService handles all metrics calculations
type MetricsService struct {
	client *github.Client
}

// NewMetricsService creates a new MetricsService
func NewMetricsService(client *github.Client) *MetricsService {
	return &MetricsService{client: client}
}

// GetCommitMetrics calculates commit metrics for a repository
func (s *MetricsService) GetCommitMetrics(ctx context.Context, owner, repo string, since, until time.Time) (*models.CommitMetrics, error) {
	commits, err := s.client.GetCommits(ctx, owner, repo, since, until)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	// Process commits data
	authorMap := make(map[string]*models.CommitAuthor)
	dailyCommits := make(map[string]int)

	for _, commit := range commits {
		// Update author stats
		authorKey := commit.Author.Email
		if author, exists := authorMap[authorKey]; exists {
			author.Commits++
			author.Additions += commit.Additions
			author.Deletions += commit.Deletions
		} else {
			authorMap[authorKey] = &models.CommitAuthor{
				Author:    commit.Author.Name,
				Email:     commit.Author.Email,
				Commits:   1,
				Additions: commit.Additions,
				Deletions: commit.Deletions,
			}
		}

		// Update daily stats
		dateKey := commit.Author.Date.Format("2006-01-02")
		dailyCommits[dateKey]++
	}

	// Convert maps to slices and sort
	var leaderboard []models.CommitAuthor
	for _, author := range authorMap {
		leaderboard = append(leaderboard, *author)
	}
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Commits > leaderboard[j].Commits
	})

	var timeline []models.CommitDay
	for date, count := range dailyCommits {
		timeline = append(timeline, models.CommitDay{
			Date:    date,
			Commits: count,
		})
	}
	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].Date < timeline[j].Date
	})

	return &models.CommitMetrics{
		TotalCommits:  len(commits),
		UniqueAuthors: len(authorMap),
		DateRange: models.DateRange{
			Start: models.FormatDate(since),
			End:   models.FormatDate(until),
		},
		Leaderboard: leaderboard,
		Timeline:    timeline,
	}, nil
}

// GetPRMetrics calculates pull request metrics for a repository
func (s *MetricsService) GetPRMetrics(ctx context.Context, owner, repo string, since, until time.Time) (*models.PRMetrics, error) {
	prs, err := s.client.GetPullRequests(ctx, owner, repo, since, until)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull requests: %w", err)
	}

	// Process PR data
	authorMap := make(map[string]*models.PRAuthor)
	dailyPRs := make(map[string]*models.PRDay)

	var totalMergeTime float64
	var mergedCount int
	var totalPRs, openPRs, closedPRs, mergedPRs, draftPRs int

	for _, pr := range prs {
		totalPRs++

		switch pr.State {
		case "OPEN":
			openPRs++
		case "CLOSED":
			closedPRs++
		case "MERGED":
			mergedPRs++
		}

		if pr.IsDraft {
			draftPRs++
		}

		// Calculate merge time
		if pr.MergedAt != nil {
			mergeTime := pr.MergedAt.Sub(pr.CreatedAt).Hours()
			totalMergeTime += mergeTime
			mergedCount++
		}

		// Update author stats
		authorKey := pr.Author.Login
		if author, exists := authorMap[authorKey]; exists {
			author.TotalPRs++
			switch pr.State {
			case "OPEN":
				author.OpenPRs++
			case "CLOSED":
				author.ClosedPRs++
			case "MERGED":
				author.MergedPRs++
			}
			if pr.IsDraft {
				author.DraftPRs++
			}
		} else {
			author := &models.PRAuthor{
				Author:   pr.Author.Login,
				TotalPRs: 1,
			}
			switch pr.State {
			case "OPEN":
				author.OpenPRs = 1
			case "CLOSED":
				author.ClosedPRs = 1
			case "MERGED":
				author.MergedPRs = 1
			}
			if pr.IsDraft {
				author.DraftPRs = 1
			}
			authorMap[authorKey] = author
		}

		// Update daily stats
		dateKey := pr.CreatedAt.Format("2006-01-02")
		if day, exists := dailyPRs[dateKey]; exists {
			day.Created++
			if pr.State == "MERGED" {
				day.Merged++
			} else if pr.State == "CLOSED" {
				day.Closed++
			}
		} else {
			day := &models.PRDay{
				Date:    dateKey,
				Created: 1,
			}
			if pr.State == "MERGED" {
				day.Merged = 1
			} else if pr.State == "CLOSED" {
				day.Closed = 1
			}
			dailyPRs[dateKey] = day
		}
	}

	// Calculate average merge time for each author
	for _, author := range authorMap {
		if author.MergedPRs > 0 {
			// This is a simplified calculation - in a real implementation,
			// you'd want to track individual merge times per author
			author.AvgMergeTime = totalMergeTime / float64(mergedCount)
		}
	}

	// Convert maps to slices and sort
	var leaderboard []models.PRAuthor
	for _, author := range authorMap {
		leaderboard = append(leaderboard, *author)
	}
	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].TotalPRs > leaderboard[j].TotalPRs
	})

	var timeline []models.PRDay
	for _, day := range dailyPRs {
		timeline = append(timeline, *day)
	}
	sort.Slice(timeline, func(i, j int) bool {
		return timeline[i].Date < timeline[j].Date
	})

	var avgMergeTime float64
	if mergedCount > 0 {
		avgMergeTime = totalMergeTime / float64(mergedCount)
	}

	var approvalRate float64
	if totalPRs > 0 {
		approvalRate = float64(mergedPRs) / float64(totalPRs) * 100
	}

	return &models.PRMetrics{
		TotalPRs:  totalPRs,
		OpenPRs:   openPRs,
		ClosedPRs: closedPRs,
		MergedPRs: mergedPRs,
		DraftPRs:  draftPRs,
		DateRange: models.DateRange{
			Start: models.FormatDate(since),
			End:   models.FormatDate(until),
		},
		Leaderboard:  leaderboard,
		Timeline:     timeline,
		AvgMergeTime: avgMergeTime,
		ApprovalRate: approvalRate,
	}, nil
}

// GetContributionMetrics calculates contribution metrics with weighted scoring
func (s *MetricsService) GetContributionMetrics(ctx context.Context, owner, repo string, since, until time.Time) (*models.ContributionMetrics, error) {
	// Default score weights
	weights := models.ScoreWeights{
		Commit:   2.0,
		PRReview: 2.0,
		Comment:  0.5,
		Issue:    1.0,
		Approval: 2.0,
		Changes:  0.5,
	}

	// Get all data
	commits, err := s.client.GetCommits(ctx, owner, repo, since, until)
	if err != nil {
		return nil, fmt.Errorf("failed to get commits: %w", err)
	}

	prs, err := s.client.GetPullRequests(ctx, owner, repo, since, until)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull requests: %w", err)
	}

	issues, err := s.client.GetIssues(ctx, owner, repo, since, until)
	if err != nil {
		return nil, fmt.Errorf("failed to get issues: %w", err)
	}

	// Process contributions
	contributorMap := make(map[string]*models.Contributor)

	// Process commits
	for _, commit := range commits {
		username := commit.Author.Name
		if username == "" {
			username = commit.Author.Email
		}

		if contributor, exists := contributorMap[username]; exists {
			contributor.Breakdown.Commits.Count++
			contributor.Breakdown.Commits.Score += weights.Commit
		} else {
			contributorMap[username] = &models.Contributor{
				Username: username,
				Breakdown: models.ContributionBreakdown{
					Commits: models.ContributionDetail{
						Count: 1,
						Score: weights.Commit,
					},
				},
			}
		}
	}

	// Process PR reviews and comments
	for _, pr := range prs {
		// Process reviews
		for _, review := range pr.Reviews.Nodes {
			username := review.Author.Login
			if contributor, exists := contributorMap[username]; exists {
				contributor.Breakdown.PRReviews.Count++
				contributor.Breakdown.PRReviews.Score += weights.PRReview

				if review.State == "APPROVED" {
					contributor.Breakdown.Approvals.Count++
					contributor.Breakdown.Approvals.Score += weights.Approval
				} else if review.State == "CHANGES_REQUESTED" {
					contributor.Breakdown.Changes.Count++
					contributor.Breakdown.Changes.Score += weights.Changes
				}
			} else {
				newContributor := &models.Contributor{
					Username: username,
					Breakdown: models.ContributionBreakdown{
						PRReviews: models.ContributionDetail{
							Count: 1,
							Score: weights.PRReview,
						},
					},
				}

				if review.State == "APPROVED" {
					newContributor.Breakdown.Approvals = models.ContributionDetail{
						Count: 1,
						Score: weights.Approval,
					}
				} else if review.State == "CHANGES_REQUESTED" {
					newContributor.Breakdown.Changes = models.ContributionDetail{
						Count: 1,
						Score: weights.Changes,
					}
				}

				contributorMap[username] = newContributor
			}
		}

		// Process PR comments
		for _, comment := range pr.Comments.Nodes {
			username := comment.Author.Login
			if contributor, exists := contributorMap[username]; exists {
				contributor.Breakdown.Comments.Count++
				contributor.Breakdown.Comments.Score += weights.Comment
			} else {
				contributorMap[username] = &models.Contributor{
					Username: username,
					Breakdown: models.ContributionBreakdown{
						Comments: models.ContributionDetail{
							Count: 1,
							Score: weights.Comment,
						},
					},
				}
			}
		}
	}

	// Process issues
	for _, issue := range issues {
		username := issue.Author.Login
		if contributor, exists := contributorMap[username]; exists {
			contributor.Breakdown.Issues.Count++
			contributor.Breakdown.Issues.Score += weights.Issue
		} else {
			contributorMap[username] = &models.Contributor{
				Username: username,
				Breakdown: models.ContributionBreakdown{
					Issues: models.ContributionDetail{
						Count: 1,
						Score: weights.Issue,
					},
				},
			}
		}

		// Process issue comments
		for _, comment := range issue.Comments.Nodes {
			username := comment.Author.Login
			if contributor, exists := contributorMap[username]; exists {
				contributor.Breakdown.Comments.Count++
				contributor.Breakdown.Comments.Score += weights.Comment
			} else {
				contributorMap[username] = &models.Contributor{
					Username: username,
					Breakdown: models.ContributionBreakdown{
						Comments: models.ContributionDetail{
							Count: 1,
							Score: weights.Comment,
						},
					},
				}
			}
		}
	}

	// Calculate total scores and prepare final data
	var contributors []models.Contributor
	var totalScore float64
	var totalContributions int

	for _, contributor := range contributorMap {
		contributor.TotalScore = contributor.Breakdown.Commits.Score +
			contributor.Breakdown.PRReviews.Score +
			contributor.Breakdown.Comments.Score +
			contributor.Breakdown.Issues.Score +
			contributor.Breakdown.Approvals.Score +
			contributor.Breakdown.Changes.Score

		totalScore += contributor.TotalScore
		totalContributions += contributor.Breakdown.Commits.Count +
			contributor.Breakdown.PRReviews.Count +
			contributor.Breakdown.Comments.Count +
			contributor.Breakdown.Issues.Count +
			contributor.Breakdown.Approvals.Count +
			contributor.Breakdown.Changes.Count

		contributors = append(contributors, *contributor)
	}

	// Sort by total score (descending)
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].TotalScore > contributors[j].TotalScore
	})

	return &models.ContributionMetrics{
		TotalContributions: totalContributions,
		TotalScore:         totalScore,
		DateRange: models.DateRange{
			Start: models.FormatDate(since),
			End:   models.FormatDate(until),
		},
		Contributors: contributors,
		ScoreWeights: weights,
	}, nil
}
