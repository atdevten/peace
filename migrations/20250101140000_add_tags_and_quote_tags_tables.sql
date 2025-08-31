-- +goose Up
-- Create tags table
CREATE TABLE IF NOT EXISTS tags (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Create quote_tags junction table for many-to-many relationship
CREATE TABLE IF NOT EXISTS quote_tags (
    id SERIAL PRIMARY KEY,
    quote_id INTEGER NOT NULL REFERENCES quotes(id) ON DELETE CASCADE,
    tag_id INTEGER NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(quote_id, tag_id)
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_tags_name ON tags(name);
CREATE INDEX IF NOT EXISTS idx_tags_deleted_at ON tags(deleted_at) WHERE deleted_at IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_quote_tags_quote_id ON quote_tags(quote_id);
CREATE INDEX IF NOT EXISTS idx_quote_tags_tag_id ON quote_tags(tag_id);

-- Add comments
COMMENT ON TABLE tags IS 'Stores tags for categorizing quotes';
COMMENT ON COLUMN tags.id IS 'Unique auto-increment identifier for the tag';
COMMENT ON COLUMN tags.name IS 'Tag name (unique)';
COMMENT ON COLUMN tags.description IS 'Optional description of the tag';
COMMENT ON COLUMN tags.created_at IS 'When the tag was created';
COMMENT ON COLUMN tags.updated_at IS 'When the tag was last updated';
COMMENT ON COLUMN tags.deleted_at IS 'Soft delete timestamp';

COMMENT ON TABLE quote_tags IS 'Junction table linking quotes and tags (many-to-many)';
COMMENT ON COLUMN quote_tags.id IS 'Unique auto-increment identifier for the quote-tag relationship';
COMMENT ON COLUMN quote_tags.quote_id IS 'Reference to quotes table';
COMMENT ON COLUMN quote_tags.tag_id IS 'Reference to tags table';
COMMENT ON COLUMN quote_tags.created_at IS 'When the quote-tag relationship was created';

-- +goose Down
-- Drop indexes
DROP INDEX IF EXISTS idx_quote_tags_tag_id;
DROP INDEX IF EXISTS idx_quote_tags_quote_id;
DROP INDEX IF EXISTS idx_tags_deleted_at;
DROP INDEX IF EXISTS idx_tags_name;

-- Drop tables
DROP TABLE IF EXISTS quote_tags;
DROP TABLE IF EXISTS tags;
