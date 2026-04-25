CREATE TABLE agent_history (
    id TEXT PRIMARY KEY,
    prompt TEXT NOT NULL,
    response TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);
