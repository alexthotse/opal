-- name: GetHistory :many
SELECT * FROM agent_history ORDER BY created_at DESC LIMIT 100;

-- name: CreateHistory :exec
INSERT INTO agent_history (id, prompt, response) VALUES (?, ?, ?);
