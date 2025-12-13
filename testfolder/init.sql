CREATE TABLE IF NOT EXISTS demo_items (
  id   bigserial PRIMARY KEY,
  name text NOT NULL UNIQUE,
  note text NOT NULL DEFAULT ''
);
-- // 差分づくり用
-- // 差分づくり用