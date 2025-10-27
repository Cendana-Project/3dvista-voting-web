-- Migration: Change to one vote per IP globally
-- Instead of one vote per IP per innovation, now one IP can only vote once in total

-- First, remove duplicate votes (keep only the first vote per IP)
DELETE FROM votes
WHERE id NOT IN (
    SELECT MIN(id)
    FROM votes
    GROUP BY voter_ip_hash
);

-- Drop old constraint
ALTER TABLE votes DROP CONSTRAINT IF EXISTS votes_unique_per_ip_per_innovation;

-- Add new constraint: one vote per IP globally
ALTER TABLE votes ADD CONSTRAINT votes_unique_per_ip UNIQUE (voter_ip_hash);

-- Note: innovation_id column remains to track which innovation was voted for
-- But the uniqueness is now on voter_ip_hash alone
