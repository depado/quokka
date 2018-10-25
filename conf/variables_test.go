package conf

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	tbool     = true
	empty     = &Variable{Name: "empty"}
	hasvalue  = &Variable{Name: "value", Result: "result"}
	hasbool   = &Variable{Name: "bool", Confirm: &tbool}
	subvar    = &Variable{Name: "sub", Result: "ok"}
	parentvar = &Variable{Name: "parent", Variables: Variables{subvar}}
)

func TestVariables_Ctx(t *testing.T) {
	tests := []struct {
		name   string
		fields Variables
		want   map[string]interface{}
	}{
		{"should get one", Variables{empty},
			map[string]interface{}{empty.Name: ""}},
		{"should get two vars", Variables{empty, hasvalue},
			map[string]interface{}{empty.Name: "", hasvalue.Name: hasvalue.Result}},
		{"should get bool", Variables{hasbool},
			map[string]interface{}{hasbool.Name: true}},
		{"should get all vars", Variables{empty, hasvalue, hasbool},
			map[string]interface{}{empty.Name: "", hasvalue.Name: hasvalue.Result, hasbool.Name: true}},
		{"should get even sub vars", Variables{parentvar},
			map[string]interface{}{parentvar.Name: "", parentvar.Name + "_" + subvar.Name: "ok"}},
		{"should get even all with sub vars", Variables{empty, hasvalue, hasbool, parentvar},
			map[string]interface{}{empty.Name: "", hasvalue.Name: hasvalue.Result, hasbool.Name: true, parentvar.Name: "", parentvar.Name + "_" + subvar.Name: "ok"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.Ctx(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Variables.Ctx() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariables_AddToCtx(t *testing.T) {
	basectx := map[string]interface{}{"one": "one", "two": true}
	type args struct {
		prefix string
		ctx    map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  Variables
		args    args
		expects map[string]interface{}
	}{
		{"should add empty to context", Variables{empty}, args{"", basectx},
			map[string]interface{}{empty.Name: "", "one": "one", "two": true}},
		{"should add boolean to context", Variables{hasbool}, args{"", basectx},
			map[string]interface{}{hasbool.Name: true, "one": "one", "two": true}},
		{"should add boolean with key to context", Variables{hasbool}, args{"sub", basectx},
			map[string]interface{}{"sub_" + hasbool.Name: true, "one": "one", "two": true}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cp := make(map[string]interface{})
			for k, v := range tt.args.ctx {
				cp[k] = v
			}
			tt.fields.AddToCtx(tt.args.prefix, cp)
			assert.Equal(t, tt.expects, cp)
		})
	}
}

func TestVariables_FindNamed(t *testing.T) {
	vv := Variables{hasbool, hasvalue, empty}
	type args struct {
		s string
	}
	tests := []struct {
		name string
		vv   Variables
		args args
		want *Variable
	}{
		{"should find even empty", vv, args{hasvalue.Name}, hasvalue},
		{"should find bool", vv, args{hasbool.Name}, hasbool},
		{"should find value", vv, args{hasvalue.Name}, hasvalue},
		{"shouldn't find", vv, args{"random.jpg"}, nil},
		{"should find nested", Variables{parentvar}, args{"parent_sub"}, subvar},
		{"shouldn't find nested raw", Variables{parentvar}, args{"sub"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.vv.FindNamed(tt.args.s)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Variables.FindNamed() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariable_True(t *testing.T) {
	var fbool bool
	type fields struct {
		Default      string
		CustomPrompt string
		Values       []string
		Help         string
		Required     bool
		Confirm      *bool
		Variables    Variables
		Result       string
		Name         string
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{"should be false when confirm is false", fields{Confirm: &fbool}, false},
		{"should be false when result is empty", fields{Result: ""}, false},
		{"should be true when confirm is true", fields{Confirm: &tbool}, true},
		{"should be true when result is true", fields{Result: "toto"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Variable{
				Default:      tt.fields.Default,
				CustomPrompt: tt.fields.CustomPrompt,
				Values:       tt.fields.Values,
				Help:         tt.fields.Help,
				Required:     tt.fields.Required,
				Confirm:      tt.fields.Confirm,
				Variables:    tt.fields.Variables,
				Result:       tt.fields.Result,
				Name:         tt.fields.Name,
			}
			if got := v.True(); got != tt.want {
				t.Errorf("Variable.True() = %v, want %v", got, tt.want)
			}
		})
	}
}
