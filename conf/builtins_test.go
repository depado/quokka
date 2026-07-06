package conf

import (
	"os"
	"os/user"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultBuiltins(t *testing.T) {
	builtins := DefaultBuiltins("/some/path/myproject", nil)

	assert.NotEmpty(t, builtins["year"])
	assert.NotEmpty(t, builtins["date"])
	assert.NotEmpty(t, builtins["datetime"])
	assert.NotEmpty(t, builtins["timestamp"])
	assert.Equal(t, "myproject", builtins["output"])
	assert.Equal(t, runtime.GOOS, builtins["os"])
	assert.Equal(t, runtime.GOARCH, builtins["arch"])
	assert.NotEmpty(t, builtins["uuid"])

	if u, err := user.Current(); err == nil {
		assert.Equal(t, u.Username, builtins["username"])
		assert.Equal(t, u.HomeDir, builtins["home"])
	}

	if h, err := os.Hostname(); err == nil {
		assert.Equal(t, h, builtins["hostname"])
	}

	if cwd, err := os.Getwd(); err == nil {
		assert.Equal(t, cwd, builtins["cwd"])
	}

	assert.Nil(t, builtins["template_name"])
	assert.Nil(t, builtins["template_version"])
}

func TestDefaultBuiltins_WithRoot(t *testing.T) {
	root := &Root{Name: "my-template", Version: "1.0.0"}
	builtins := DefaultBuiltins("/output", root)

	assert.Equal(t, "my-template", builtins["template_name"])
	assert.Equal(t, "1.0.0", builtins["template_version"])
}

func TestResolveDefault_NoPrefix(t *testing.T) {
	v := &Variable{Default: "hello"}
	builtins := map[string]interface{}{"hello": "world"}
	resolveDefault(v, builtins)
	assert.Equal(t, "hello", v.Default)
}

func TestResolveDefault_WithPrefix(t *testing.T) {
	v := &Variable{Default: "$username"}
	builtins := map[string]interface{}{"username": "alice"}
	resolveDefault(v, builtins)
	assert.Equal(t, "alice", v.Default)
}

func TestResolveDefault_Unknown(t *testing.T) {
	v := &Variable{Default: "$unknown"}
	builtins := map[string]interface{}{"username": "alice"}
	resolveDefault(v, builtins)
	assert.Equal(t, "$unknown", v.Default)
}

func TestResolveDefault_NilBuiltins(t *testing.T) {
	v := &Variable{Default: "$username"}
	resolveDefault(v, nil)
	assert.Equal(t, "$username", v.Default)
}

func TestResolveDefault_IntValue(t *testing.T) {
	v := &Variable{Default: "$count"}
	builtins := map[string]interface{}{"count": 42}
	resolveDefault(v, builtins)
	assert.Equal(t, "42", v.Default)
}

func TestResolveDefault_Int64Value(t *testing.T) {
	v := &Variable{Default: "$ts"}
	builtins := map[string]interface{}{"ts": int64(1710000000)}
	resolveDefault(v, builtins)
	assert.Equal(t, "1710000000", v.Default)
}

func TestResolveDefault_Float64Value(t *testing.T) {
	v := &Variable{Default: "$pi"}
	builtins := map[string]interface{}{"pi": 3.14}
	resolveDefault(v, builtins)
	assert.Equal(t, "3.14", v.Default)
}

func TestResolveDefault_BoolValue(t *testing.T) {
	v := &Variable{Default: "$flag"}
	builtins := map[string]interface{}{"flag": true}
	resolveDefault(v, builtins)
	assert.Equal(t, "true", v.Default)
}

func TestTemplateFuncUUID(t *testing.T) {
	fn, ok := templateFuncMaps["uuid"].(func() string)
	assert.True(t, ok, "uuid should be a template function")

	u1 := fn()
	u2 := fn()
	assert.Len(t, u1, 36)
	assert.Len(t, u2, 36)
	assert.NotEmpty(t, u1)
	assert.NotEmpty(t, u2)
	assert.NotEqual(t, u1, u2, "each call should generate a new UUID")
}
