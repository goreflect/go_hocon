package hocon

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type fields struct {
	Tokenizer *Tokenizer
}

type testData struct {
	name   string
	fields fields
	want   bool
}

func getTestData(wrongToken string, append_ bool, correctTokens ...string) []testData {
	tests := []testData{
		{
			name: "returns false with nil tokenizer",
			want: false,
		},
		{
			name: "returns false with no needed token",
			fields: fields{
				Tokenizer: NewTokenizer(wrongToken),
			},
			want: false,
		},
	}

	for _, correctToken := range correctTokens {
		if append_ {
			correctToken += wrongToken
		}

		correctTokenTests := []testData{
			{
				name: "returns false needed token not first",
				fields: fields{
					Tokenizer: NewTokenizer(wrongToken + correctToken),
				},
				want: false,
			},
			{
				name: "returns true with needed token",
				fields: fields{
					Tokenizer: NewTokenizer(correctToken),
				},
				want: true,
			},
		}

		tests = append(tests, correctTokenTests...)
	}

	return tests
}

func TestHoconTokenizer_IsArrayEnd(t *testing.T) {
	tests := getTestData("a", false, arrayEndToken)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsArrayEnd(); got != tt.want {
				t.Errorf("IsArrayEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsArrayStart(t *testing.T) {
	tests := getTestData("a", false, arrayStartToken)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsArrayStart(); got != tt.want {
				t.Errorf("IsArrayStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsAssignment(t *testing.T) {
	tests := getTestData("a", true, assignmentTokens...)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsAssignment(); got != tt.want {
				t.Errorf("IsAssignment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsComma(t *testing.T) {
	tests := getTestData("a", false, commaToken)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsComma(); got != tt.want {
				t.Errorf("IsComma() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsDot(t *testing.T) {
	tests := getTestData("a", false, dotToken)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsDot(); got != tt.want {
				t.Errorf("IsDot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsEndOfObject(t *testing.T) {
	tests := getTestData("a", false, endOfObjectToken)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsEndOfObject(); got != tt.want {
				t.Errorf("IsEndOfObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsInclude(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "returns false if does not have any continuation",
			fields: fields{
				Tokenizer: NewTokenizer(includeSpecial),
			},
			want: false,
		},
		{
			name: "returns false if continues with whitespaces only",
			fields: fields{
				Tokenizer: NewTokenizer(includeSpecial + "  "),
			},
			want: false,
		},
		{
			name: "returns false if continues with whitespaces and comments",
			fields: fields{
				Tokenizer: NewTokenizer(includeSpecial + "  // comment"),
			},
			want: false,
		},
		{
			name: "returns false if does not have whitespace before quoted text",
			fields: fields{
				Tokenizer: NewTokenizer(includeSpecial + `"text"`),
			},
			want: false,
		},
		{
			name: "returns false if continued by whitespace and unquoted text",
			fields: fields{
				Tokenizer: NewTokenizer(includeSpecial + ` text`),
			},
			want: false,
		},
		{
			name: "returns true if continued by whitespace and quoted text",
			fields: fields{
				Tokenizer: NewTokenizer(includeSpecial + ` "text"`),
			},
			want: true,
		},
		{
			name: "returns true if continued by comment and quoted text",
			fields: fields{
				Tokenizer: NewTokenizer(includeSpecial + "//comment\n\"text\""),
			},
			want: true,
		},
		{
			name: "returns false if contains errors",
			fields: fields{
				Tokenizer: NewTokenizer(includeSpecial + " \"text\\z"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsInclude(); got != tt.want {
				t.Errorf("IsInclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsNewline(t *testing.T) {
	tests := getTestData("a", false, newLineToken)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsNewline(); got != tt.want {
				t.Errorf("IsNewline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsObjectStart(t *testing.T) {
	tests := getTestData("a", false, objectStartToken)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsObjectStart(); got != tt.want {
				t.Errorf("IsObjectStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsPlusAssignment(t *testing.T) {
	tests := getTestData("a", false, plusAssignmentToken)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsPlusAssignment(); got != tt.want {
				t.Errorf("IsPlusAssignment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsSpaceOrTab(t *testing.T) {
	tests := getTestData("a", true, spaceOrTabTokens...)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsSpaceOrTab(); got != tt.want {
				t.Errorf("IsSpaceOrTab() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsStartOfComment(t *testing.T) {
	tests := getTestData("a", true, startOfCommentTokens...)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsStartOfComment(); got != tt.want {
				t.Errorf("IsStartOfComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsStartOfQuotedText(t *testing.T) {
	tests := getTestData("a", false, startOfQuotedTextToken)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsStartOfQuotedText(); got != tt.want {
				t.Errorf("IsStartOfQuotedText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsStartOfTripleQuotedText(t *testing.T) {
	tests := getTestData("a", false, startOfTripleQuotedTextToken)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsStartOfTripleQuotedText(); got != tt.want {
				t.Errorf("IsStartOfTripleQuotedText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsStartSimpleValue(t *testing.T) {
	tests := getTestData("\\", true, "a")
	tests = append(tests, testData{
		name: "returns true if starts with more than one space",
		fields: fields{
			Tokenizer: NewTokenizer("  "),
		},
		want: true,
	})
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsStartSimpleValue(); got != tt.want {
				t.Errorf("IsStartSimpleValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsSubstitutionStart(t *testing.T) {
	tests := getTestData("a", true, substitutionStartTokens...)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsSubstitutionStart(); got != tt.want {
				t.Errorf("IsSubstitutionStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsUnquotedKey(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsUnquotedKey(); got != tt.want {
				t.Errorf("IsUnquotedKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsUnquotedKeyStart(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsUnquotedKeyStart(); got != tt.want {
				t.Errorf("IsUnquotedKeyStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsWhitespace(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsWhitespace(); got != tt.want {
				t.Errorf("IsWhitespace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_IsWhitespaceOrComment(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.IsWhitespaceOrComment(); got != tt.want {
				t.Errorf("IsWhitespaceOrComment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullArrayEnd(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   *Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.PullArrayEnd(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullArrayEnd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullArrayStart(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   *Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.PullArrayStart(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullArrayStart() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullAssignment(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   *Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.PullAssignment(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullAssignment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullComma(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   *Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.PullComma(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullComma() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullComment(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   *Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got := p.PullComment()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullComment() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullDot(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   *Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.PullDot(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullDot() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullEndOfObject(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   *Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.PullEndOfObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullEndOfObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullInclude(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got, err := p.PullInclude()
			if (err != nil) != tt.wantErr {
				t.Errorf("PullInclude() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullInclude() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullNewline(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   *Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.PullNewline(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullNewline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullNext(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got, err := p.PullNext()
			if (err != nil) != tt.wantErr {
				t.Errorf("PullNext() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullNext() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullPlusAssignment(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   *Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.PullPlusAssignment(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullPlusAssignment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullQuotedKey(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got, err := p.PullQuotedKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("PullQuotedKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullQuotedKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullQuotedText(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got, err := p.PullQuotedText()
			if (err != nil) != tt.wantErr {
				t.Errorf("PullQuotedText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullQuotedText() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullRestOfLine(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got := p.PullRestOfLine()
			if got != tt.want {
				t.Errorf("PullRestOfLine() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullSimpleValue(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got, err := p.PullSimpleValue()
			if (err != nil) != tt.wantErr {
				t.Errorf("PullSimpleValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullSimpleValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullSpaceOrTab(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got, err := p.PullSpaceOrTab()
			if (err != nil) != tt.wantErr {
				t.Errorf("PullSpaceOrTab() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullSpaceOrTab() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullStartOfObject(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   *Token
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.PullStartOfObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullStartOfObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullTripleQuotedText(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got, err := p.PullTripleQuotedText()
			if (err != nil) != tt.wantErr {
				t.Errorf("PullTripleQuotedText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullTripleQuotedText() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullUnquotedKey(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got, err := p.PullUnquotedKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("PullUnquotedKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullUnquotedKey() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullValue(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got, err := p.PullValue()
			if (err != nil) != tt.wantErr {
				t.Errorf("PullValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PullValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_PullWhitespaceAndComments(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			p.PullWhitespaceAndComments()
		})
	}
}

func TestHoconTokenizer_isStartOfQuotedKey(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.isStartOfQuotedKey(); got != tt.want {
				t.Errorf("isStartOfQuotedKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_isUnquotedText(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.isUnquotedText(); got != tt.want {
				t.Errorf("isUnquotedText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_isValue(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			if got := p.isValue(); got != tt.want {
				t.Errorf("isValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_pullEscapeSequence(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got, err := p.pullEscapeSequence()
			if (err != nil) != tt.wantErr {
				t.Errorf("pullEscapeSequence() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("pullEscapeSequence() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_pullSubstitution(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got, err := p.pullSubstitution()
			if (err != nil) != tt.wantErr {
				t.Errorf("pullSubstitution() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pullSubstitution() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconTokenizer_pullUnquotedText(t *testing.T) {
	type fields struct {
		Tokenizer *Tokenizer
	}
	tests := []struct {
		name    string
		fields  fields
		want    *Token
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			got, err := p.pullUnquotedText()
			if (err != nil) != tt.wantErr {
				t.Errorf("pullUnquotedText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pullUnquotedText() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHoconTokenizer(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want *HoconTokenizer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHoconTokenizer(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHoconTokenizer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTokenizer(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want *Tokenizer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTokenizer(tt.args.text); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenizer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenizer_EOF(t *testing.T) {
	type fields struct {
		text       string
		index      int
		indexStack *Stack
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Tokenizer{
				text:       tt.fields.text,
				index:      tt.fields.index,
				indexStack: tt.fields.indexStack,
			}
			if got := p.EOF(); got != tt.want {
				t.Errorf("EOF() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenizer_Matches(t *testing.T) {
	type fields struct {
		text       string
		index      int
		indexStack *Stack
	}
	type args struct {
		pattern string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Tokenizer{
				text:       tt.fields.text,
				index:      tt.fields.index,
				indexStack: tt.fields.indexStack,
			}
			if got := p.Matches(tt.args.pattern); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenizer_MatchesMore(t *testing.T) {
	type fields struct {
		text       string
		index      int
		indexStack *Stack
	}
	type args struct {
		patterns []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Tokenizer{
				text:       tt.fields.text,
				index:      tt.fields.index,
				indexStack: tt.fields.indexStack,
			}
			if got := p.MatchesMore(tt.args.patterns); got != tt.want {
				t.Errorf("MatchesMore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenizer_Peek(t *testing.T) {
	type fields struct {
		text       string
		index      int
		indexStack *Stack
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Tokenizer{
				text:       tt.fields.text,
				index:      tt.fields.index,
				indexStack: tt.fields.indexStack,
			}
			if got := p.Peek(); got != tt.want {
				t.Errorf("Peek() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenizer_Pop(t *testing.T) {
	type fields struct {
		text       string
		index      int
		indexStack *Stack
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Tokenizer{
				text:       tt.fields.text,
				index:      tt.fields.index,
				indexStack: tt.fields.indexStack,
			}
			if err := p.Pop(); (err != nil) != tt.wantErr {
				t.Errorf("Pop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestTokenizer_PullWhitespace(t *testing.T) {
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconTokenizer{
				Tokenizer: tt.fields.Tokenizer,
			}
			assert.NotPanics(t, p.PullWhitespace, "PullWhitespace() panicked")
		})
	}
}

func TestTokenizer_Push(t *testing.T) {
	type fields struct {
		text       string
		index      int
		indexStack *Stack
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Tokenizer{
				text:       tt.fields.text,
				index:      tt.fields.index,
				indexStack: tt.fields.indexStack,
			}
			assert.NotPanics(t, p.Push, "Push() panicked")
		})
	}
}

func TestTokenizer_Take(t *testing.T) {
	type fields struct {
		text       string
		index      int
		indexStack *Stack
	}
	type args struct {
		length int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Tokenizer{
				text:       tt.fields.text,
				index:      tt.fields.index,
				indexStack: tt.fields.indexStack,
			}
			if got := p.Take(tt.args.length); got != tt.want {
				t.Errorf("Take() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTokenizer_TakeOne(t *testing.T) {
	type fields struct {
		text       string
		index      int
		indexStack *Stack
	}
	tests := []struct {
		name   string
		fields fields
		want   byte
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Tokenizer{
				text:       tt.fields.text,
				index:      tt.fields.index,
				indexStack: tt.fields.indexStack,
			}
			if got := p.TakeOne(); got != tt.want {
				t.Errorf("TakeOne() = %v, want %v", got, tt.want)
			}
		})
	}
}
