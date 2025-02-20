#!/bin/bash

# Remove existing database
rm -f words.db

# Create new database
sqlite3 words.db < db/migrations/0001_init.sql

# Insert sample data
sqlite3 words.db << 'EOF'
INSERT INTO groups (name) VALUES ('Basic Greetings');
INSERT INTO groups (name) VALUES ('Numbers');
INSERT INTO groups (name) VALUES ('Colors');

INSERT INTO words (japanese, romaji, english, parts) VALUES 
('こんにちは', 'konnichiwa', 'hello', '{"type": "greeting", "formality": "neutral"}'),
('さようなら', 'sayounara', 'goodbye', '{"type": "greeting", "formality": "formal"}'),
('おはよう', 'ohayou', 'good morning', '{"type": "greeting", "formality": "informal"}');

INSERT INTO word_groups (word_id, group_id) VALUES (1, 1), (2, 1), (3, 1);

INSERT INTO study_activities (name, thumbnail_url, description) VALUES 
('Vocabulary Quiz', 'https://example.com/thumbnail.jpg', 'Practice your vocabulary with flashcards');

INSERT INTO study_sessions (group_id, study_activity_id, start_time, end_time) VALUES 
(1, 1, datetime('now', '-1 day'), datetime('now', '-1 day', '+30 minutes')),
(1, 1, datetime('now'), NULL);

INSERT INTO word_review_items (word_id, study_session_id, correct, created_at) VALUES 
(1, 1, 1, datetime('now', '-1 day')),
(2, 1, 1, datetime('now', '-1 day')),
(3, 1, 0, datetime('now', '-1 day'));
EOF
