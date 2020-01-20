package hocon

import (
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	_Num1000 = big.NewInt(1000)
	_Num1024 = big.NewInt(1024)
)

var (
	_IByte = big.NewInt(1)
	_KByte = (&big.Int{}).Mul(_IByte, _Num1000)
	_MByte = (&big.Int{}).Mul(_KByte, _Num1000)
	_GByte = (&big.Int{}).Mul(_MByte, _Num1000)
	_TByte = (&big.Int{}).Mul(_GByte, _Num1000)
	_PByte = (&big.Int{}).Mul(_TByte, _Num1000)
	_EByte = (&big.Int{}).Mul(_PByte, _Num1000)
	_ZByte = (&big.Int{}).Mul(_EByte, _Num1000)
	_YByte = (&big.Int{}).Mul(_ZByte, _Num1000)
)

var (
	_Byte   = big.NewInt(1)
	_KiByte = (&big.Int{}).Mul(_Byte, _Num1024)
	_MiByte = (&big.Int{}).Mul(_KiByte, _Num1024)
	_GiByte = (&big.Int{}).Mul(_MiByte, _Num1024)
	_TiByte = (&big.Int{}).Mul(_GiByte, _Num1024)
	_PiByte = (&big.Int{}).Mul(_TiByte, _Num1024)
	_EiByte = (&big.Int{}).Mul(_PiByte, _Num1024)
	_ZiByte = (&big.Int{}).Mul(_EiByte, _Num1024)
	_YiByte = (&big.Int{}).Mul(_ZiByte, _Num1024)
)

type HoconValue struct {
	values   []HoconElement
	oldValue *HoconValue
}

func NewHoconValue() *HoconValue {
	return &HoconValue{}
}

func (p *HoconValue) IsEmpty() bool {
	if len(p.values) == 0 {
		return true
	}

	if first, ok := p.values[0].(*HoconObject); ok {
		if len(first.items) == 0 {
			return true
		}
	}
	return false
}

func (p *HoconValue) AtKey(key string) *HoconRoot {
	obj := NewHoconObject()
	obj.GetOrCreateKey(key)
	obj.items[key] = p
	r := NewHoconValue()
	r.AppendValue(obj)
	return NewHoconRoot(r)
}

func (p *HoconValue) IsString() bool {
	strCount := 0

	for _, v := range p.values {
		v = p.topValueOfSub(v)
		if v.IsString() {
			strCount += 1
		}
	}

	if strCount > 0 && strCount == len(p.values) {
		return true
	}

	return false
}

func (p *HoconValue) topValueOfSub(v interface{}) HoconElement {
	if v == nil {
		return nil
	}

	if sub, ok := v.(*HoconSubstitution); ok {
		if sub.ResolvedValue != nil && sub.ResolvedValue.oldValue != nil {
			return sub.ResolvedValue.oldValue
		}
		if sub.ResolvedValue == nil && p.oldValue != nil {
			return p.oldValue
		}
	}

	return v.(HoconElement)
}

func (p *HoconValue) concatString() (string, error) {
	var concat string
	for _, v := range p.values {
		v = p.topValueOfSub(v)
		stringV, err := v.GetString()
		if err != nil {
			return "", err
		}

		concat += stringV
	}

	if concat == "null" {
		concat = ""
	}

	return strings.TrimSpace(concat), nil
}

