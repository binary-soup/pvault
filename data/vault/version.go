package vault

import (
	"fmt"
	"os"
	"path/filepath"
	"pvault/data/version"

	"github.com/binary-soup/go-commando/alert"
)

const VERSION_FILE = "VERSION.txt"

func (v Vault) writeVersion() error {
	err := os.WriteFile(filepath.Join(v.Path, VERSION_FILE), []byte(fmt.Sprintf("Version: %d", version.VAULT)), 0666)
	if err != nil {
		return alert.ChainError(err, "error writing vault version file")
	}
	return nil
}

func (v Vault) ReadVersion() (uint16, error) {
	file, err := os.Open(filepath.Join(v.Path, VERSION_FILE))
	if err != nil {
		return 0, alert.ChainError(err, "error opening vault version file")
	}
	defer file.Close()

	var version uint16
	fmt.Fscanf(file, "Version: %d", &version)

	return version, nil
}
