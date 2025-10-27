package domain

import (
	"context"
	"log/slog"
	"os"
	"testing"
)

// Mock repository for testing
type mockRepository struct {
	innovations map[string]*Innovation
	votes       map[string]bool // key: innovationID + ipHash
	voteCounts  map[string]int64
}

func newMockRepository() *mockRepository {
	return &mockRepository{
		innovations: make(map[string]*Innovation),
		votes:       make(map[string]bool),
		voteCounts:  make(map[string]int64),
	}
}

func (m *mockRepository) GetInnovationBySlug(ctx context.Context, groupSlug, slug string) (*Innovation, error) {
	key := groupSlug + ":" + slug
	if innovation, ok := m.innovations[key]; ok {
		return innovation, nil
	}
	return nil, ErrInnovationNotFound
}

func (m *mockRepository) InsertVote(ctx context.Context, vote *Vote) (bool, error) {
	key := vote.InnovationID + ":" + string(vote.VoterIPHash)
	if m.votes[key] {
		return false, nil // Already voted
	}
	m.votes[key] = true
	m.voteCounts[vote.InnovationID]++
	return true, nil
}

func (m *mockRepository) GetVoteCount(ctx context.Context, innovationID string) (int64, error) {
	return m.voteCounts[innovationID], nil
}

// Mock IP hasher
type mockIPHasher struct{}

func (m *mockIPHasher) HashIP(ip string) []byte {
	return []byte(ip) // Simple mock - just use IP as hash
}

func TestVoteService_SubmitVote(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	repo := newMockRepository()
	hasher := &mockIPHasher{}
	service := NewVoteService(repo, hasher, logger)

	// Add test innovation
	innovation := &Innovation{
		ID:        "test-id-1",
		GroupSlug: "test-group",
		Slug:      "test-innovation",
		Name:      "Test Innovation",
	}
	repo.innovations["test-group:test-innovation"] = innovation

	ctx := context.Background()

	t.Run("first vote succeeds", func(t *testing.T) {
		req := VoteRequest{
			GroupSlug: "test-group",
			Slug:      "test-innovation",
			ClientIP:  "192.168.1.1",
			UserAgent: "Test Agent",
		}

		result, err := service.SubmitVote(ctx, req)
		if err != nil {
			t.Fatalf("SubmitVote() error = %v", err)
		}

		if !result.Success {
			t.Error("Expected vote to succeed")
		}

		if result.AlreadyVoted {
			t.Error("Expected AlreadyVoted to be false")
		}

		if result.VoteCount != 1 {
			t.Errorf("Expected VoteCount = 1, got %d", result.VoteCount)
		}
	})

	t.Run("duplicate vote from same IP", func(t *testing.T) {
		req := VoteRequest{
			GroupSlug: "test-group",
			Slug:      "test-innovation",
			ClientIP:  "192.168.1.1",
			UserAgent: "Test Agent",
		}

		result, err := service.SubmitVote(ctx, req)
		if err != nil {
			t.Fatalf("SubmitVote() error = %v", err)
		}

		if result.Success {
			t.Error("Expected vote to fail (already voted)")
		}

		if !result.AlreadyVoted {
			t.Error("Expected AlreadyVoted to be true")
		}

		if result.VoteCount != 1 {
			t.Errorf("Expected VoteCount to remain 1, got %d", result.VoteCount)
		}
	})

	t.Run("vote from different IP succeeds", func(t *testing.T) {
		req := VoteRequest{
			GroupSlug: "test-group",
			Slug:      "test-innovation",
			ClientIP:  "192.168.1.2",
			UserAgent: "Test Agent",
		}

		result, err := service.SubmitVote(ctx, req)
		if err != nil {
			t.Fatalf("SubmitVote() error = %v", err)
		}

		if !result.Success {
			t.Error("Expected vote to succeed")
		}

		if result.VoteCount != 2 {
			t.Errorf("Expected VoteCount = 2, got %d", result.VoteCount)
		}
	})

	t.Run("innovation not found", func(t *testing.T) {
		req := VoteRequest{
			GroupSlug: "nonexistent",
			Slug:      "nonexistent",
			ClientIP:  "192.168.1.1",
			UserAgent: "Test Agent",
		}

		_, err := service.SubmitVote(ctx, req)
		if err != ErrInnovationNotFound {
			t.Errorf("Expected ErrInnovationNotFound, got %v", err)
		}
	})
}

