-- Drop foreign keys
ALTER TABLE "matches" DROP CONSTRAINT IF EXISTS matches_user_id_fkey;
ALTER TABLE "matches" DROP CONSTRAINT IF EXISTS matches_matched_user_id_fkey;

-- Drop unique constraint
ALTER TABLE "matches" DROP CONSTRAINT IF EXISTS unique_like_pair;

-- Drop indexes
DROP INDEX IF EXISTS idx_matches_user_a;
DROP INDEX IF EXISTS idx_matches_user_b;

-- Drop the table
DROP TABLE IF EXISTS matches;
