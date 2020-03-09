package hocon

import (
	"reflect"
	"testing"
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
				ResolvedValue: wrapInValue(simpleLiteral1),
			},
		},
		{
			name: "returns nil if it contains an array after another value",
			fields: fields{
				ResolvedValue: wrapInValue(simpleLiteral1, simpleTwoValuesArray),
			},
			want: simpleTwoValuesArray.values,
		},
		{
			name: "returns array if it contains array",
			fields: fields{
				ResolvedValue: wrapInValue(simpleTwoValuesArray),
			},
			want: simpleTwoValuesArray.values,
		},
		{
			name: "returns merged array if it contains more than one array",
			fields: fields{
				ResolvedValue: wrapInValue(simpleTwoValuesArray, simpleTwoValuesArray),
			},
			want: append(simpleTwoValuesArray.values, simpleTwoValuesArray.values...),
		},
		{
			name: "returns values array only not oldValues",
			fields: fields{
				ResolvedValue: wrapInValue(simpleArrayWithOldValue),
			},
			want: simpleArrayWithOldValue.values,
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
				ResolvedValue: wrapInValue(simpleLiteral1),
			},
			want: nil,
		},
		{
			name: "returns nil if it contains an element before an object",
			fields: fields{
				ResolvedValue: wrapInValue(simpleLiteral1, simpleObject),
			},
			want: nil,
		},
		{
			name: "returns object if it contains object",
			fields: fields{
				ResolvedValue: wrapInValue(simpleObject),
			},
			want: simpleObject,
		},
		{
			name: "fails if contains cycled reference",
			fields: fields{
				ResolvedValue: &HoconValue{values: []HoconElement{getCycledSubstitution()}},
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
				ResolvedValue: wrapInValue(simpleObject),
			},
			want: "",
		},
		{
			name: "returns string if it contains string",
			fields: fields{
				ResolvedValue: wrapInValue(simpleLiteral1),
			},
			want: simpleLiteral1.value,
		},
		{
			name: "fails if contains cycled substitution",
			fields: fields{
				ResolvedValue: wrapInValue(getCycledSubstitution()),
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
			name: "returns false if contains nothing",
			want: false,
		},
		{
			name: "returns true if contains array",
			fields: fields{
				ResolvedValue: wrapInValue(simpleTwoValuesArray),
			},
			want: true,
		},
		{
			name: "returns true if contains substitution with array",
			fields: fields{
				ResolvedValue: wrapInValue(wrapInSubstitution(simpleTwoValuesArray)),
			},
			want: true,
		},
		{
			name: "returns true if contains nested substitution with array",
			fields: fields{
				ResolvedValue: wrapInValue(wrapInSubstitution(wrapInSubstitution(simpleTwoValuesArray))),
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
			name: "returns false if contains nothing",
			want: false,
		},
		{
			name: "returns true if contains object",
			fields: fields{
				ResolvedValue: wrapInValue(simpleObject),
			},
			want: true,
		},
		{
			name: "returns true if contains substitution with object",
			fields: fields{
				ResolvedValue: wrapInValue(wrapInSubstitution(simpleObject)),
			},
			want: true,
		},
		{
			name: "returns true if contains nested substitution with object",
			fields: fields{
				ResolvedValue: wrapInValue(wrapInSubstitution(wrapInSubstitution(simpleObject))),
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
			name:   "returns false if contains nothing",
			fields: fields{},
			want:   false,
		},
		{
			name: "returns true if contains string",
			fields: fields{
				ResolvedValue: wrapInValue(simpleLiteral1),
			},
			want: true,
		},
		{
			name: "returns true if contains substitution with string",
			fields: fields{
				ResolvedValue: wrapInValue(wrapInSubstitution(simpleLiteral1)),
			},
			want: true,
		},
		{
			name: "returns true if contains nested substitution with string",
			fields: fields{
				ResolvedValue: wrapInValue(wrapInSubstitution(wrapInSubstitution(simpleLiteral1))),
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
			name: "returns nil if does not contain cycle reference",
			fields: fields{
				ResolvedValue: wrapInValue(simpleLiteral1),
			},
			wantErr: false,
		},
		{
			name: "returns error if contains substitution with cycle reference",
			fields: fields{
				ResolvedValue: wrapInValue(getCycledSubstitution()),
			},
			wantErr: true,
		},
		{
			name: "returns error if contains nested substitution with cycle reference",
			fields: fields{
				ResolvedValue: wrapInValue(wrapInSubstitution(getCycledSubstitution())),
			},
			wantErr: true,
		},
		{
			name: "returns error if contains double nested substitution with cycle reference",
			fields: fields{
				ResolvedValue: wrapInValue(wrapInSubstitution(wrapInSubstitution(getCycledSubstitution()))),
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
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "returns false if contains nothing",
			want: false,
		},
		{
			name: "returns true if has direct cycled reference",
			fields: fields{
				ResolvedValue: wrapInValue(getCycledSubstitution()),
			},
			args: args{
				dup: map[HoconElement]int{},
			},
			want: true,
		},
		{
			name: "returns true if has nested cycled reference",
			fields: fields{
				ResolvedValue: wrapInValue(wrapInSubstitution(getCycledSubstitution())),
			},
			args: args{
				dup: map[HoconElement]int{},
			},
			want: true,
		},
		{
			name: "returns false if has cycled reference nested in value",
			fields: fields{
				ResolvedValue: wrapInValue(wrapInValue(getCycledSubstitution())),
			},
			args: args{
				dup: map[HoconElement]int{},
			},
			want: false,
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
			if got := p.hasCycleRef(tt.args.dup, tt.args.level); got != tt.want {
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
		{
			name: "returns empty substitution",
			args: args{},
			want: &HoconSubstitution{},
		},
		{
			name: "returns empty substitution",
			args: args{
				path:       simpleKey1,
				isOptional: true,
			},
			want: &HoconSubstitution{
				Path:          simpleKey1,
				ResolvedValue: nil,
				IsOptional:    true,
				OriginalPath:  simpleKey1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHoconSubstitution(tt.args.path, tt.args.isOptional); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHoconSubstitution() = %v, want %v", got, tt.want)
			}
		})
	}
}