func (p *HoconValue) GetByteSize() (*big.Int, error) {
	res, err := p.GetString()
	if err != nil {
		return nil, err
	}

	groups, matched := findStringSubmatchMap(res, `^(?P<value>([0-9]+(\.[0-9]+)?))\s*(?P<unit>(B|b|byte|bytes|kB|kilobyte|kilobytes|MB|megabyte|megabytes|GB|gigabyte|gigabytes|TB|terabyte|terabytes|PB|petabyte|petabytes|EB|exabyte|exabytes|ZB|zettabyte|zettabytes|YB|yottabyte|yottabytes|K|k|Ki|KiB|kibibyte|kibibytes|M|m|Mi|MiB|mebibyte|mebibytes|G|g|Gi|GiB|gibibyte|gibibytes|T|t|Ti|TiB|tebibyte|tebibytes|P|p|Pi|PiB|pebibyte|pebibytes|E|e|Ei|EiB|exbibyte|exbibytes|Z|z|Zi|ZiB|zebibyte|zebibytes|Y|y|Yi|YiB|yobibyte|yobibytes))$`)

	if matched {
		foundUnit, foundFloat := groups["unit"], groups["value"]
		positiveV, err := parsePositiveValue(foundFloat)
		if err != nil {
			return nil, err
		}

		intV := int64(positiveV) // todo 1.5 TB is not going to work due to floor to 1
		bigInt := big.NewInt(intV)

		switch foundUnit {
		case "B", "b", "byte", "bytes":
			return (&big.Int{}).Mul(bigInt, _IByte), nil
		case "kB", "kilobyte", "kilobytes":
			return (&big.Int{}).Mul(bigInt, _KByte), nil
		case "MB", "megabyte", "megabytes":
			return (&big.Int{}).Mul(bigInt, _MByte), nil
		case "GB", "gigabyte", "gigabytes":
			return (&big.Int{}).Mul(bigInt, _GByte), nil
		case "TB", "terabyte", "terabytes":
			return (&big.Int{}).Mul(bigInt, _TByte), nil
		case "PB", "petabyte", "petabytes":
			return (&big.Int{}).Mul(bigInt, _PByte), nil
		case "EB", "exabyte", "exabytes":
			return (&big.Int{}).Mul(bigInt, _EByte), nil
		case "ZB", "zettabyte", "zettabytes":
			return (&big.Int{}).Mul(bigInt, _ZByte), nil
		case "YB", "yottabyte", "yottabytes":
			return (&big.Int{}).Mul(bigInt, _YByte), nil
		case "K", "k", "Ki", "KiB", "kibibyte", "kibibytes":
			return (&big.Int{}).Mul(bigInt, _KiByte), nil
		case "M", "m", "Mi", "MiB", "mebibyte", "mebibytes":
			return (&big.Int{}).Mul(bigInt, _MiByte), nil
		case "G", "g", "Gi", "GiB", "gibibyte", "gibibytes":
			return (&big.Int{}).Mul(bigInt, _GiByte), nil
		case "T", "t", "Ti", "TiB", "tebibyte", "tebibytes":
			return (&big.Int{}).Mul(bigInt, _TiByte), nil
		case "P", "p", "Pi", "PiB", "pebibyte", "pebibytes":
			return (&big.Int{}).Mul(bigInt, _PiByte), nil
		case "E", "e", "Ei", "EiB", "exbibyte", "exbibytes":
			return (&big.Int{}).Mul(bigInt, _EiByte), nil
		case "Z", "z", "Zi", "ZiB", "zebibyte", "zebibytes":
			return (&big.Int{}).Mul(bigInt, _ZiByte), nil
		case "Y", "y", "Yi", "YiB", "yobibyte", "yobibytes":
			return (&big.Int{}).Mul(bigInt, _YiByte), nil
		}
	}

	return nil, errors.New("unknown byte size unit")
}

func (p *HoconValue) String() string {
	stringV, err := p.ToString(0)
	if err != nil {
		return fmt.Sprintf("cannot get string: %s", err.Error())
	}

	return stringV
}

func (p *HoconValue) ToString(indent int) (string, error) {
	if p.IsString() {
		stringV, err := p.GetString()
		if err != nil {
			return "", err
		}

		return p.quoteIfNeeded(stringV), nil
	}

	isObject, err := p.IsObject()
	if err != nil {
		return "", err
	}

	if isObject {
		indentString := strings.Repeat(" ", indent*2)
		objectV, err := p.GetObject()
		if err != nil {
			return "", err
		}

		stringV, err := objectV.ToString(indent + 1)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("{%s%s%s}", newLine, stringV, indentString), nil
	}

	if p.IsArray() {
		arrayV, err := p.GetArray()
		if err != nil {
			return "", err
		}

		var sstr []string
		for _, v := range arrayV {
			stringV, err := v.ToString(indent + 1)
			if err != nil {
				return "", err
			}

			sstr = append(sstr, stringV)
		}

		return "[" + strings.Join(sstr, ",") + "]", nil
	}

	if p.IsEmpty() {
		return "", nil
	}

	return "<<unknown value>>", nil
}

