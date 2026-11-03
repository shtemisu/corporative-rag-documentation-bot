package usecase

import (
	"reflect"
	"testing"
)

func TestParseHeading(t *testing.T) {
	tests := []struct {
		num       int
		in        string
		wantLevel int
		wantTitle string
		wantOK    bool
	}{
		{1, "## Деплой", 2, "Деплой", true},
		{2, "##  Деплой", 2, "Деплой", true},
		{3, "### Деплой", 3, "Деплой", true},
		{4, "#### Деплой", 4, "Деплой", true},
		{5, "##### Деплой", 5, "Деплой", true},
		{6, "###### Деплой", 6, "Деплой", true},
		{7, "####### Деплой", 0, "", false},
		{8, "#hashtag", 0, "", false},
		{9, "заголовок", 0, "", false},
		{10, " отступ", 0, "", false},
		{11, "## ", 0, "", false},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			level, title, ok := parseHeading(tt.in)
			if ok != tt.wantOK || title != tt.wantTitle || tt.wantLevel != level {
				t.Errorf(`TEST NUMBER %d
						want level: %d, get level: %d
						want title: "%s", get title: "%s"
						want OK: %t, get OK: %t`, tt.num, tt.wantLevel,
					level, tt.wantTitle, title, tt.wantOK, ok)
			}
		})
	}
}

func TestSplitSections(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want []section
	}{
		{
			name: "пустой вход",
			in:   "",
			want: []section{},
		},
		{
			name: "текст до первого заголовка и возврат на уровень выше",
			in:   "Введение\n# A\nтекст a\n## B\nтекст b\n# C\nтекст c",
			want: []section{
				{heading: "", body: "Введение"},
				{heading: "A", body: "текст a"},
				{heading: "A > B", body: "текст b"},
				{heading: "C", body: "текст c"},
			},
		},
		{
			name: "решётка внутри блока кода",
			in:   "# A\n```bash\n# это комментарий\n```\nтекст",
			want: []section{
				{heading: "A", body: "```bash\n# это комментарий\n```\nтекст"},
			},
		},
		{
			name: "глубокая вложенность",
			in:   "# A\nа\n## B\nб\n### C\nв\n#### D\nг",
			want: []section{
				{heading: "A", body: "а"},
				{heading: "A > B", body: "б"},
				{heading: "A > B > C", body: "в"},
				{heading: "A > B > C > D", body: "г"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitSections(tt.in)
			ok := reflect.DeepEqual(got, tt.want)
			if !ok {
				t.Errorf("want: %q, but got: %q", tt.want, got)
			}
		})
	}
}

func TestPackParagraphs(t *testing.T) {
	tests := []struct {
		name   string
		in     []string
		want   []string
		maxLen int
	}{
		{
			name:   "пустой",
			in:     []string{},
			want:   []string{},
			maxLen: 10,
		},
		{
			name:   "Все влезает",
			in:     []string{"aa", "bb"},
			want:   []string{"aa\n\nbb"},
			maxLen: 100,
		},
		{
			name:   "Много параграфов",
			in:     []string{"aaaa", "bbb", "cccccc", "dd"},
			want:   []string{"aaaa\n\nbbb", "cccccc\n\ndd"},
			maxLen: 10,
		},
		{
			name:   "Ровно на границе",
			in:     []string{"aaaa", "bbbb"},
			want:   []string{"aaaa\n\nbbbb"},
			maxLen: 10,
		},
		{
			name:   "На 1 больше",
			in:     []string{"aaaaa", "bbbb"},
			want:   []string{"aaaaa", "bbbb"},
			maxLen: 10,
		},
		{
			name:   "кириллица",
			in:     []string{"абвг", "деёж"},
			want:   []string{"абвг\n\nдеёж"},
			maxLen: 10,
		},
		{
			name:   "длинный параграф",
			in:     []string{"aaaaaaaaaaaaaaaa"},
			want:   []string{"aaaaaaaaaaaaaaaa"},
			maxLen: 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := packParagraphs(tt.in, tt.maxLen)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want: %q, but got: %q", tt.want, got)
			}
		})
	}
}

func TestTailRunes(t *testing.T) {
	tests := []struct {
		name  string
		inStr string
		inN   int
		want  string
	}{
		{
			name:  "empty string",
			inStr: "",
			inN:   10,
			want:  "",
		},
		{
			name:  "non-empty string",
			inStr: "Документация",
			inN:   5,
			want:  "тация",
		},
		{
			name:  "with space",
			inStr: "Documen tation",
			inN:   8,
			want:  "n tation",
		},
		{
			name:  "len string less than n",
			inStr: "Hello!",
			inN:   10,
			want:  "Hello!",
		},
		{
			name:  "n is equal zero",
			inStr: "Hello!",
			inN:   0,
			want:  "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tailRunes(tt.inStr, tt.inN)
			if tt.want != got {
				t.Errorf("want: %q, but got: %qd", tt.want, got)
			}
		})
	}
}

