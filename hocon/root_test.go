package hocon

import (
	"reflect"
	"testing"
)

func TestHoconRoot_Substitutions(t *testing.T) {
	type fields struct {
		value         *HoconValue
		substitutions []*HoconSubstitution
	}
	tests := []struct {
		name   string
		fields fields
		want   []*HoconSubstitution
	}{
		{
			name: "returns nil if no substitutions",
			want: nil,
		},
		{
			name: "returns substitution correctly",
			fields: fields{
				value:         nil,
				substitutions: []*HoconSubstitution{getCycledSubstitution()},
			},
			want: []*HoconSubstitution{getCycledSubstitution()},
		},
		{
			name: "returns substitutions correctly",
			fields: fields{
				value:         nil,
				substitutions: []*HoconSubstitution{getCycledSubstitution(), getCycledSubstitution()},
			},
			want: []*HoconSubstitution{getCycledSubstitution(), getCycledSubstitution()},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := HoconRoot{
				value:         tt.fields.value,
				substitutions: tt.fields.substitutions,
			}
			if got := p.Substitutions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Substitutions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHoconRoot_Value(t *testing.T) {
	type fields struct {
		value         *HoconValue
		substitutions []*HoconSubstitution
	}
	tests := []struct {
		name   string
		fields fields
		want   *HoconValue
	}{
		{
			name: "returns nil if no value",
			want: nil,
		},
		{
			name: "returns value correctly",
			fields: fields{
				value:         &HoconValue{values: []HoconElement{simpleLiteral}},
				substitutions: nil,
			},
			want: &HoconValue{values: []HoconElement{simpleLiteral}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := HoconRoot{
				value:         tt.fields.value,
				substitutions: tt.fields.substitutions,
			}
			if got := p.Value(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Value() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewHoconRoot(t *testing.T) {
	type args struct {
		value         *HoconValue
		substitutions []*HoconSubstitution
	}
	tests := []struct {
		name string
		args args
		want *HoconRoot
	}{
		{
			name: "returns empty hocon root",
			want: &HoconRoot{},
		},
		{
			name: "returns value correctly",
			args: args{
				value: &HoconValue{values: []HoconElement{simpleLiteral}},
			},
			want: &HoconRoot{value: &HoconValue{values: []HoconElement{simpleLiteral}}},
		},
		{
			name: "returns substitutions correctly",
			args: args{
				substitutions: []*HoconSubstitution{getCycledSubstitution(), getCycledSubstitution()},
			},
			want: &HoconRoot{substitutions: []*HoconSubstitution{getCycledSubstitution(), getCycledSubstitution()}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewHoconRoot(tt.args.value, tt.args.substitutions...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHoconRoot() = %v, want %v", got, tt.want)
			}
		})
	}
}
