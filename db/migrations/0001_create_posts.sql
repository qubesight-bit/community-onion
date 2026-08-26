BEGIN;

CREATE TABLE posts (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    title varchar(160) NOT NULL,
    body text NOT NULL,
    status varchar(16) NOT NULL DEFAULT 'draft',
    created_at timestamptz NOT NULL DEFAULT now(),
    published_at timestamptz,
    CONSTRAINT posts_title_not_blank
        CHECK (length(btrim(title)) BETWEEN 1 AND 160),
    CONSTRAINT posts_body_not_blank
        CHECK (length(btrim(body)) BETWEEN 1 AND 50000),
    CONSTRAINT posts_status_valid
        CHECK (status IN ('draft', 'published', 'archived')),
    CONSTRAINT posts_publication_date_valid
        CHECK (
            (status = 'published' AND published_at IS NOT NULL)
            OR
            (status <> 'published')
        )
);

CREATE INDEX posts_published_at_idx
    ON posts (published_at DESC)
    WHERE status = 'published';

CREATE VIEW published_posts
WITH (security_barrier = true)
AS
SELECT id, title, body, published_at
FROM posts
WHERE status = 'published'
  AND published_at <= now();

REVOKE ALL ON posts FROM community_app;
GRANT SELECT ON published_posts TO community_app;
REVOKE CREATE ON SCHEMA public FROM community_app;

COMMIT;
