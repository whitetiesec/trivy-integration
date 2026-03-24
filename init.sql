CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    owner VARCHAR(100) NOT NULL,
    size BIGINT DEFAULT 0,
    mime_type VARCHAR(100) DEFAULT 'application/pdf',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO documents (name, owner, size, mime_type) VALUES
    ('report-q1.pdf', 'alice', 102400, 'application/pdf'),
    ('invoice-2024.pdf', 'bob', 51200, 'application/pdf'),
    ('contract-draft.pdf', 'charlie', 204800, 'application/pdf');
