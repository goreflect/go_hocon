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
		{name: "empty", fields: fields{}, want: nil, wantErr: true},
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
			name: "simple",
			fields: fields{
				items: map[string]*HoconValue{"a": {values: []HoconElement{NewHoconLiteral("b")}}},
				keys:  []string{"a"},
			},
			args: args{key: "a"},
			want: &HoconValue{values: []HoconElement{NewHoconLiteral("b")}},
		},
		{
			name:   "unknown",
			fields: fields{},
			args:   args{key: "a"},
			want:   nil,
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
			name: "simple",
			fields: fields{
				keys: []string{"a", "c", "d"},
			},
			want: []string{"a", "c", "d"},
		},
		{
			name:   "unknown",
			fields: fields{},
			want:   nil,
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
			name: "simple",
			fields: fields{
				items: map[string]*HoconValue{"a": {values: []HoconElement{NewHoconLiteral("b")}}},
				keys:  []string{"a"},
			},
			args: args{key: "a"},
			want: &HoconValue{oldValue: &HoconValue{values: []HoconElement{NewHoconLiteral("b")}}},
		},
		{
			name: "existedMap",
			fields: fields{
				items: map[string]*HoconValue{"a": {values: []HoconElement{NewHoconLiteral("b")}}},
				keys:  []string{"a"},
			},
			args: args{key: "b"},
			want: &HoconValue{},
		},
		{
			name:   "noMapNewValue",
			fields: fields{},
			args:   args{key: "a"},
			want:   &HoconValue{},
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		// TODO: Add test cases.
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
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			if err := p.Merge(tt.args.other); (err != nil) != tt.wantErr {
				t.Errorf("Merge() error = %v, wantErr %v", err, tt.wantErr)
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
		name    string
		fields  fields
		args    args
		want    *HoconObject
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			got, err := p.MergeImmutable(tt.args.other)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeImmutable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		// TODO: Add test cases.
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
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			got, err := p.ToString(tt.args.indent)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
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
		// TODO: Add test cases.
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

func TestHoconObject_quoteIfNeeded(t *testing.T) {
	type fields struct {
		items map[string]*HoconValue
		keys  []string
	}
	type args struct {
		text string
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
			p := &HoconObject{
				items: tt.fields.items,
				keys:  tt.fields.keys,
			}
			if got := p.quoteIfNeeded(tt.args.text); got != tt.want {
				t.Errorf("quoteIfNeeded() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHoconObject(t *testing.T) {
	tests := []struct {
		name string
		want *HoconObject
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHoconObject(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHoconObject() = %v, want %v", got, tt.want)
			}
		})
	}
}
