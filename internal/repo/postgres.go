package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"voteweb/internal/domain"
)

type postgresRepository struct {
	pool *pgxpool.Pool
}

// NewPostgresRepository creates a new PostgreSQL repository
func NewPostgresRepository(pool *pgxpool.Pool) domain.Repository {
	return &postgresRepository{pool: pool}
}

func (r *postgresRepository) GetInnovationBySlug(ctx context.Context, groupSlug, slug string) (*domain.Innovation, error) {
	query := `
		SELECT id, group_slug, slug, name, division, entity_name, pic, description,
		       logo_innovation_url, logo_entity_url, video_url, slide_url, ig_url, yt_url,
		       created_at, updated_at
		FROM innovations
		WHERE group_slug = $1 AND slug = $2
	`

	var innovation domain.Innovation
	err := r.pool.QueryRow(ctx, query, groupSlug, slug).Scan(
		&innovation.ID,
		&innovation.GroupSlug,
		&innovation.Slug,
		&innovation.Name,
		&innovation.Division,
		&innovation.EntityName,
		&innovation.PIC,
		&innovation.Description,
		&innovation.LogoInnovationURL,
		&innovation.LogoEntityURL,
		&innovation.VideoURL,
		&innovation.SlideURL,
		&innovation.IgURL,
		&innovation.YtURL,
		&innovation.CreatedAt,
		&innovation.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInnovationNotFound
		}
		return nil, fmt.Errorf("query innovation: %w", err)
	}

	return &innovation, nil
}

func (r *postgresRepository) InsertVote(ctx context.Context, vote *domain.Vote) (bool, error) {
	query := `
		INSERT INTO votes (innovation_id, voter_ip_hash, user_agent, created_at)
		VALUES ($1, $2, $3, NOW())
		ON CONFLICT (voter_ip_hash) DO NOTHING
		RETURNING id
	`

	var id int64
	err := r.pool.QueryRow(ctx, query, vote.InnovationID, vote.VoterIPHash, vote.UserAgent).Scan(&id)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No rows returned means conflict occurred - this IP has already voted
			return false, nil
		}
		return false, fmt.Errorf("insert vote: %w", err)
	}

	// Row was inserted successfully
	return true, nil
}

func (r *postgresRepository) GetVoteCount(ctx context.Context, innovationID string) (int64, error) {
	query := `SELECT COUNT(*) FROM votes WHERE innovation_id = $1`

	var count int64
	err := r.pool.QueryRow(ctx, query, innovationID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count votes: %w", err)
	}

	return count, nil
}

func (r *postgresRepository) ListInnovations(ctx context.Context) ([]*domain.Innovation, error) {
	query := `
		SELECT id, group_slug, slug, name, division, entity_name, pic, description,
		       logo_innovation_url, logo_entity_url, video_url, slide_url, ig_url, yt_url,
		       created_at, updated_at
		FROM innovations
		ORDER BY group_slug, name
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query innovations: %w", err)
	}
	defer rows.Close()

	var innovations []*domain.Innovation
	for rows.Next() {
		var innovation domain.Innovation
		err := rows.Scan(
			&innovation.ID,
			&innovation.GroupSlug,
			&innovation.Slug,
			&innovation.Name,
			&innovation.Division,
			&innovation.EntityName,
			&innovation.PIC,
			&innovation.Description,
			&innovation.LogoInnovationURL,
			&innovation.LogoEntityURL,
			&innovation.VideoURL,
			&innovation.SlideURL,
			&innovation.IgURL,
			&innovation.YtURL,
			&innovation.CreatedAt,
			&innovation.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scan innovation: %w", err)
		}
		innovations = append(innovations, &innovation)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate rows: %w", err)
	}

	return innovations, nil
}

func (r *postgresRepository) HasVoted(ctx context.Context, innovationID string, voterIPHash []byte) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM votes 
			WHERE innovation_id = $1 AND voter_ip_hash = $2
		)
	`

	var exists bool
	err := r.pool.QueryRow(ctx, query, innovationID, voterIPHash).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check vote: %w", err)
	}

	return exists, nil
}

func (r *postgresRepository) GetTotalVoters(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(DISTINCT voter_ip_hash) FROM votes`

	var count int64
	err := r.pool.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count total voters: %w", err)
	}

	return count, nil
}

func (r *postgresRepository) HasVotedGlobally(ctx context.Context, voterIPHash []byte) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM votes WHERE voter_ip_hash = $1)`

	var exists bool
	err := r.pool.QueryRow(ctx, query, voterIPHash).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("check global vote: %w", err)
	}

	return exists, nil
}

func (r *postgresRepository) GetVotedInnovation(ctx context.Context, voterIPHash []byte) (*domain.Innovation, error) {
	query := `
		SELECT i.id, i.group_slug, i.slug, i.name, i.division, i.entity_name, i.pic, i.description,
		       i.logo_innovation_url, i.logo_entity_url, i.video_url, i.slide_url, i.ig_url, i.yt_url,
		       i.created_at, i.updated_at
		FROM votes v
		JOIN innovations i ON v.innovation_id = i.id
		WHERE v.voter_ip_hash = $1
		LIMIT 1
	`

	var innovation domain.Innovation
	err := r.pool.QueryRow(ctx, query, voterIPHash).Scan(
		&innovation.ID,
		&innovation.GroupSlug,
		&innovation.Slug,
		&innovation.Name,
		&innovation.Division,
		&innovation.EntityName,
		&innovation.PIC,
		&innovation.Description,
		&innovation.LogoInnovationURL,
		&innovation.LogoEntityURL,
		&innovation.VideoURL,
		&innovation.SlideURL,
		&innovation.IgURL,
		&innovation.YtURL,
		&innovation.CreatedAt,
		&innovation.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrInnovationNotFound
		}
		return nil, fmt.Errorf("get voted innovation: %w", err)
	}

	return &innovation, nil
}
