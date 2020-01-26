package hocon

import (
	"reflect"
	"testing"
)

func TestNewToken(t *testing.T) {
	type args struct {
		v interface{}
	}
	tests := []struct {
		name string
		args args
		want *Token
	}{
		{
			name: "returns nil with empty args",
			args: args{},
			want: nil,
		},
		{
			name: "returns nil with unknown type",
			args: args{
				v: 123,
			},
			want: nil,
		},
		{
			name: "returns token type literal with string",
			args: args{
				v: simpleKey1,
			},
			want: &Token{tokenType: TokenTypeLiteralValue, value: simpleKey1},
		},
		{
			name: "returns token type none correctly",
			args: args{
				v: TokenTypeNone,
			},
			want: &Token{tokenType: TokenTypeNone},
		},
		{
			name: "returns token type dot correctly",
			args: args{
				v: TokenTypeDot,
			},
			want: &Token{tokenType: TokenTypeDot},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewToken(tt.args.v); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTokenInclude(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want *Token
	}{
		{
			name: "creates simple token with no arguments",
			want: &Token{tokenType: TokenTypeInclude},
		},
		{
			name: "creates simple token",
			args: args{path: simpleKey1},
			want: &Token{tokenType: TokenTypeInclude, value: simpleKey1},
		},
		{
			name: "creates token with specials",
			args: args{path: specials},
			want: &Token{tokenType: TokenTypeInclude, value: specials},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTokenInclude(tt.args.path); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenInclude() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTokenKey(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want *Token
	}{
		{
			name: "creates simple token with no arguments",
			want: &Token{tokenType: TokenTypeKey},
		},
		{
			name: "creates simple token",
			args: args{key: simpleKey1},
			want: &Token{tokenType: TokenTypeKey, value: simpleKey1},
		},
		{
			name: "creates token with specials",
			args: args{key: specials},
			want: &Token{tokenType: TokenTypeKey, value: specials},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTokenKey(tt.args.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTokenLiteralValue(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name string
		args args
		want *Token
	}{
		{
			name: "creates simple token with no arguments",
			want: &Token{tokenType: TokenTypeLiteralValue},
		},
		{
			name: "creates simple token",
			args: args{value: simpleKey1},
			want: &Token{tokenType: TokenTypeLiteralValue, value: simpleKey1},
		},
		{
			name: "creates token with specials",
			args: args{value: specials},
			want: &Token{tokenType: TokenTypeLiteralValue, value: specials},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTokenLiteralValue(tt.args.value); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenLiteralValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewTokenSubstitution(t *testing.T) {
	type args struct {
		path       string
		isOptional bool
	}
	tests := []struct {
		name string
		args args
		want *Token
	}{
		{
			name: "creates simple token with no arguments",
			want: &Token{tokenType: TokenTypeSubstitute},
		},
		{
			name: "creates simple token",
			args: args{
				path: simpleKey1,
			},
			want: &Token{tokenType: TokenTypeSubstitute, value: simpleKey1},
		},
		{
			name: "creates token with specials",
			args: args{
				path: specials,
			},
			want: &Token{tokenType: TokenTypeSubstitute, value: specials},
		},
		{
			name: "creates token with specials",
			args: args{
				path:       specials,
				isOptional: true,
			},
			want: &Token{tokenType: TokenTypeSubstitute, isOptional: true, value: specials},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTokenSubstitution(tt.args.path, tt.args.isOptional); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTokenSubstitution() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringTokenType(t *testing.T) {
	type args struct {
		tokenType TokenType
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "returns TokenTypeNone with no arguments",
			want: "TokenTypeNone",
		},
		{
			name: "returns TokenTypeNone",
			args: args{
				tokenType: TokenTypeNone,
			},
			want: "TokenTypeNone",
		},
		{
			name: "returns TokenTypeComment",
			args: args{
				tokenType: TokenTypeComment,
			},
			want: "TokenTypeComment",
		},
		{
			name: "returns TokenTypeKey",
			args: args{
				tokenType: TokenTypeKey,
			},
			want: "TokenTypeKey",
		},
		{
			name: "returns TokenTypeLiteralValue",
			args: args{
				tokenType: TokenTypeLiteralValue,
			},
			want: "TokenTypeLiteralValue",
		},
		{
			name: "returns TokenTypeAssign",
			args: args{
				tokenType: TokenTypeAssign,
			},
			want: "TokenTypeAssign",
		},
		{
			name: "returns TokenTypePlusAssign",
			args: args{
				tokenType: TokenTypePlusAssign,
			},
			want: "TokenTypePlusAssign",
		},
		{
			name: "returns TokenTypeObjectStart",
			args: args{
				tokenType: TokenTypeObjectStart,
			},
			want: "TokenTypeObjectStart",
		},
		{
			name: "returns TokenTypeObjectEnd",
			args: args{
				tokenType: TokenTypeObjectEnd,
			},
			want: "TokenTypeObjectEnd",
		},
		{
			name: "returns TokenTypeDot",
			args: args{
				tokenType: TokenTypeDot,
			},
			want: "TokenTypeDot",
		},
		{
			name: "returns TokenTypeNewline",
			args: args{
				tokenType: TokenTypeNewline,
			},
			want: "TokenTypeNewline",
		},
		{
			name: "returns TokenTypeEoF",
			args: args{
				tokenType: TokenTypeEoF,
			},
			want: "TokenTypeEoF",
		},
		{
			name: "returns TokenTypeArrayStart",
			args: args{
				tokenType: TokenTypeArrayStart,
			},
			want: "TokenTypeArrayStart",
		},
		{
			name: "returns TokenTypeArrayEnd",
			args: args{
				tokenType: TokenTypeArrayEnd,
			},
			want: "TokenTypeArrayEnd",
		},
		{
			name: "returns TokenTypeComma",
			args: args{
				tokenType: TokenTypeComma,
			},
			want: "TokenTypeComma",
		},
		{
			name: "returns TokenTypeSubstitute",
			args: args{
				tokenType: TokenTypeSubstitute,
			},
			want: "TokenTypeSubstitute",
		},
		{
			name: "returns TokenTypeInclude",
			args: args{
				tokenType: TokenTypeInclude,
			},
			want: "TokenTypeInclude",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringTokenType(tt.args.tokenType); got != tt.want {
				t.Errorf("StringTokenType() = %v, want %v", got, tt.want)
			}
		})
	}
}
