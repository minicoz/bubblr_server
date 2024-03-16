CREATE TABLE IF NOT EXISTS votes (
    id SERIAL PRIMARY KEY, 
    prospective_user_id UUID,
    user_id UUID,
    vote_yes BOOLEAN Not NULL,
    created_at timestamptz NOT NULL DEFAULT(now()),
    CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(user_id)
);
