package sections

import (
	"math"
	"testing"
	"time"

	gh "github.com/jeziellopes/livemark/internal/github"
)

func repo(stars, size int, pushedDaysAgo float64) gh.Repo {
	return gh.Repo{
		StargazersCount: stars,
		Size:            size,
		PushedAt:        time.Now().Add(-time.Duration(pushedDaysAgo*24) * time.Hour),
	}
}

func TestRepoScore_RecentBeatsStale(t *testing.T) {
	fresh := repo(10, 100, 1)
	stale := repo(10, 100, 700)
	if repoScore(fresh) <= repoScore(stale) {
		t.Errorf("fresh repo (1d) should score higher than stale (700d): %.2f vs %.2f",
			repoScore(fresh), repoScore(stale))
	}
}

func TestRepoScore_ZeroStarsFreshScoresPositive(t *testing.T) {
	r := repo(0, 100, 1)
	if repoScore(r) <= 0 {
		t.Errorf("zero-star repo pushed yesterday should score > 0, got %.4f", repoScore(r))
	}
}

func TestRepoScore_MoreStarsHigherScore(t *testing.T) {
	low := repo(5, 100, 7)
	high := repo(50, 100, 7)
	if repoScore(high) <= repoScore(low) {
		t.Errorf("higher-star repo should score higher: %.2f vs %.2f", repoScore(high), repoScore(low))
	}
}

func TestRepoScore_730DaysScoresNearZero(t *testing.T) {
	r := repo(100, 1000, 730)
	if repoScore(r) > 0.01 {
		t.Errorf("730-day-old repo should score ~0, got %.4f", repoScore(r))
	}
}

func TestRepoScore_MultiplicativeRecency(t *testing.T) {
	// Same repo: pushed today vs pushed 365 days ago.
	// 365 days ago → recencyFactor = (730-365)/730 ≈ 0.5
	// Score should be roughly half.
	today := repo(20, 500, 0)
	halfAge := repo(20, 500, 365)
	ratio := repoScore(halfAge) / repoScore(today)
	// Allow ±5% tolerance around 0.5
	if math.Abs(ratio-0.5) > 0.05 {
		t.Errorf("365-day-old repo should score ~50%% of today's; got ratio=%.3f", ratio)
	}
}

func TestRepoScore_FallsBackToUpdatedAt(t *testing.T) {
	// When PushedAt is zero, UpdatedAt should be used.
	r := gh.Repo{
		StargazersCount: 5,
		Size:            100,
		UpdatedAt:       time.Now().Add(-24 * time.Hour),
		// PushedAt is zero value
	}
	if repoScore(r) <= 0 {
		t.Error("should use UpdatedAt when PushedAt is zero")
	}
}
