-- CREATE USER shortlink;
-- CREATE DATABASE shortlink;
-- GRANT ALL PRIVILEGES ON DATABASE shortlink TO shortlink;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE links
(
    id         uuid DEFAULT uuid_generate_v4(),
    hash       VARCHAR NOT NULL,
    admin_link VARCHAR NOT NULL,
    link       VARCHAR NOT NULL,
    count      BIGINT DEFAULT 0,
    created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp NULL,
    deleted_at timestamp NULL,
    PRIMARY KEY (id)
);

CREATE UNIQUE INDEX idx_links_link ON links(link);

CREATE TABLE link_stats
(
  link_id uuid NOT NULL,
  ip varchar(30) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE link_stats ADD CONSTRAINT fk_links FOREIGN KEY (link_id) REFERENCES links (id) ON DELETE CASCADE;
