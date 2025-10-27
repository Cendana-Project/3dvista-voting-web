package domain

import (
	"context"
	"fmt"
	"log/slog"
)

// VoteService handles the business logic for voting
type VoteService interface {
	GetInnovation(ctx context.Context, groupSlug, slug string) (*Innovation, error)
	SubmitVote(ctx context.Context, req VoteRequest) (*VoteResponse, error)
	GetVoteCount(ctx context.Context, innovationID string) (int64, error)
	ListInnovations(ctx context.Context) ([]*Innovation, error)
	CheckHasVoted(ctx context.Context, innovationID, clientIP string) (bool, error)
	GetTotalVoters(ctx context.Context) (int64, error)
}

type voteService struct {
	repo   Repository
	hasher IPHasher
	logger *slog.Logger
}

// NewVoteService creates a new VoteService
func NewVoteService(repo Repository, hasher IPHasher, logger *slog.Logger) VoteService {
	return &voteService{
		repo:   repo,
		hasher: hasher,
		logger: logger,
	}
}

func (s *voteService) GetInnovation(ctx context.Context, groupSlug, slug string) (*Innovation, error) {
	innovation, err := s.repo.GetInnovationBySlug(ctx, groupSlug, slug)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get innovation",
			"group_slug", groupSlug,
			"slug", slug,
			"error", err)
		return nil, err
	}
	return innovation, nil
}

func (s *voteService) SubmitVote(ctx context.Context, req VoteRequest) (*VoteResponse, error) {
	// Hash IP
	ipHash := s.hasher.HashIP(req.ClientIP)

	// Check if IP has already voted (globally - one vote per IP)
	hasVoted, err := s.repo.HasVotedGlobally(ctx, ipHash)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to check global vote status",
			"error", err)
		return nil, fmt.Errorf("failed to check vote status: %w", err)
	}

	if hasVoted {
		// Get the innovation they already voted for
		votedInnovation, err := s.repo.GetVotedInnovation(ctx, ipHash)
		if err != nil {
			s.logger.ErrorContext(ctx, "failed to get voted innovation",
				"error", err)
			return &VoteResponse{
				Success:      false,
				AlreadyVoted: true,
				VoteCount:    0,
				Message:      "Anda sudah pernah vote untuk inovasi lain. Hanya 1 vote per IP.",
			}, nil
		}

		s.logger.InfoContext(ctx, "duplicate vote attempt - already voted globally",
			"group_slug", req.GroupSlug,
			"slug", req.Slug)
		return &VoteResponse{
			Success:      false,
			AlreadyVoted: true,
			VoteCount:    0,
			Message:      fmt.Sprintf("Anda sudah pernah vote untuk '%s'. Hanya 1 vote per IP yang diizinkan.", votedInnovation.Name),
		}, nil
	}

	// Get innovation to vote for
	innovation, err := s.repo.GetInnovationBySlug(ctx, req.GroupSlug, req.Slug)
	if err != nil {
		s.logger.ErrorContext(ctx, "innovation not found",
			"group_slug", req.GroupSlug,
			"slug", req.Slug,
			"error", err)
		return nil, err
	}

	// Insert vote
	vote := &Vote{
		InnovationID: innovation.ID,
		VoterIPHash:  ipHash,
		UserAgent:    req.UserAgent,
	}

	inserted, err := s.repo.InsertVote(ctx, vote)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to insert vote",
			"innovation_id", innovation.ID,
			"error", err)
		return nil, fmt.Errorf("failed to insert vote: %w", err)
	}

	if !inserted {
		return &VoteResponse{
			Success:      false,
			AlreadyVoted: true,
			VoteCount:    0,
			Message:      "Vote gagal diproses",
		}, nil
	}

	// Get current vote count for this innovation
	count, err := s.repo.GetVoteCount(ctx, innovation.ID)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to get vote count",
			"innovation_id", innovation.ID,
			"error", err)
		return nil, fmt.Errorf("failed to get vote count: %w", err)
	}

	s.logger.InfoContext(ctx, "vote recorded successfully",
		"innovation_id", innovation.ID,
		"group_slug", req.GroupSlug,
		"slug", req.Slug,
		"vote_count", count)

	return &VoteResponse{
		Success:      true,
		AlreadyVoted: false,
		VoteCount:    count,
		Message:      "Vote berhasil dicatat",
	}, nil
}

func (s *voteService) GetVoteCount(ctx context.Context, innovationID string) (int64, error) {
	return s.repo.GetVoteCount(ctx, innovationID)
}

func (s *voteService) ListInnovations(ctx context.Context) ([]*Innovation, error) {
	return s.repo.ListInnovations(ctx)
}

func (s *voteService) CheckHasVoted(ctx context.Context, innovationID, clientIP string) (bool, error) {
	ipHash := s.hasher.HashIP(clientIP)
	return s.repo.HasVotedGlobally(ctx, ipHash)
}

func (s *voteService) GetTotalVoters(ctx context.Context) (int64, error) {
	return s.repo.GetTotalVoters(ctx)
}

// Repository defines the data access interface
type Repository interface {
	GetInnovationBySlug(ctx context.Context, groupSlug, slug string) (*Innovation, error)
	InsertVote(ctx context.Context, vote *Vote) (bool, error)
	GetVoteCount(ctx context.Context, innovationID string) (int64, error)
	ListInnovations(ctx context.Context) ([]*Innovation, error)
	HasVoted(ctx context.Context, innovationID string, voterIPHash []byte) (bool, error)
	GetTotalVoters(ctx context.Context) (int64, error)
	HasVotedGlobally(ctx context.Context, voterIPHash []byte) (bool, error)
	GetVotedInnovation(ctx context.Context, voterIPHash []byte) (*Innovation, error)
}

// IPHasher defines the interface for IP hashing
type IPHasher interface {
	HashIP(ip string) []byte
}