func TestTailRunesWithSpace(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "empty string",
			in:   "",
			want: "",
		},
		{
			name: "non-empty string",
			in:   "Доку ментация",
			want: "ментация",
		},
		{
			name: "without space",
			in:   "Documentation",
			want: "Documentation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tailRunesWithSpace(tt.in)
			if tt.want != got {
				t.Errorf("want: %s, but got: %s", tt.want, got)
			}
		})
	}
}

var source string = "doc.md"

func TestSplitMarkdown(t *testing.T) {
	tests := []struct {
		name    string
		text    string
		maxLen  int
		overlap int
		want    []Chunk
	}{
		{
			name:    "empty text",
			text:    "",
			maxLen:  100,
			overlap: 0,
			want:    []Chunk{},
		},
		{
			name:    "only one section",
			text:    "A\naaaa\n\nbbbb",
			maxLen:  100,
			overlap: 0,
			want: []Chunk{
				{
					Source:  source,
					Heading: "",
					Index:   0,
					Content: "A\naaaa\n\nbbbb",
				},
			},
		},
		{
			name:    "without header",
			text:    "just a text",
			maxLen:  100,
			overlap: 0,
			want: []Chunk{
				{
					Source:  source,
					Heading: "",
					Index:   0,
					Content: "just a text",
				},
			},
		},
		{
			name:    "sequential index",
			text:    "# A\nx\n# B\ny",
			maxLen:  100,
			overlap: 0,
			want: []Chunk{
				{
					Source:  source,
					Heading: "A",
					Index:   0,
					Content: "A\n\nx",
				},
				{
					Source:  source,
					Heading: "B",
					Index:   1,
					Content: "B\n\ny",
				},
			},
		},
		{
			name:    "nested header",
			text:    "# A\n\n## B\nтекст",
			maxLen:  100,
			overlap: 0,
			want: []Chunk{
				{
					Source:  source,
					Heading: "A > B",
					Index:   0,
					Content: "A > B\n\nтекст",
				},
			},
		},
		{
			name:    "overlap into section",
			text:    "# A\naaaa\n\nbbbb",
			maxLen:  4,
			overlap: 2,
			want: []Chunk{
				{
					Source:  source,
					Heading: "A",
					Index:   0,
					Content: "A\n\naaaa",
				},
				{
					Source:  source,
					Heading: "A",
					Index:   1,
					Content: "A\n\naa\n\nbbbb",
				},
			},
		},
		{
			name:    "no overlap between sections",
			text:    "# A\naaaa\n# B\nbbbb",
			maxLen:  100,
			overlap: 2,
			want: []Chunk{
				{
					Source:  source,
					Heading: "A",
					Index:   0,
					Content: "A\n\naaaa",
				},
				{
					Source:  source,
					Heading: "B",
					Index:   1,
					Content: "B\n\nbbbb",
				},
			},
		},
		{
			name:    "overlap with cyrilic",
			text:    "# Раздел\nабвгд\n\nежзий",
			maxLen:  5,
			overlap: 3,
			want: []Chunk{
				{
					Source:  source,
					Heading: "Раздел",
					Index:   0,
					Content: "Раздел\n\nабвгд",
				},
				{
					Source:  source,
					Heading: "Раздел",
					Index:   1,
					Content: "Раздел\n\nвгд\n\nежзий",
				},
			},
		},
		{
			name:    "empty paragraphs",
			text:    "# A\n\naa\n\n\n\n\nbb",
			maxLen:  100,
			overlap: 0,
			want: []Chunk{
				{
					Source:  source,
					Heading: "A",
					Index:   0,
					Content: "A\n\naa\n\nbb",
				},
			},
		},
		{
			name:    "maxlen is default",
			text:    "# A\n\naaaa\n\nbbbb",
			maxLen:  0,
			overlap: 0,
			want: []Chunk{
				{
					Source:  source,
					Heading: "A",
					Index:   0,
					Content: "A\n\naaaa\n\nbbbb",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SplitMarkdown(source, tt.text, tt.maxLen, tt.overlap)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("want: %#v, but got: %#v", tt.want, got)
			}
		})
	}
}
