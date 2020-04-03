package hocon

import (
	"reflect"
	"testing"
)

func TestHoconObject_GetArray(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []*HoconValue
		wantErr bool
	}{
		{
			name:    "cannot return an array",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
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

func TestHoconObject_GetKey(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *HoconValue
	}{
		{
			name: "returns correct value by key",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
				keys:  getArrayOfTwoSimpleKeys(),
			},
			args: args{key: simpleKey1},
			want: wrapInValue(simpleLiteral1),
		},
		{
			name: "returns nil by unknown key",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
				keys:  getArrayOfTwoSimpleKeys(),
			},
			args: args{key: simpleKey3},
			want: nil,
		},
		{
			name: "returns nil out of nil key/value",
			args: args{key: simpleKey1},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			if got := p.GetKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconObject_GetKeys(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		{
			name: "returns existed keys correctly",
			fields: fields{
				keys: getArrayOfTwoSimpleKeys(),
			},
			want: getArrayOfTwoSimpleKeys(),
		},
		{
			name: "returns empty slice correctly",
			fields: fields{
				keys: []string{},
			},
			want: []string{},
		},
		{
			name: "returns nil instead of keys", // todo maybe should return empty list
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			if got := p.GetKeys(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKeys() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconObject_GetOrCreateKey(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *HoconValue
	}{
		{
			name: "returns current value as oldValue",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
				keys:  getArrayOfTwoSimpleKeys(),
			},
			args: args{key: simpleKey1},
			want: &HoconValue{oldValue: wrapInValue(simpleLiteral1)},
		},
		{
			name: "returns empty value if it didn't exist",
			fields: fields{
				items: map[string]*HoconValue{simpleKey1: wrapInValue(simpleLiteral1)},
				keys:  []string{simpleKey1},
			},
			args: args{key: simpleKey2},
			want: &HoconValue{},
		},
		{
			name: "returns empty value if it didn't have any fields",
			args: args{key: simpleKey1},
			want: &HoconValue{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			if got := p.GetOrCreateKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetOrCreateKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconObject_GetString(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
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
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
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

func TestHoconObject_IsArray(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
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
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			if got := p.IsArray(); got != tt.want {
				t.Errorf("IsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconObject_IsString(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
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
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			if got := p.IsString(); got != tt.want {
				t.Errorf("IsString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconObject_Items(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
	}
	tests := []struct {
		name   string
		fields fields
		want   map[string]*HoconValue
	}{
		{
			name: "returns its values correctly",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
			},
			want: getMapOfTwoSimpleLiterals(),
		},
		{
			name: "returns nil items if it's empty",
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			if got := p.Items(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Items() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconObject_Merge(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
	}
	type args struct {
		other *HoconObject
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *HoconObject
	}{
		{
			name: "merges with nil correctly",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
				keys:  getArrayOfTwoSimpleKeys(),
			},
			args: args{other: nil},
			want: wrapAllInObject(
				getArrayOfTwoSimpleKeys(),
				[]HoconElement{simpleLiteral1, simpleLiteral2}),
		},
		{
			name: "merges with empty object correctly",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
				keys:  getArrayOfTwoSimpleKeys(),
			},
			args: args{other: wrapAllInObject([]string{}, []HoconElement{})},
			want: wrapAllInObject(
				getArrayOfTwoSimpleKeys(),
				[]HoconElement{simpleLiteral1, simpleLiteral2}),
		},
		{
			name: "merges with other correctly",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
				keys:  getArrayOfTwoSimpleKeys(),
			},
			args: args{
				other: wrapAllInObject(
					[]string{simpleKey3},
					[]HoconElement{simpleLiteral3}),
			},
			want: wrapAllInObject(
				[]string{simpleKey1, simpleKey2, simpleKey3},
				[]HoconElement{simpleLiteral1, simpleLiteral2, simpleLiteral3}),
		},
		{
			name: "fails to merge cycled object with the same key",
			fields: fields{
				items: map[string]*HoconValue{
					simpleKey1: wrapInValue(getCycledSubstitution()),
				},
				keys: []string{simpleKey1},
			},
			args: args{
				other: wrapAllInObject(
					[]string{simpleKey1},
					[]HoconElement{wrapInValue(simpleLiteral1)}),
			},
			want: &HoconObject{
				items: map[string]*HoconValue{
					simpleKey1: wrapInValue(getCycledSubstitution()),
				},
				keys: []string{simpleKey1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}

			p.Merge(tt.args.other)

			if !reflect.DeepEqual(p, tt.want) {
				t.Errorf("Merge() got = %v, want %v", p, tt.want)
			}
		})
	}
}

func TestHoconObject_MergeImmutable(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
	}
	type args struct {
		other *HoconObject
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *HoconObject
	}{
		{
			name: "merges with nil correctly",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
				keys:  getArrayOfTwoSimpleKeys(),
			},
			want: wrapAllInObject(
				getArrayOfTwoSimpleKeys(),
				[]HoconElement{simpleLiteral1, simpleLiteral2}),
		},
		{
			name: "merges with empty object correctly",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
				keys:  getArrayOfTwoSimpleKeys(),
			},
			args: args{other: makeHoconObject([]string{}, []string{})},
			want: wrapAllInObject(
				getArrayOfTwoSimpleKeys(),
				[]HoconElement{simpleLiteral1, simpleLiteral2}),
		},
		{
			name: "merges with other correctly",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
				keys:  getArrayOfTwoSimpleKeys(),
			},
			args: args{other: wrapInObject(simpleKey3, simpleLiteral3)},
			want: wrapAllInObject(
				[]string{simpleKey1, simpleKey2, simpleKey3},
				[]HoconElement{simpleLiteral1, simpleLiteral2, simpleLiteral3}),
		},
		{
			name: "object with nested objects merges with other correctly",
			fields: fields{
				items: map[string]*HoconValue{
					simpleKey1: wrapInValue(wrapInObject(simpleKey2, simpleLiteral2)),
				},
				keys: []string{simpleKey1},
			},
			args: args{
				other: wrapInObject(simpleKey1, wrapInValue(wrapInObject(simpleKey3, simpleLiteral3))),
			},
			want: &HoconObject{
				keys: []string{simpleKey1},
				items: map[string]*HoconValue{
					simpleKey1: wrapInValue(wrapAllInObject(
						[]string{simpleKey2, simpleKey3},
						[]HoconElement{simpleLiteral2, simpleLiteral3})),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			got := p.MergeImmutable(tt.args.other)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MergeImmutable() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconObject_String(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "empty object returns empty string",
			want: "",
		},
		{
			name: "returns its fields with values",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
				keys:  getArrayOfTwoSimpleKeys(),
			},
			want: simpleKey1 + " : " + simpleLiteral1.value + newLine +
				simpleKey2 + " : " + simpleLiteral2.value + newLine,
		},
		{
			name: "returns nested objects in brackets",
			fields: fields{
				items: map[string]*HoconValue{
					simpleKey1: wrapInValue(wrapInObject(simpleKey2, simpleLiteral2)),
				},
				keys: []string{simpleKey1},
			},
			want: simpleKey1 + " : {" + newLine +
				"  " + simpleKey2 + " : " + simpleLiteral2.value + newLine +
				"}" + newLine,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			if got := p.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconObject_ToString(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
	}
	type args struct {
		indent int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "empty object returns empty string",
			fields: fields{},
			want:   "",
		},
		{
			name:   "empty object with non 0 indent returns empty string",
			fields: fields{},
			args:   args{indent: 1},
			want:   "",
		},
		{
			name: "returns its fields with values",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
				keys:  getArrayOfTwoSimpleKeys(),
			},
			want: simpleKey1 + " : " + simpleLiteral1.value + newLine +
				simpleKey2 + " : " + simpleLiteral2.value + newLine,
		},
		{
			name: "object with non 0 indent returns its fields with values",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
				keys:  getArrayOfTwoSimpleKeys(),
			},
			args: args{indent: 2},
			want: "    " + simpleKey1 + " : " + simpleLiteral1.value + newLine +
				"    " + simpleKey2 + " : " + simpleLiteral2.value + newLine,
		},
		{
			name: "returns nested objects in brackets",
			fields: fields{
				items: map[string]*HoconValue{
					simpleKey1: wrapInValue(wrapInObject(simpleKey2, simpleLiteral1)),
				},
				keys: []string{simpleKey1},
			},
			want: simpleKey1 + " : {" + newLine +
				"  " + simpleKey2 + " : " + simpleLiteral1.value + newLine +
				"}" + newLine,
		},
		{
			name: "object with non 0 indent returns nested objects in brackets",
			fields: fields{
				items: map[string]*HoconValue{
					simpleKey1: wrapInValue(wrapInObject(simpleKey2, simpleLiteral1)),
				},
				keys: []string{simpleKey1},
			},
			args: args{indent: 3},
			want: "      " + simpleKey1 + " : {" + newLine +
				"        " + simpleKey2 + " : " + simpleLiteral1.value + newLine +
				"      }" + newLine,
		},
	}
	wrapInObject(simpleKey1, wrapInObject(simpleKey2, simpleLiteral1))
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			got := p.ToString(tt.args.indent)
			if got != tt.want {
				t.Errorf("ToString() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconObject_Unwrapped(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:    "empty object cannot be unwrapped",
			fields:  fields{},
			wantErr: true,
		},
		{
			name: "returns its items",
			fields: fields{
				items: getMapOfTwoSimpleLiterals(),
				keys:  getArrayOfTwoSimpleKeys(),
			},
			want: map[string]interface{}{
				simpleKey1: wrapInValue(simpleLiteral1),
				simpleKey2: wrapInValue(simpleLiteral2),
			},
		},
		{
			name: "returns its item with nested object unwrapped",
			fields: fields{
				items: map[string]*HoconValue{
					simpleKey1: wrapInValue(wrapInObject(simpleKey2, simpleLiteral2)),
				},
				keys: []string{simpleKey1},
			},
			want: map[string]interface{}{
				simpleKey1: map[string]interface{}{simpleKey2: wrapInValue(simpleLiteral2)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			got, err := p.Unwrapped()
			if (err != nil) != tt.wantErr {
				t.Errorf("Unwrapped() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Unwrapped() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHoconObject(t *testing.T) {
	tests := []struct {
		name string
		want *HoconObject
	}{
		{
			name: "returns object with empty items and null keys",
			want: &HoconObject{items: map[string]*HoconValue{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHoconObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHoconObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
