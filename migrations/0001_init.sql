-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Innovations table
CREATE TABLE IF NOT EXISTS innovations (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  group_slug TEXT NOT NULL,
  slug TEXT NOT NULL,
  name TEXT NOT NULL,
  division TEXT,
  entity_name TEXT,
  pic TEXT,
  description TEXT,
  logo_innovation_url TEXT,
  logo_entity_url TEXT,
  video_url TEXT,
  slide_url TEXT,
  ig_url TEXT,
  yt_url TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT innovations_group_slug_slug_uk UNIQUE (group_slug, slug)
);

-- Votes table
CREATE TABLE IF NOT EXISTS votes (
  id BIGSERIAL PRIMARY KEY,
  innovation_id UUID NOT NULL REFERENCES innovations(id) ON DELETE CASCADE,
  voter_ip_hash BYTEA NOT NULL,
  user_agent TEXT,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  CONSTRAINT votes_unique_per_ip_per_innovation UNIQUE (innovation_id, voter_ip_hash)
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_votes_innovation ON votes(innovation_id);
CREATE INDEX IF NOT EXISTS idx_innovations_group_slug ON innovations(group_slug);


