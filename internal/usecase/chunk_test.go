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
