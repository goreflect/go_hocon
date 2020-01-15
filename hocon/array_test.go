package hocon

import (
	"reflect"
	"testing"
)

const newLine string = "\r\n"

func TestHoconArray_GetArray(t *testing.T) {
	type fields struct {
		values []*HoconValue
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*HoconValue
		wantErr bool
	}{
		{
			name:    "array without array returns nil",
			fields:  fields{},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "array with empty nested array return empty array",
			fields:  fields{values: []*HoconValue{}},
			want:    []*HoconValue{},
			wantErr: false,
		},
		{
			name:    "array with no elements returns nested array with no elements",
			fields:  fields{values: []*HoconValue{{}}},
			want:    []*HoconValue{{}},
			wantErr: false,
		},
		{
			name:    "array with empty element returns empty element",
			fields:  fields{values: []*HoconValue{{values: []HoconElement{}}}},
			want:    []*HoconValue{{values: []HoconElement{}}},
			wantErr: false,
		},
		{
			name:    "array with values return its values",
			fields:  fields{values: []*HoconValue{{values: []HoconElement{NewHoconLiteral("a")}}}},
			want:    []*HoconValue{{values: []HoconElement{NewHoconLiteral("a")}}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconArray{
				values: tt.fields.values,
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

func TestHoconArray_GetString(t *testing.T) {
	type fields struct {
		values []*HoconValue
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{name: "array cannot return a string", fields: fields{}, want: "", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconArray{
				values: tt.fields.values,
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

func TestHoconArray_IsArray(t *testing.T) {
	type fields struct {
		values []*HoconValue
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "array is array", fields: fields{}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconArray{
				values: tt.fields.values,
			}
			if got := p.IsArray(); got != tt.want {
				t.Errorf("IsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconArray_IsString(t *testing.T) {
	type fields struct {
		values []*HoconValue
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{name: "array is not a string", fields: fields{}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconArray{
				values: tt.fields.values,
			}
			if got := p.IsString(); got != tt.want {
				t.Errorf("IsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconArray_String(t *testing.T) {
	type fields struct {
		values []*HoconValue
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{name: "empty array returns empty brackets", fields: fields{}, want: "[]"},
		{name: "array with no elements returns empty brackets", fields: fields{values: []*HoconValue{}}, want: "[]"},
		{name: "array with empty element returns empty brackets", fields: fields{values: []*HoconValue{{}}}, want: "[]"},
		{
			name:   "array with element with empty value returns empty brackets",
			fields: fields{values: []*HoconValue{{values: []HoconElement{}}}},
			want:   "[]",
		}, {
			name:   "array returns its text elements in brackets",
			fields: fields{values: []*HoconValue{{values: []HoconElement{NewHoconLiteral("a")}}}},
			want:   "[a]",
		},
		{
			name: "array returns its text elements divided by comma",
			fields: fields{values: []*HoconValue{
				{values: []HoconElement{NewHoconLiteral("a")}},
				{values: []HoconElement{NewHoconLiteral("b")}},
				{values: []HoconElement{NewHoconLiteral("1")}},
			}},
			want: "[a,b,1]",
		},
		{
			name: "array returns its object elements divided by comma",
			fields: fields{
				values: []*HoconValue{
					{values: []HoconElement{makeHoconObject([]string{"a", "c"}, []string{"b", "d"})}},
					{values: []HoconElement{makeHoconObject([]string{"e"}, []string{"f"})}},
				},
			},
			want: "[{" + newLine + "  a : b" + newLine + "  c : d" + newLine + "},{" + newLine + "  e : f" + newLine + "}]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconArray{
				values: tt.fields.values,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHoconArray(t *testing.T) {
	tests := []struct {
		name string
		want *HoconArray
	}{
		{name: "NewHoconArray returns an empty array", want: &HoconArray{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHoconArray(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHoconArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
