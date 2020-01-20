package hocon

import (
	"reflect"
	"testing"
)

var (
	simpleArray = &HoconArray{values: []*HoconValue{
		{values: []HoconElement{NewHoconLiteral("a")}},
		{values: []HoconElement{NewHoconLiteral("b")}},
	}}

	simpleArray2 = &HoconArray{values: []*HoconValue{
		{
			values:   []HoconElement{NewHoconLiteral("current")},
			oldValue: &HoconValue{values: []HoconElement{NewHoconLiteral("old")}},
		},
	}}

	simpleObject = &HoconObject{
		keys:  []string{"a"},
		items: map[string]*HoconValue{"a": {values: []HoconElement{NewHoconLiteral("b")}}},
	}

	simpleLiteral = NewHoconLiteral("a")

	simpleNestedObject = &HoconObject{
		keys:  []string{"a"},
		items: map[string]*HoconValue{"a": {values: []HoconElement{simpleObject}}},
	}
)

func TestHoconSubstitution_GetArray(t *testing.T) {
	cycledSubstitution := &HoconSubstitution{}
	cycledSubstitution.ResolvedValue = &HoconValue{values: []HoconElement{cycledSubstitution}}

	type fields struct {
		Path          string
		ResolvedValue *HoconValue
		IsOptional    bool
		OriginalPath  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*HoconValue
		wantErr bool
	}{
		{
			name: "returns nil if it contains nothing",
			fields: fields{
				ResolvedValue: nil,
			},
			want: nil,
		},
		{
			name: "returns nil if it contains not an array",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleLiteral}},
			},
		},
		{
			name: "returns nil if it contains an array after another value",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleLiteral, simpleArray}},
			},
			want: simpleArray.values,
		},
		{
			name: "returns array if it contains array",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleArray}},
			},
			want: simpleArray.values,
		},
		{
			name: "returns merged array if it contains more than one array",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleArray, simpleArray}},
			},
			want: append(simpleArray.values, simpleArray.values...),
		},
		{
			name: "returns values array only not oldValues",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleArray2}},
			},
			want: simpleArray2.values,
		},
		{
			name: "returns nil if contains cycled reference",
			fields: fields{
				ResolvedValue: cycledSubstitution.ResolvedValue,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconSubstitution{
				Path:          tt.fields.Path,
				ResolvedValue: tt.fields.ResolvedValue,
				IsOptional:    tt.fields.IsOptional,
				OriginalPath:  tt.fields.OriginalPath,
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

func TestHoconSubstitution_GetObject(t *testing.T) {
	cycledSubstitution := &HoconSubstitution{}
	cycledSubstitution.ResolvedValue = &HoconValue{values: []HoconElement{cycledSubstitution}}

	type fields struct {
		Path          string
		ResolvedValue *HoconValue
		IsOptional    bool
		OriginalPath  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *HoconObject
		wantErr bool
	}{
		{
			name: "returns nil if it contains nothing",
			fields: fields{
				ResolvedValue: nil,
			},
			want: nil,
		},
		{
			name: "returns nil if it contains not an object",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleLiteral}},
			},
			want: nil,
		},
		{
			name: "returns nil if it contains an element before an object",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleLiteral, simpleObject}},
			},
			want: nil,
		},
		{
			name: "returns object if it contains object",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleObject}},
			},
			want: simpleObject,
		},
		{
			name: "fails if contains cycled reference",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{cycledSubstitution}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconSubstitution{
				Path:          tt.fields.Path,
				ResolvedValue: tt.fields.ResolvedValue,
				IsOptional:    tt.fields.IsOptional,
				OriginalPath:  tt.fields.OriginalPath,
			}
			got, err := p.GetObject()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetObject() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconSubstitution_GetString(t *testing.T) {
	type fields struct {
		Path          string
		ResolvedValue *HoconValue
		IsOptional    bool
		OriginalPath  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "returns empty string if contains nothing",
			fields: fields{
				ResolvedValue: nil,
			},
			want: "",
		},
		{
			name: "returns empty string if it does not contain string",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleObject}},
			},
			want: "",
		},
		{
			name: "returns string if it contains string",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleLiteral}},
			},
			want: simpleLiteral.String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconSubstitution{
				Path:          tt.fields.Path,
				ResolvedValue: tt.fields.ResolvedValue,
				IsOptional:    tt.fields.IsOptional,
				OriginalPath:  tt.fields.OriginalPath,
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

func TestHoconSubstitution_IsArray(t *testing.T) {
	simpleSubstitution := &HoconSubstitution{
		ResolvedValue: &HoconValue{values: []HoconElement{simpleArray}},
	}
	wrapperSubstitution := &HoconSubstitution{
		ResolvedValue: &HoconValue{values: []HoconElement{simpleSubstitution}},
	}
	type fields struct {
		Path          string
		ResolvedValue *HoconValue
		IsOptional    bool
		OriginalPath  string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "returns true if contains array",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleArray}},
			},
			want: true,
		},
		{
			name: "returns true if contains substitution with array",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleSubstitution}},
			},
			want: true,
		},
		{
			name: "returns true if contains 2 substitution with array",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{wrapperSubstitution}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconSubstitution{
				Path:          tt.fields.Path,
				ResolvedValue: tt.fields.ResolvedValue,
				IsOptional:    tt.fields.IsOptional,
				OriginalPath:  tt.fields.OriginalPath,
			}
			if got := p.IsArray(); got != tt.want {
				t.Errorf("IsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconSubstitution_IsObject(t *testing.T) {
	simpleSubstitution := &HoconSubstitution{
		ResolvedValue: &HoconValue{values: []HoconElement{simpleObject}},
	}
	wrapperSubstitution := &HoconSubstitution{
		ResolvedValue: &HoconValue{values: []HoconElement{simpleSubstitution}},
	}

	type fields struct {
		Path          string
		ResolvedValue *HoconValue
		IsOptional    bool
		OriginalPath  string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "returns true if contains object",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleObject}},
			},
			want: true,
		},
		{
			name: "returns true if contains substitution with object",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleSubstitution}},
			},
			want: true,
		},
		{
			name: "returns true if contains 2 substitution with object",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{wrapperSubstitution}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconSubstitution{
				Path:          tt.fields.Path,
				ResolvedValue: tt.fields.ResolvedValue,
				IsOptional:    tt.fields.IsOptional,
				OriginalPath:  tt.fields.OriginalPath,
			}
			if got, _ := p.IsObject(); got != tt.want {
				t.Errorf("IsObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconSubstitution_IsString(t *testing.T) {
	simpleSubstitution := &HoconSubstitution{
		ResolvedValue: &HoconValue{values: []HoconElement{simpleLiteral}},
	}
	wrapperSubstitution := &HoconSubstitution{
		ResolvedValue: &HoconValue{values: []HoconElement{simpleSubstitution}},
	}

	type fields struct {
		Path          string
		ResolvedValue *HoconValue
		IsOptional    bool
		OriginalPath  string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "returns true if contains string",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleLiteral}},
			},
			want: true,
		},
		{
			name: "returns true if contains substitution with string",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleSubstitution}},
			},
			want: true,
		},
		{
			name: "returns true if contains 2 substitution with string",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{wrapperSubstitution}},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconSubstitution{
				Path:          tt.fields.Path,
				ResolvedValue: tt.fields.ResolvedValue,
				IsOptional:    tt.fields.IsOptional,
				OriginalPath:  tt.fields.OriginalPath,
			}
			if got := p.IsString(); got != tt.want {
				t.Errorf("IsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconSubstitution_checkCycleRef(t *testing.T) {
	cycledSubstitution := &HoconSubstitution{}
	cycledSubstitution.ResolvedValue = &HoconValue{values: []HoconElement{cycledSubstitution}}

	type fields struct {
		Path          string
		ResolvedValue *HoconValue
		IsOptional    bool
		OriginalPath  string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "returns false if does not contain cycle reference",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{simpleLiteral}},
			},
			wantErr: false,
		},
		{
			name: "returns true if contains substitution with cycle reference",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{cycledSubstitution}},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconSubstitution{
				Path:          tt.fields.Path,
				ResolvedValue: tt.fields.ResolvedValue,
				IsOptional:    tt.fields.IsOptional,
				OriginalPath:  tt.fields.OriginalPath,
			}
			if err := p.checkCycleRef(); (err != nil) != tt.wantErr {
				t.Errorf("checkCycleRef() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHoconSubstitution_hasCycleRef(t *testing.T) {
	type fields struct {
		Path          string
		ResolvedValue *HoconValue
		IsOptional    bool
		OriginalPath  string
	}
	type args struct {
		dup   map[HoconElement]int
		level int
		v     interface{}
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
			p := &HoconSubstitution{
				Path:          tt.fields.Path,
				ResolvedValue: tt.fields.ResolvedValue,
				IsOptional:    tt.fields.IsOptional,
				OriginalPath:  tt.fields.OriginalPath,
			}
			if got := p.hasCycleRef(tt.args.dup, tt.args.level, tt.args.v); got != tt.want {
				t.Errorf("hasCycleRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHoconSubstitution(t *testing.T) {
	type args struct {
		path       string
		isOptional bool
	}
	tests := []struct {
		name string
		args args
		want *HoconSubstitution
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHoconSubstitution(tt.args.path, tt.args.isOptional); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHoconSubstitution() = %v, want %v", got, tt.want)
			}
		})
	}
}
