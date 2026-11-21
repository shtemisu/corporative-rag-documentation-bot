-- Чистим дубликаты, накопившиеся до идемпотентного ингеста.
DELETE FROM chunks c
USING chunks d
WHERE c.id > d.id
  AND c.source = d.source
  AND c.chunks_index = d.chunks_index;

-- Защита от повторной вставки одного и того же чанка документа.
ALTER TABLE chunks
    ADD CONSTRAINT chunks_source_index_key UNIQUE (source, chunks_index);
