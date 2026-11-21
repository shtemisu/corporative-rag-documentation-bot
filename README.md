# corporative-rag-documentation-bot

This bot need for easily start into project

for start Ollama service:
```bash
sudo systemctl start ollama
```

for stop Ollama service:
```bash
sudo systemctl stop ollama
```

## Порог релевантности

В контекст модели попадают только чанки, косинусное расстояние которых
(`embedding <=> query`) не превышает порог `DefaultMaxDistance = 0.5`
(см. `internal/usecase/rag.go`). Порог передаётся последним аргументом
`NewRAGService`; `0` означает «использовать значение по умолчанию».

Порог зависит от модели эмбеддингов и корпуса, поэтому его нужно калибровать:

1. Возьмите 10–20 вопросов, на которые документация точно отвечает, и столько же вопросов не по теме.
2. Прогоните их через `RAGService.Answer`.
3. Смотрите логи:
   `rag: context built, best_distance=... kept=... retrieved=...` —
   вопрос прошёл порог;
   `rag: question rejected, best_distance=... max_distance=...` — не прошёл.
4. Если у релевантных вопросов `best_distance` стабильно выше порога — увеличьте порог.
   Если у нерелевантных `best_distance` ниже порога — уменьшите.
5. Ориентир для `bge-m3` на русскоязычной документации: 0.4–0.6.
   Абсолютные значения несравнимы между разными моделями эмбеддингов.