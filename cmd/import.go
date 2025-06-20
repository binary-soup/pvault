package cmd

import (
	"bufio"
	"fmt"
	"os"
	cmdstyle "pvault/cmd/style"
	"pvault/data"
	vw "pvault/workflows/vault"
	"strings"

	"github.com/binary-soup/go-command/command"
	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type ImportCommand struct {
	command.CommandBase
}

func NewImportCommand() ImportCommand {
	return ImportCommand{
		CommandBase: command.NewCommandBase("import", "import many passwords from CSV [name|password|username|url]. All items will use the same passkey"),
	}
}

func (cmd ImportCommand) Run(args []string) error {
	path := cmd.Flags.String("p", "", "path to the import CSV file")
	cmd.Flags.Parse(args)

	cfg, err := data.LoadConfig()
	if err != nil {
		return err
	}

	if *path == "" {
		return util.Error("(p)ath missing or invalid")
	}

	passwords, err := cmd.loadImportCSV(*path)
	if err != nil {
		return err
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)
	defer cfg.Vault.Close()

	var passkey string
	err = workflow.ChooseOrVerifyPasskey(&passkey)
	if err != nil {
		return err
	}

	for _, password := range passwords {
		fmt.Printf("%s -> ", cmdstyle.NAME_STYLE.Format(password.Name))

		err = cmd.savePassword(workflow, password, passkey)
		if err != nil {
			fmt.Printf("%s %s\n", style.BoldError.Format("[ERROR]"), err)
			continue
		}

		style.BoldInfo.Println("Imported")
	}

	return nil
}

func (cmd ImportCommand) savePassword(workflow vw.VaultWorkflow, password *data.Password, passkey string) error {
	if workflow.Vault.Index.NameExists(password.Name) {
		return util.Error(fmt.Sprintf("name \"%s\" already exists", password.Name))
	}

	return workflow.Encrypt(password, data.NewPasswordCache(passkey))
}

func (cmd ImportCommand) loadImportCSV(path string) ([]*data.Password, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, util.ChainError(err, "error opening import file")
	}
	defer file.Close()

	passwords := []*data.Password{}

	scanner := bufio.NewScanner(file)
	line := 0

	for scanner.Scan() {
		line++
		password := cmd.parsePasswordFromCSV(strings.TrimSpace(scanner.Text()))

		err = password.Validate()
		if err != nil {
			return nil, util.ChainErrorF(err, "[line %d] error validating password", line)
		}

		passwords = append(passwords, password)
	}

	if err := scanner.Err(); err != nil {
		return nil, util.ChainError(err, "error parsing CSV file")
	}

	return passwords, nil
}

func (cmd ImportCommand) parsePasswordFromCSV(line string) *data.Password {
	tokens := strings.SplitN(line, ",", 5)

	password := &data.Password{
		Name:          tokens[0],
		RecoveryCodes: []string{},
	}

	if len(tokens) > 1 {
		password.Password = tokens[1]
	}

	if len(tokens) > 2 {
		password.Username = tokens[2]
	}

	if len(tokens) > 3 {
		password.URL = tokens[3]
	}

	if len(tokens) > 4 {
		codes := strings.Split(tokens[4], ":")

		password.RecoveryCodes = make([]string, len(codes))
		copy(password.RecoveryCodes, codes)
	}

	return password
}
