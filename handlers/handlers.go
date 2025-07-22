package handlers

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github-analytics-backend/github"
	"github-analytics-backend/models"
	"github-analytics-backend/services"

	"github.com/gin-gonic/gin"
)

// Handler contains all the HTTP handlers
type Handler struct {
	githubClient   *github.Client
	metricsService *services.MetricsService
}

// NewHandler creates a new handler instance
func NewHandler(githubClient *github.Client) *Handler {
	return &Handler{
		githubClient:   githubClient,
		metricsService: services.NewMetricsService(githubClient),
	}
}

// parseQueryParams parses common query parameters
func (h *Handler) parseQueryParams(c *gin.Context) (*models.QueryParams, error) {
	var params models.QueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		return nil, err
	}
	return &params, nil
}

// parseDateRange parses start and end dates with defaults
func (h *Handler) parseDateRange(startDate, endDate string) (time.Time, time.Time, error) {
	var since, until time.Time
	var err error

	if startDate == "" {
		startDate = "2025-07-01"
	}
	if endDate == "" {
		endDate = "2025-07-31"
	}

	since, err = models.ParseDate(startDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	until, err = models.ParseDate(endDate)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Add 24 hours to until date to include the entire end date
	until = until.Add(24 * time.Hour)

	return since, until, nil
}

// GetRepositories returns repositories for an owner (user or organization)
func (h *Handler) GetRepositories(c *gin.Context) {
	owner := c.Param("owner")
	if owner == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "owner is required",
		})
		return
	}

	// Parse limit parameter
	limitStr := c.DefaultQuery("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	repos, err := h.githubClient.GetRepositories(c.Request.Context(), owner, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Convert to our model format
	var result []models.Repository
	for _, repo := range repos {
		language := ""
		if len(repo.Languages.Nodes) > 0 {
			language = repo.Languages.Nodes[0].Name
		}

		result = append(result, models.Repository{
			Owner:       repo.Owner.Login,
			Name:        repo.Name,
			FullName:    repo.FullName,
			Description: repo.Description,
			Private:     repo.IsPrivate,
			Fork:        repo.IsFork,
			CreatedAt:   repo.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   repo.UpdatedAt.Format(time.RFC3339),
			Language:    language,
			Stars:       repo.StargazerCount,
			Forks:       repo.ForkCount,
		})
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    result,
	})
}

// GetCommitMetrics returns commit metrics for a repository
func (h *Handler) GetCommitMetrics(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	if owner == "" || repo == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "owner and repo are required",
		})
		return
	}

	params, err := h.parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	since, until, err := h.parseDateRange(params.StartDate, params.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid date format, use YYYY-MM-DD",
		})
		return
	}

	metrics, err := h.metricsService.GetCommitMetrics(c.Request.Context(), owner, repo, since, until)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    metrics,
	})
}

// GetPRMetrics returns pull request metrics for a repository
func (h *Handler) GetPRMetrics(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	if owner == "" || repo == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "owner and repo are required",
		})
		return
	}

	params, err := h.parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	since, until, err := h.parseDateRange(params.StartDate, params.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid date format, use YYYY-MM-DD",
		})
		return
	}

	metrics, err := h.metricsService.GetPRMetrics(c.Request.Context(), owner, repo, since, until)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    metrics,
	})
}

// GetContributionMetrics returns contribution metrics for a repository
func (h *Handler) GetContributionMetrics(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	if owner == "" || repo == "" {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "owner and repo are required",
		})
		return
	}

	params, err := h.parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	since, until, err := h.parseDateRange(params.StartDate, params.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid date format, use YYYY-MM-DD",
		})
		return
	}

	metrics, err := h.metricsService.GetContributionMetrics(c.Request.Context(), owner, repo, since, until)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    metrics,
	})
}

// GetCommitsLeaderboard returns commits leaderboard chart data
func (h *Handler) GetCommitsLeaderboard(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	params, err := h.parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	since, until, err := h.parseDateRange(params.StartDate, params.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid date format, use YYYY-MM-DD",
		})
		return
	}

	metrics, err := h.metricsService.GetCommitMetrics(c.Request.Context(), owner, repo, since, until)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Format for chart.js
	chartData := map[string]interface{}{
		"type": "bar",
		"data": map[string]interface{}{
			"labels": func() []string {
				var labels []string
				for _, author := range metrics.Leaderboard {
					labels = append(labels, author.Author)
				}
				return labels
			}(),
			"datasets": []map[string]interface{}{
				{
					"label":           "Commits",
					"backgroundColor": "rgba(54, 162, 235, 0.6)",
					"borderColor":     "rgba(54, 162, 235, 1)",
					"borderWidth":     1,
					"data": func() []int {
						var data []int
						for _, author := range metrics.Leaderboard {
							data = append(data, author.Commits)
						}
						return data
					}(),
				},
			},
		},
		"options": map[string]interface{}{
			"scales": map[string]interface{}{
				"y": map[string]interface{}{
					"beginAtZero": true,
				},
			},
			"plugins": map[string]interface{}{
				"title": map[string]interface{}{
					"display": true,
					"text":    "Commits Leaderboard",
				},
			},
		},
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    chartData,
	})
}

