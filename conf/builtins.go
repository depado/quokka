package conf

import (
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func DefaultBuiltins(output string, root *Root) map[string]interface{} {
	builtins := map[string]interface{}{
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

	if name, err := gitConfig("user.name"); err == nil && name != "" {
		builtins["git_user"] = name
	}
	if email, err := gitConfig("user.email"); err == nil && email != "" {
		builtins["git_email"] = email
	}

	if root != nil {
		builtins["template_name"] = root.Name
		builtins["template_version"] = root.Version
	}

	return builtins
}

func gitConfig(key string) (string, error) {
	cmd := exec.Command("git", "config", key)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out[:len(out)-1]), nil
}

