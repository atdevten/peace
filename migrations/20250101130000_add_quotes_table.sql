-- +goose Up
-- Create quotes table
CREATE TABLE IF NOT EXISTS quotes (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    author TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes for quotes table
CREATE INDEX IF NOT EXISTS idx_quotes_author ON quotes(author);
CREATE INDEX IF NOT EXISTS idx_quotes_created_at ON quotes(created_at);
CREATE INDEX IF NOT EXISTS idx_quotes_deleted_at ON quotes(deleted_at) WHERE deleted_at IS NOT NULL;

-- Add comments
COMMENT ON TABLE quotes IS 'Stores inspirational and motivational quotes';
COMMENT ON COLUMN quotes.id IS 'Unique auto-increment identifier for the quote';
COMMENT ON COLUMN quotes.content IS 'The quote text content';
COMMENT ON COLUMN quotes.author IS 'Author of the quote (optional)';
COMMENT ON COLUMN quotes.created_at IS 'When the quote was added to the system';
COMMENT ON COLUMN quotes.updated_at IS 'When the quote was last updated';
COMMENT ON COLUMN quotes.deleted_at IS 'Soft delete timestamp';

-- +goose Down
-- Drop indexes
DROP INDEX IF EXISTS idx_quotes_author;
DROP INDEX IF EXISTS idx_quotes_created_at;
DROP INDEX IF EXISTS idx_quotes_deleted_at;

-- Drop table
DROP TABLE IF EXISTS quotes; 