func (p *HoconValue) GetObject() (*HoconObject, error) {
	if len(p.values) == 0 {
		return nil, nil
	}

	raw := p.values[0]
	if o, ok := raw.(*HoconObject); ok {
		return o, nil
	}

	raw = p.topValueOfSub(raw)

	if s, ok := raw.(*HoconSubstitution); ok {
		if s.ResolvedValue == nil {
			return nil, nil
		}
	}

	if sub, ok := raw.(MightBeAHoconObject); ok {
		if sub != nil {
			isObject, err := sub.IsObject()
			if err != nil {
				return nil, err
			}

			if isObject {
				return sub.GetObject()
			}
		}
	}

	return nil, nil
}

func (p *HoconValue) IsObject() (bool, error) {
	objectV, err := p.GetObject()
	if err != nil {
		return false, err
	}

	return objectV != nil, nil
}

func (p *HoconValue) AppendValue(value HoconElement) {
	p.values = append(p.values, value)
}

func (p *HoconValue) Clear() {
	p.values = []HoconElement{}
}

func (p *HoconValue) NewValue(value HoconElement) {
	p.values = []HoconElement{}
	p.values = append(p.values, value)
}

func (p *HoconValue) GetBoolean() (bool, error) {
	stringV, err := p.GetString()
	if err != nil {
		return false, err
	}

	switch strings.ToLower(stringV) {
	case "on", "true", "yes":
		return true, nil
	case "off", "false", "no":
		return false, nil
	}
	return false, fmt.Errorf("unknown boolean format: %s", stringV)
}

func (p *HoconValue) GetString() (string, error) {
	if p.IsString() {
		return p.concatString()
	}
	return "", nil
}

func (p *HoconValue) GetFloat64() (float64, error) {
	stringV, err := p.GetString()
	if err != nil {
		return 0, err
	}

	return strconv.ParseFloat(stringV, 64)
}

func (p *HoconValue) GetFloat32() (float32, error) {
	stringV, err := p.GetString()
	if err != nil {
		return 0, err
	}

	floatV, err := strconv.ParseFloat(stringV, 32)
	if err != nil {
		return 0, err
	}

	return float32(floatV), nil
}

func (p *HoconValue) GetInt64() (int64, error) {
	stringV, err := p.GetString()
	if err != nil {
		return 0, nil
	}

	return strconv.ParseInt(stringV, 10, 64)
}

