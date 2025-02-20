-- Insert sample groups
INSERT INTO groups (name) VALUES 
('Basic Greetings'),
('Numbers'),
('Colors');

-- Insert sample words
INSERT INTO words (japanese, romaji, english, parts) VALUES 
('こんにちは', 'konnichiwa', 'hello', '{"type": "greeting", "formality": "neutral"}'),
('さようなら', 'sayounara', 'goodbye', '{"type": "greeting", "formality": "formal"}'),
('おはよう', 'ohayou', 'good morning', '{"type": "greeting", "formality": "informal"}'),
('一', 'ichi', 'one', '{"type": "number", "category": "cardinal"}'),
('二', 'ni', 'two', '{"type": "number", "category": "cardinal"}'),
('赤', 'aka', 'red', '{"type": "color", "category": "basic"}'),
('青', 'ao', 'blue', '{"type": "color", "category": "basic"}');

-- Link words to groups
INSERT INTO word_groups (word_id, group_id) VALUES 
(1, 1), -- konnichiwa -> Basic Greetings
(2, 1), -- sayounara -> Basic Greetings
(3, 1), -- ohayou -> Basic Greetings
(4, 2), -- ichi -> Numbers
(5, 2), -- ni -> Numbers
(6, 3), -- aka -> Colors
(7, 3); -- ao -> Colors

-- Create some study activities
INSERT INTO study_activities (group_id, study_session_id, created_at) VALUES 
(1, 1, datetime('now', '-5 days')),
(2, 2, datetime('now', '-3 days')),
(3, 3, datetime('now', '-1 day'));

-- Create study sessions
INSERT INTO study_sessions (group_id, study_session_id, correct, study_activity_id, created_at) VALUES 
(1, 1, 1, 1, datetime('now', '-5 days')),
(2, 2, 1, 2, datetime('now', '-3 days')),
(3, 3, 1, 3, datetime('now', '-1 day'));

-- Add word review items
INSERT INTO word_reviews_items (word_id, study_session_id, correct, created_at) VALUES 
(1, 1, 1, datetime('now', '-5 days')), -- correct
(2, 1, 0, datetime('now', '-5 days')), -- incorrect
(3, 1, 1, datetime('now', '-5 days')), -- correct
(4, 2, 1, datetime('now', '-3 days')), -- correct
(5, 2, 1, datetime('now', '-3 days')), -- correct
(6, 3, 0, datetime('now', '-1 day')),  -- incorrect
(7, 3, 1, datetime('now', '-1 day'));  -- correct
