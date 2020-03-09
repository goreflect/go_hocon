package hocon

import (
	"reflect"
	"testing"
)

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
			name: "returns nil if empty",
			want: nil,
		},
		{
			name:   "return empty array correctly",
			fields: fields{values: []*HoconValue{}},
			want:   []*HoconValue{},
		},
		{
			name:   "returns nested array with no elements",
			fields: fields{values: []*HoconValue{{}}},
			want:   []*HoconValue{{}},
		},
		{
			name:   "returns empty element array with empty element",
			fields: fields{values: []*HoconValue{{values: []HoconElement{}}}},
			want:   []*HoconValue{{values: []HoconElement{}}},
		},
		{
			name:   "return its values correctly",
			fields: fields{values: []*HoconValue{wrapInValue(simpleLiteral1)}},
			want:   []*HoconValue{wrapInValue(simpleLiteral1)},
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
		{
			name:    "cannot return a string",
			wantErr: true,
		},
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
		{
			name: "returns true",
			want: true,
		},
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
		{
			name: "always returns false",
			want: false,
		},
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
		{
			name: "empty array returns empty brackets",
			want: "[]",
		},
		{
			name:   "returns empty brackets when do not have elements",
			fields: fields{values: []*HoconValue{}},
			want:   "[]",
		},
		{
			name:   "returns empty brackets when contains element with empty value",
			fields: fields{values: []*HoconValue{{values: []HoconElement{}}}},
			want:   "[]",
		}, {
			name:   "returns its text elements in brackets",
			fields: fields{values: []*HoconValue{wrapInValue(simpleLiteral1)}},
			want:   "[" + simpleLiteral1.value + "]",
		},
		{
			name: "returns its text elements divided by comma",
			fields: fields{values: []*HoconValue{
				wrapInValue(simpleLiteral1),
				wrapInValue(simpleLiteral2),
				wrapInValue(simpleLiteral3),
			}},
			want: "[" + simpleLiteral1.value + "," +
				simpleLiteral2.value + "," +
				simpleLiteral3.value + "]",
		},
		{
			name: "returns its object elements divided by comma",
			fields: fields{
				values: []*HoconValue{
					wrapInValue(
						wrapAllInObject(getArrayOfTwoSimpleKeys(), []HoconElement{simpleLiteral1, simpleLiteral2})),
					wrapInValue(
						wrapInObject(simpleKey3, simpleLiteral3)),
				},
			},
			want: "[{" + newLine +
				"  " + simpleKey1 + " : " + simpleLiteral1.value + newLine +
				"  " + simpleKey2 + " : " + simpleLiteral2.value + newLine +
				"},{" + newLine +
				"  " + simpleKey3 + " : " + simpleLiteral3.value + newLine + "}]",
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
		{name: "returns an empty array", want: &HoconArray{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHoconArray(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHoconArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