func (p *HoconValue) GetInt32() (int32, error) {
	stringV, err := p.GetString()
	if err != nil {
		return 0, nil
	}

	intV, err := strconv.ParseInt(stringV, 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(intV), nil
}

func (p *HoconValue) GetByte() (byte, error) {
	stringV, err := p.GetString()
	if err != nil {
		return 0, nil
	}

	intV, err := strconv.ParseInt(stringV, 10, 8)
	if err != nil {
		return 0, err
	}

	return byte(intV), nil
}

func (p *HoconValue) GetByteList() ([]byte, error) {
	arrayV, err := p.GetArray()
	if err != nil {
		return nil, err
	}

	var items []byte
	for _, v := range arrayV {
		byteV, err := v.GetByte()
		if err != nil {
			return nil, err
		}

		items = append(items, byteV)
	}

	return items, nil
}

func (p *HoconValue) GetInt32List() ([]int32, error) {
	arrayV, err := p.GetArray()
	if err != nil {
		return nil, err
	}

	var items []int32
	for _, v := range arrayV {
		intV, err := v.GetInt32()
		if err != nil {
			return nil, err
		}

		items = append(items, intV)
	}

	return items, nil
}

func (p *HoconValue) GetInt64List() ([]int64, error) {
	arrayV, err := p.GetArray()
	if err != nil {
		return nil, err
	}

	var items []int64
	for _, v := range arrayV {
		intV, err := v.GetInt64()
		if err != nil {
			return nil, nil
		}

		items = append(items, intV)
	}

	return items, nil
}

func (p *HoconValue) GetBooleanList() ([]bool, error) {
	arrayV, err := p.GetArray()
	if err != nil {
		return nil, err
	}

	var items []bool
	for _, v := range arrayV {
		booleanV, err := v.GetBoolean()
		if err != nil {
			return nil, err
		}

		items = append(items, booleanV)
	}

	return items, nil
}

func (p *HoconValue) GetFloat32List() ([]float32, error) {
	arrayV, err := p.GetArray()
	if err != nil {
		return nil, err
	}

	var items []float32
	for _, v := range arrayV {
		floatV, err := v.GetFloat32()
		if err != nil {
			return nil, err
		}

		items = append(items, floatV)
	}

	return items, nil
}

func (p *HoconValue) GetFloat64List() ([]float64, error) {
	arrayV, err := p.GetArray()
	if err != nil {
		return nil, err
	}

	var items []float64
	for _, v := range arrayV {
		floatV, err := v.GetFloat64()
		if err != nil {
			return nil, err
		}

		items = append(items, floatV)
	}

	return items, nil
}

func (p *HoconValue) GetStringList() ([]string, error) {
	arrayV, err := p.GetArray()
	if err != nil {
		return nil, err
	}

	var items []string
	for _, v := range arrayV {
		stringV, err := v.GetString()
		if err != nil {
			return nil, err
		}

		items = append(items, stringV)
	}

	return items, nil
}

func (p *HoconValue) GetArray() ([]*HoconValue, error) {
	var items []*HoconValue

	if len(p.values) == 0 {
		return items, nil
	}

	for _, v := range p.values {
		v = p.topValueOfSub(v)
		if v.IsArray() {
			arrayV, err := v.GetArray()
			if err != nil {
				return nil, err
			}

			items = append(items, arrayV...)
		}
	}

	return items, nil
}

func (p *HoconValue) GetChildObject(key string) (*HoconValue, error) {
	objectV, err := p.GetObject()
	if err != nil {
		return nil, err
	}

	if objectV == nil {
		return nil, nil
	}

	return objectV.GetKey(key), nil
}

func (p *HoconValue) IsArray() bool {
	arr, err := p.GetArray()
	if err != nil {
		return false
	}

	return arr != nil
}

func (p *HoconValue) GetTimeDuration(allowInfinite bool) (time.Duration, error) {
	stringV, err := p.GetString()
	if err != nil {
		return 0, err
	}

	groups, matched := findStringSubmatchMap(stringV, `^(?P<value>([0-9]+(\.[0-9]+)?))\s*(?P<unit>(nanoseconds|nanosecond|nanos|nano|ns|microseconds|microsecond|micros|micro|us|milliseconds|millisecond|millis|milli|ms|seconds|second|s|minutes|minute|m|hours|hour|h|days|day|d))$`)

	if matched {
		foundUnit, foundFloat := groups["unit"], groups["value"]
		floatV, err := parsePositiveValue(foundFloat)
		if err != nil {
			return 0, err
		}

		switch foundUnit {
		case "nanoseconds", "nanosecond", "nanos", "nano", "ns":
			return time.Duration(float64(time.Nanosecond) * floatV), nil
		case "microseconds", "microsecond", "micros", "micro":
			return time.Duration(float64(time.Microsecond) * floatV), nil
		case "milliseconds", "millisecond", "millis", "milli", "ms":
			return time.Duration(float64(time.Millisecond) * floatV), nil
		case "seconds", "second", "s":
			return time.Duration(float64(time.Second) * floatV), nil
		case "minutes", "minute", "m":
			return time.Duration(float64(time.Minute) * floatV), nil
		case "hours", "hour", "h":
			return time.Duration(float64(time.Hour) * floatV), nil
		case "days", "day", "d":
			return time.Duration(float64(time.Hour*24) * floatV), nil
		}

		return 0, fmt.Errorf("cannot parse time value: %s", stringV)
	}

	if strings.ToLower(stringV) == "infinite" {
		if allowInfinite {
			return time.Duration(-1), nil
		}
		return 0, errors.New("infinite time duration not allowed")
	}

	floatV, err := parsePositiveValue(stringV)
	if err != nil {
		return 0, err
	}

	return time.Duration(float64(time.Millisecond) * floatV), nil
}

func (p *HoconValue) quoteIfNeeded(text string) string {
	if len(text) == 0 {
		return `""`
	}

	if strings.IndexByte(text, ' ') >= 0 ||
		strings.IndexByte(text, '\t') >= 0 {
		return fmt.Sprintf(`"%s"`, text)
	}

	return text
}

func findStringSubmatchMap(s, exp string) (map[string]string, bool) {
	reg := regexp.MustCompile(exp)
	captures := make(map[string]string)

	match := reg.FindStringSubmatch(s)
	if match == nil {
		return captures, false
	}

	for i, name := range reg.SubexpNames() {
		if i == 0 || name == "" {
			continue
		}
		captures[name] = match[i]
	}
	return captures, true
}

func parsePositiveValue(v string) (float64, error) {
	value, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, err
	}

	if value < 0 {
		return 0, fmt.Errorf("expected a positive value instead of %s", v)
	}
	return value, nil
}
