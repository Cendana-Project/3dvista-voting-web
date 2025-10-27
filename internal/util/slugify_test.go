package util

import "testing"

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "simple text",
			input: "Hello World",
			want:  "hello-world",
		},
		{
			name:  "with special characters",
			input: "Hello, World!",
			want:  "hello-world",
		},
		{
			name:  "with numbers",
			input: "Test 123",
			want:  "test-123",
		},
		{
			name:  "with multiple spaces",
			input: "Hello   World",
			want:  "hello-world",
		},
		{
			name:  "with leading/trailing spaces",
			input: "  Hello World  ",
			want:  "hello-world",
		},
		{
			name:  "with hyphens",
			input: "Hello-World",
			want:  "hello-world",
		},
		{
			name:  "with underscores",
			input: "Hello_World",
			want:  "hello-world",
		},
		{
			name:  "Indonesian text",
			input: "Jabar Digital Academy",
			want:  "jabar-digital-academy",
		},
		{
			name:  "with parentheses",
			input: "Data Potensi Digital Desa ( TAPAL DESA )",
			want:  "data-potensi-digital-desa-tapal-desa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Slugify(tt.input)
			if got != tt.want {
				t.Errorf("Slugify() = %v, want %v", got, tt.want)
			}
		})
	}
}

