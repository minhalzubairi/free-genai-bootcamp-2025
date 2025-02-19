CREATE TABLE words (
    id INTEGER PRIMARY KEY,
    arabic TEXT NOT NULL,
    roman TEXT NOT NULL,
    english TEXT NOT NULL,
    parts TEXT
);

CREATE TABLE groups (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE words_groups (
    id INTEGER PRIMARY KEY,
    word_id INTEGER NOT NULL,
    group_id INTEGER NOT NULL,
    FOREIGN KEY (word_id) REFERENCES words(id),
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE study_sessions (
    id INTEGER PRIMARY KEY,
    group_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    study_activity_id INTEGER NOT NULL,
    FOREIGN KEY (group_id) REFERENCES groups(id)
);

CREATE TABLE study_activities (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    thumbnail_url TEXT,
    description TEXT,
    launch_url TEXT
);

CREATE TABLE word_review_items (
    id INTEGER PRIMARY KEY,
    word_id INTEGER NOT NULL,
    study_session_id INTEGER NOT NULL,
    correct BOOLEAN NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (word_id) REFERENCES words(id),
    FOREIGN KEY (study_session_id) REFERENCES study_sessions(id)
); 
