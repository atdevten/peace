-- +goose Up
-- Create mental_health_records table
CREATE TABLE IF NOT EXISTS mental_health_records (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    happy_level INTEGER NOT NULL,
    energy_level INTEGER NOT NULL,
    notes TEXT,
    status char(25) NOT NULL DEFAULT 'CLOSED',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes for mental_health_records table
CREATE INDEX IF NOT EXISTS idx_mental_health_records_user_id ON mental_health_records(user_id);

-- Add foreign key constraint to mental_health_records table
ALTER TABLE mental_health_records ADD CONSTRAINT fk_mental_health_records_user_id 
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

-- Add comments
COMMENT ON TABLE mental_health_records IS 'Stores mental health level records for users';
COMMENT ON COLUMN mental_health_records.id IS 'Unique identifier for the mental health record';
COMMENT ON COLUMN mental_health_records.user_id IS 'Reference to the user who recorded this mental health level';
COMMENT ON COLUMN mental_health_records.happy_level IS 'Overall happiness level (1-10 scale)';
COMMENT ON COLUMN mental_health_records.energy_level IS 'Energy level (1-10 scale)';
COMMENT ON COLUMN mental_health_records.notes IS 'Additional notes about mental health state';
COMMENT ON COLUMN mental_health_records.status IS 'Status of the mental health record';

-- +goose Down
-- Drop foreign key constraint
ALTER TABLE mental_health_records DROP CONSTRAINT IF EXISTS fk_mental_health_records_user_id;

-- Drop indexes
DROP INDEX IF EXISTS idx_mental_health_records_user_id;

-- Drop table
DROP TABLE IF EXISTS mental_health_records; 