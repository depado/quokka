package conf

import (
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func DefaultBuiltins(output string, root *Root) map[string]any {
	builtins := map[string]any{
		"year":      strconv.Itoa(time.Now().Year()),
		"date":      time.Now().Format("2006-01-02"),
		"datetime":  time.Now().Format(time.RFC3339),
		"timestamp": strconv.FormatInt(time.Now().Unix(), 10),
		"output":    filepath.Base(output),
		"os":        runtime.GOOS,
		"arch":      runtime.GOARCH,
		"uuid":      uuid.New().String(),
	}

	if u, err := user.Current(); err == nil {
		builtins["username"] = u.Username
		builtins["home"] = u.HomeDir
	}

	if h, err := os.Hostname(); err == nil {
		builtins["hostname"] = h
	}

	if cwd, err := os.Getwd(); err == nil {
		builtins["cwd"] = cwd
	}

	if root != nil {
		builtins["template_name"] = root.Name
		builtins["template_version"] = root.Version
	}

	return builtins
}
