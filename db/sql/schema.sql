CREATE TABLE IF NOT EXISTS executions (
  id INTEGER PRIMARY KEY,
  date DATE NOT NULL DEFAULT CURRENT_DATE
);
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY,
  uuid VARCHAR (255) NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS matches (
  id INTEGER PRIMARY KEY,
  a INTEGER NOT NULL REFERENCES users (id),
  b INTEGER NOT NULL REFERENCES users (id),
  freq INTEGER NOT NULL DEFAULT 0,
  CHECK (a < b)
);
CREATE UNIQUE INDEX IF NOT EXISTS matches_a_b_uidx ON matches (a, b);