// GetCommitsTimeline returns commits timeline chart data
func (h *Handler) GetCommitsTimeline(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	params, err := h.parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	since, until, err := h.parseDateRange(params.StartDate, params.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid date format, use YYYY-MM-DD",
		})
		return
	}

	metrics, err := h.metricsService.GetCommitMetrics(c.Request.Context(), owner, repo, since, until)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Format for chart.js
	chartData := map[string]interface{}{
		"type": "line",
		"data": map[string]interface{}{
			"labels": func() []string {
				var labels []string
				for _, day := range metrics.Timeline {
					labels = append(labels, day.Date)
				}
				return labels
			}(),
			"datasets": []map[string]interface{}{
				{
					"label":           "Commits",
					"backgroundColor": "rgba(75, 192, 192, 0.2)",
					"borderColor":     "rgba(75, 192, 192, 1)",
					"borderWidth":     2,
					"fill":            false,
					"data": func() []int {
						var data []int
						for _, day := range metrics.Timeline {
							data = append(data, day.Commits)
						}
						return data
					}(),
				},
			},
		},
		"options": map[string]interface{}{
			"scales": map[string]interface{}{
				"y": map[string]interface{}{
					"beginAtZero": true,
				},
			},
			"plugins": map[string]interface{}{
				"title": map[string]interface{}{
					"display": true,
					"text":    "Commits Timeline",
				},
			},
		},
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    chartData,
	})
}

// GetPRsTimeline returns PRs timeline chart data
func (h *Handler) GetPRsTimeline(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	params, err := h.parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	since, until, err := h.parseDateRange(params.StartDate, params.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid date format, use YYYY-MM-DD",
		})
		return
	}

	metrics, err := h.metricsService.GetPRMetrics(c.Request.Context(), owner, repo, since, until)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Format for chart.js
	chartData := map[string]interface{}{
		"type": "line",
		"data": map[string]interface{}{
			"labels": func() []string {
				var labels []string
				for _, day := range metrics.Timeline {
					labels = append(labels, day.Date)
				}
				return labels
			}(),
			"datasets": []map[string]interface{}{
				{
					"label":           "Created",
					"backgroundColor": "rgba(255, 99, 132, 0.2)",
					"borderColor":     "rgba(255, 99, 132, 1)",
					"borderWidth":     2,
					"fill":            false,
					"data": func() []int {
						var data []int
						for _, day := range metrics.Timeline {
							data = append(data, day.Created)
						}
						return data
					}(),
				},
				{
					"label":           "Merged",
					"backgroundColor": "rgba(54, 162, 235, 0.2)",
					"borderColor":     "rgba(54, 162, 235, 1)",
					"borderWidth":     2,
					"fill":            false,
					"data": func() []int {
						var data []int
						for _, day := range metrics.Timeline {
							data = append(data, day.Merged)
						}
						return data
					}(),
				},
			},
		},
		"options": map[string]interface{}{
			"scales": map[string]interface{}{
				"y": map[string]interface{}{
					"beginAtZero": true,
				},
			},
			"plugins": map[string]interface{}{
				"title": map[string]interface{}{
					"display": true,
					"text":    "Pull Requests Timeline",
				},
			},
		},
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    chartData,
	})
}

// GetPRsLeaderboard returns PRs leaderboard chart data
func (h *Handler) GetPRsLeaderboard(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	params, err := h.parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	since, until, err := h.parseDateRange(params.StartDate, params.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid date format, use YYYY-MM-DD",
		})
		return
	}

	metrics, err := h.metricsService.GetPRMetrics(c.Request.Context(), owner, repo, since, until)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Format for chart.js
	chartData := map[string]interface{}{
		"type": "bar",
		"data": map[string]interface{}{
			"labels": func() []string {
				var labels []string
				for _, author := range metrics.Leaderboard {
					labels = append(labels, author.Author)
				}
				return labels
			}(),
			"datasets": []map[string]interface{}{
				{
					"label":           "Total PRs",
					"backgroundColor": "rgba(255, 206, 86, 0.6)",
					"borderColor":     "rgba(255, 206, 86, 1)",
					"borderWidth":     1,
					"data": func() []int {
						var data []int
						for _, author := range metrics.Leaderboard {
							data = append(data, author.TotalPRs)
						}
						return data
					}(),
				},
			},
		},
		"options": map[string]interface{}{
			"scales": map[string]interface{}{
				"y": map[string]interface{}{
					"beginAtZero": true,
				},
			},
			"plugins": map[string]interface{}{
				"title": map[string]interface{}{
					"display": true,
					"text":    "Pull Requests Leaderboard",
				},
			},
		},
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    chartData,
	})
}

