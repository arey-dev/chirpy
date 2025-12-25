-- +goose Up
CREATE TABLE chirps (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL,
  body VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

  CONSTRAINT fk_chirps_users
  FOREIGN KEY (user_id)
  REFERENCES users(id)
  ON DELETE CASCADE
);

-- +goose Down
DROP TABLE chirps;
