package usecase

import (
	"context"
	"corp-rag-tg-bot/internal/domain"
	"strings"
	"testing"
)

type fakeEmbedder struct {
	vectors [][]float32
	err     error
}

func (f *fakeEmbedder) Embed(_ context.Context, text []string) ([][]float32, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.vectors != nil {
		return f.vectors, nil
	}
	vectors := make([][]float32, len(text))
	for i := range vectors {
		vectors[i] = []float32{1, 0}
	}
	return vectors, nil
}

type fakeSearcher struct {
	chunks []domain.Chunk
	err    error
	topK   int
}

func (f *fakeSearcher) Search(_ context.Context, _ []float32, topK int) ([]domain.Chunk, error) {
	f.topK = topK
	return f.chunks, f.err
}

type fakeGenerator struct {
	answer string
	err    error
	prompt string
	calls  int
}

func (f *fakeGenerator) Generate(_ context.Context, prompt string) (string, error) {
	f.calls++
	f.prompt = prompt
	return f.answer, f.err
}

type fakeSaver struct{}

func (f *fakeSaver) Save(context.Context, []domain.Chunk, [][]float32) error { return nil }

type fakeChunker struct{}

func (f *fakeChunker) Split(string, string, int, int) []domain.Chunk { return nil }

func newAnswerService(chunks []domain.Chunk, maxDistance float64) (*RAGService, *fakeSearcher, *fakeGenerator) {
	searcher := &fakeSearcher{chunks: chunks}
	generator := &fakeGenerator{answer: "answer"}
	service := NewRAGService(
		&fakeChunker{},
		&fakeEmbedder{},
		&fakeSaver{},
		searcher,
		generator,
		maxDistance,
	)
	return service, searcher, generator
}

func TestAnswerFiltersChunksByDistance(t *testing.T) {
	chunks := []domain.Chunk{
		{Source: "doc.md", Index: 0, Content: "близкий", Distance: 0.10},
		{Source: "doc.md", Index: 1, Content: "на границе", Distance: 0.50},
		{Source: "doc.md", Index: 2, Content: "далёкий", Distance: 0.51},
	}
	service, searcher, generator := newAnswerService(chunks, 0.5)

	answer, err := service.Answer(context.Background(), "вопрос")
	if err != nil {
		t.Fatalf("Answer returned error: %v", err)
	}
	if answer != "answer" {
		t.Errorf(`want answer "answer", got %q`, answer)
	}
	if searcher.topK != searchTopK {
		t.Errorf("want topK %d, got %d", searchTopK, searcher.topK)
	}
	if generator.calls != 1 {
		t.Fatalf("want generator called once, got %d", generator.calls)
	}

	for _, want := range []string{"близкий", "на границе"} {
		if !strings.Contains(generator.prompt, want) {
			t.Errorf("prompt does not contain %q:\n%s", want, generator.prompt)
		}
	}
	if strings.Contains(generator.prompt, "далёкий") {
		t.Errorf("prompt contains chunk above threshold:\n%s", generator.prompt)
	}
}

func TestAnswerRejectsQuestionWithoutRelevantChunks(t *testing.T) {
	service, _, generator := newAnswerService([]domain.Chunk{
		{Source: "doc.md", Index: 0, Content: "нерелевантный", Distance: 0.9},
	}, 0.5)

	answer, err := service.Answer(context.Background(), "как приготовить борщ")
	if err != nil {
		t.Fatalf("Answer returned error: %v", err)
	}
	if answer != notAboutDocumentation {
		t.Errorf("want fallback answer, got %q", answer)
	}
	if generator.calls != 0 {
		t.Errorf("generator must not be called, calls = %d", generator.calls)
	}
}

func TestAnswerDeduplicatesChunks(t *testing.T) {
	service, _, generator := newAnswerService([]domain.Chunk{
		{Source: "doc.md", Index: 0, Content: "повтор", Distance: 0.1},
		{Source: "doc.md", Index: 0, Content: "повтор", Distance: 0.1},
	}, 0.5)

	if _, err := service.Answer(context.Background(), "вопрос"); err != nil {
		t.Fatalf("Answer returned error: %v", err)
	}

	if got := strings.Count(generator.prompt, "Content: повтор"); got != 1 {
		t.Errorf("want chunk included once, got %d:\n%s", got, generator.prompt)
	}
}

func TestAnswerUsesDefaultThreshold(t *testing.T) {
	service, _, generator := newAnswerService([]domain.Chunk{
		{Source: "doc.md", Index: 0, Content: "под порогом", Distance: 0.49},
		{Source: "doc.md", Index: 1, Content: "над порогом", Distance: 0.51},
	}, 0)

	if _, err := service.Answer(context.Background(), "вопрос"); err != nil {
		t.Fatalf("Answer returned error: %v", err)
	}

	if !strings.Contains(generator.prompt, "под порогом") {
		t.Error("prompt must contain chunk below default threshold")
	}
	if strings.Contains(generator.prompt, "над порогом") {
		t.Error("prompt must not contain chunk above default threshold")
	}
}

func TestAnswerEmptyQuestion(t *testing.T) {
	service, _, generator := newAnswerService(nil, 0)

	if _, err := service.Answer(context.Background(), ""); err == nil {
		t.Fatal("want error for empty question")
	}
	if generator.calls != 0 {
		t.Errorf("generator must not be called, calls = %d", generator.calls)
	}
}