// GetContributionsLeaderboard returns contributions leaderboard chart data
func (h *Handler) GetContributionsLeaderboard(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	params, err := h.parseQueryParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   err.Error(),
		})
		return
	}

	// Handle multiple repositories
	var repos []string
	if params.Repos != "" {
		repos = strings.Split(params.Repos, ",")
	} else {
		repos = []string{repo}
	}

	since, until, err := h.parseDateRange(params.StartDate, params.EndDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Error:   "invalid date format, use YYYY-MM-DD",
		})
		return
	}

	// Aggregate metrics from multiple repos
	allContributors := make(map[string]*models.Contributor)

	for _, repoName := range repos {
		metrics, err := h.metricsService.GetContributionMetrics(c.Request.Context(), owner, strings.TrimSpace(repoName), since, until)
		if err != nil {
			continue // Skip repos with errors but don't fail the entire request
		}

		// Merge contributors
		for _, contributor := range metrics.Contributors {
			if existing, exists := allContributors[contributor.Username]; exists {
				// Merge the data
				existing.TotalScore += contributor.TotalScore
				existing.Breakdown.Commits.Count += contributor.Breakdown.Commits.Count
				existing.Breakdown.Commits.Score += contributor.Breakdown.Commits.Score
				existing.Breakdown.PRReviews.Count += contributor.Breakdown.PRReviews.Count
				existing.Breakdown.PRReviews.Score += contributor.Breakdown.PRReviews.Score
				existing.Breakdown.Comments.Count += contributor.Breakdown.Comments.Count
				existing.Breakdown.Comments.Score += contributor.Breakdown.Comments.Score
				existing.Breakdown.Issues.Count += contributor.Breakdown.Issues.Count
				existing.Breakdown.Issues.Score += contributor.Breakdown.Issues.Score
				existing.Breakdown.Approvals.Count += contributor.Breakdown.Approvals.Count
				existing.Breakdown.Approvals.Score += contributor.Breakdown.Approvals.Score
				existing.Breakdown.Changes.Count += contributor.Breakdown.Changes.Count
				existing.Breakdown.Changes.Score += contributor.Breakdown.Changes.Score
			} else {
				// Create a copy
				allContributors[contributor.Username] = &models.Contributor{
					Username:   contributor.Username,
					TotalScore: contributor.TotalScore,
					Breakdown:  contributor.Breakdown,
				}
			}
		}
	}

	// Convert to slice and sort
	var contributors []models.Contributor
	for _, contributor := range allContributors {
		contributors = append(contributors, *contributor)
	}

	// Sort by total score (descending)
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].TotalScore > contributors[j].TotalScore
	})

	// Take top 20 for chart readability
	if len(contributors) > 20 {
		contributors = contributors[:20]
	}

	// Format for chart.js
	chartData := map[string]interface{}{
		"type": "bar",
		"data": map[string]interface{}{
			"labels": func() []string {
				var labels []string
				for _, contributor := range contributors {
					labels = append(labels, contributor.Username)
				}
				return labels
			}(),
			"datasets": []map[string]interface{}{
				{
					"label":           "Contribution Score",
					"backgroundColor": "rgba(153, 102, 255, 0.6)",
					"borderColor":     "rgba(153, 102, 255, 1)",
					"borderWidth":     1,
					"data": func() []float64 {
						var data []float64
						for _, contributor := range contributors {
							data = append(data, contributor.TotalScore)
						}
						return data
					}(),
				},
			},
		},
		"options": map[string]interface{}{
			"indexAxis": "y", // Horizontal bar chart for better username readability
			"scales": map[string]interface{}{
				"x": map[string]interface{}{
					"beginAtZero": true,
				},
			},
			"plugins": map[string]interface{}{
				"title": map[string]interface{}{
					"display": true,
					"text":    "Contributions Leaderboard",
				},
			},
		},
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Data:    chartData,
	})
}
