package models

import "time"

// Repository representa um repositório do GitHub
type Repository struct {
	Owner       string `json:"owner"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Private     bool   `json:"private"`
	Fork        bool   `json:"fork"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	Language    string `json:"language"`
	Stars       int    `json:"stars"`
	Forks       int    `json:"forks"`
}

// CommitMetrics representa métricas de commits
type CommitMetrics struct {
	TotalCommits  int            `json:"total_commits"`
	UniqueAuthors int            `json:"unique_authors"`
	DateRange     DateRange      `json:"date_range"`
	Leaderboard   []CommitAuthor `json:"leaderboard"`
	Timeline      []CommitDay    `json:"timeline"`
}

// CommitAuthor representa um autor de commits
type CommitAuthor struct {
	Author    string `json:"author"`
	Email     string `json:"email"`
	Commits   int    `json:"commits"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
}

// CommitDay representa commits em um dia específico
type CommitDay struct {
	Date    string `json:"date"`
	Commits int    `json:"commits"`
}

// PRMetrics representa métricas de Pull Requests
type PRMetrics struct {
	TotalPRs     int        `json:"total_prs"`
	OpenPRs      int        `json:"open_prs"`
	ClosedPRs    int        `json:"closed_prs"`
	MergedPRs    int        `json:"merged_prs"`
	DraftPRs     int        `json:"draft_prs"`
	DateRange    DateRange  `json:"date_range"`
	Leaderboard  []PRAuthor `json:"leaderboard"`
	Timeline     []PRDay    `json:"timeline"`
	AvgMergeTime float64    `json:"avg_merge_time_hours"`
	ApprovalRate float64    `json:"approval_rate"`
}

// PRAuthor representa um autor de PRs
type PRAuthor struct {
	Author       string  `json:"author"`
	TotalPRs     int     `json:"total_prs"`
	OpenPRs      int     `json:"open_prs"`
	MergedPRs    int     `json:"merged_prs"`
	ClosedPRs    int     `json:"closed_prs"`
	DraftPRs     int     `json:"draft_prs"`
	AvgMergeTime float64 `json:"avg_merge_time_hours"`
}

// PRDay representa PRs em um dia específico
type PRDay struct {
	Date    string `json:"date"`
	Created int    `json:"created"`
	Merged  int    `json:"merged"`
	Closed  int    `json:"closed"`
}

// ContributionMetrics representa métricas de contribuições
type ContributionMetrics struct {
	TotalContributions int           `json:"total_contributions"`
	TotalScore         float64       `json:"total_score"`
	DateRange          DateRange     `json:"date_range"`
	Contributors       []Contributor `json:"contributors"`
	ScoreWeights       ScoreWeights  `json:"score_weights"`
}

// Contributor representa um contribuidor
type Contributor struct {
	Username   string                `json:"username"`
	TotalScore float64               `json:"total_score"`
	Breakdown  ContributionBreakdown `json:"breakdown"`
}

// ContributionBreakdown detalha os tipos de contribuição
type ContributionBreakdown struct {
	Commits   ContributionDetail `json:"commits"`
	PRReviews ContributionDetail `json:"pr_reviews"`
	Comments  ContributionDetail `json:"comments"`
	Issues    ContributionDetail `json:"issues"`
	Approvals ContributionDetail `json:"approvals"`
	Changes   ContributionDetail `json:"changes"`
}

// ContributionDetail representa detalhes de um tipo de contribuição
type ContributionDetail struct {
	Count int     `json:"count"`
	Score float64 `json:"score"`
}

// ScoreWeights define os pesos para diferentes tipos de contribuição
type ScoreWeights struct {
	Commit   float64 `json:"commit"`
	PRReview float64 `json:"pr_review"`
	Comment  float64 `json:"comment"`
	Issue    float64 `json:"issue"`
	Approval float64 `json:"approval"`
	Changes  float64 `json:"changes"`
}

// DateRange representa um intervalo de datas
type DateRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// APIResponse representa uma resposta padrão da API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// QueryParams representa parâmetros de query comuns
type QueryParams struct {
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Status    string `form:"status"`
	Repos     string `form:"repos"`
}

// ParseDate converte string para time.Time
func ParseDate(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// FormatDate converte time.Time para string
func FormatDate(t time.Time) string {
	return t.Format("2006-01-02")
}
