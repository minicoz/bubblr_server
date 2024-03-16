DROP TABLE IF EXISTS matches;

CREATE TABLE matches (
    id SERIAL PRIMARY KEY,
    user_id UUID NOT NULL,
    matched_user_id UUID NOT NULL,
    liked_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    matched BOOLEAN DEFAULT FALSE
);

-- Indexes for faster lookups
CREATE INDEX idx_matches_user_a ON matches(user_id);
CREATE INDEX idx_matches_user_b ON matches(matched_user_id);

-- Unique constraint to prevent duplicate likes
ALTER TABLE matches ADD CONSTRAINT unique_like_pair UNIQUE (user_id, matched_user_id);


ALTER TABLE "matches" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("user_id");
ALTER TABLE "matches" ADD FOREIGN KEY ("matched_user_id") REFERENCES "users" ("user_id");