package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Depado/quokka/utils"
)

// NewQuokkaTemplate will create a new Quokka template with default params
func NewQuokkaTemplate(path, name, description, version string, yes, debug bool) {
	var err error
	var fd *os.File

	if !utils.ConfirmFileExists(path, true, yes, debug) {
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			utils.FatalPrintln("Unable to create directory")
		}
	}
	qf := filepath.Join(path, ".quokka.yml")
	utils.ConfirmFileExists(qf, false, yes, debug)

	utils.AskIfEmptyString(&name, "name", "Template name?", "Quokka Template", debug)
	utils.AskIfEmptyString(&description, "description", "Template description?", "New Quokka Template", debug)
	utils.AskIfEmptyString(&version, "version", "Template version?", "0.1.0", debug)

	if fd, err = os.Create(qf); err != nil {
		utils.FatalPrintln("Unable to create file")
	}
	defer fd.Close()
	if _, err = fd.WriteString(fmt.Sprintf("name: \"%s\"\ndescription: \"%s\"\nversion: \"%s\"\n", name, description, version)); err != nil {
		utils.FatalPrintln("Unable to write in file")
	}
}
