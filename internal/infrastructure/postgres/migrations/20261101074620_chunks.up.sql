CREATE TABLE IF NOT EXISTS chunks (
    id        BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    content   TEXT NOT NULL,
    source    TEXT NOT NULL,
    embedding vector(1024) NOT NULL,
    chunks_index INT,
    heading TEXT
);

CREATE INDEX IF NOT EXISTS chunks_embedding_idx
    ON chunks USING hnsw (embedding vector_cosine_ops);