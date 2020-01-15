package hocon

import (
	"reflect"
	"testing"
)

func TestHoconLiteral_GetArray(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*HoconValue
		wantErr bool
	}{
		{name: "literal cannot return an array", fields: fields{}, want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconLiteral{
				value: tt.fields.value,
			}
			got, err := p.GetArray()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetArray() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconLiteral_GetString(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{name: "empty literal returns empty string", fields: fields{""}, want: "", wantErr: false},
		{name: "text literal returns text", fields: fields{"abc"}, want: "abc", wantErr: false},
		{name: "array literal returns text", fields: fields{"[a,b,c,]"}, want: "[a,b,c,]", wantErr: false},
		{name: "integer literal returns text", fields: fields{"123"}, want: "123", wantErr: false},
		{name: "float literal returns text", fields: fields{"123.456"}, want: "123.456", wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconLiteral{
				value: tt.fields.value,
			}
			got, err := p.GetString()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconLiteral_IsArray(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "empty literal is not an array", fields: fields{""}, want: false},
		{name: "empty array literal is not an array", fields: fields{"[]"}, want: false},
		{name: "array literal is not an array", fields: fields{"[a,b,c,]"}, want: false},
		{name: "quoted array literal is not an array", fields: fields{`["a","b","c",]`}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconLiteral{
				value: tt.fields.value,
			}
			if got := p.IsArray(); got != tt.want {
				t.Errorf("IsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconLiteral_IsString(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "empty literal is a string", fields: fields{""}, want: true},
		{name: "text literal is a string", fields: fields{"abc"}, want: true},
		{name: "array literal is a string", fields: fields{"[a,b,c,]"}, want: true},
		{name: "integer literal is a string", fields: fields{"123"}, want: true},
		{name: "float literal is a string", fields: fields{"123.456"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconLiteral{
				value: tt.fields.value,
			}
			if got := p.IsString(); got != tt.want {
				t.Errorf("IsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconLiteral_String(t *testing.T) {
	type fields struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "empty literal returns empty string", fields: fields{""}, want: ""},
		{name: "text literal returns text", fields: fields{"abc"}, want: "abc"},
		{name: "array literal returns text", fields: fields{"[a,b,c,]"}, want: "[a,b,c,]"},
		{name: "quoted array literal returns text", fields: fields{`["a","b","c",]`}, want: `["a","b","c",]`},
		{name: "integer literal returns text", fields: fields{"123"}, want: "123"},
		{name: "float literal returns text", fields: fields{"123.456"}, want: "123.456"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconLiteral{
				value: tt.fields.value,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHoconLiteral(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want *HoconLiteral
	}{
		{name: "constructor with empty string returns empty literal", args: args{""}, want: NewHoconLiteral("")},
		{name: "constructor with text string returns text literal", args: args{"abc"}, want: NewHoconLiteral("abc")},
		{
			name: "constructor with array string returns text literal",
			args: args{"[a,b,c,]"},
			want: NewHoconLiteral("[a,b,c,]"),
		},
		{
			name: "constructor with quoted array string returns text literal",
			args: args{`["a","b","c",]`},
			want: NewHoconLiteral(`["a","b","c",]`),
		},
		{
			name: "constructor with integer string returns text literal",
			args: args{"123"},
			want: NewHoconLiteral("123"),
		},
		{
			name: "constructor with float string returns text literal",
			args: args{"123.456"},
			want: NewHoconLiteral("123.456"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHoconLiteral(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHoconLiteral() = %v, want %v", got, tt.want)
			}
		})
	}
}
