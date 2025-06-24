package cmd

import (
	"bufio"
	"fmt"
	"os"
	"pvault/data/password"
	"pvault/tools"
	vw "pvault/workflows/vault"
	"strings"

	"github.com/binary-soup/go-command/style"
	"github.com/binary-soup/go-command/util"
)

type ImportCommand struct {
	ConfigCommandBase
}

func NewImportCommand() ImportCommand {
	return ImportCommand{
		ConfigCommandBase: NewConfigCommandBase("import", "import many passwords from CSV [name|password|username|url]. All items will use the same passkey"),
	}
}

func (cmd ImportCommand) Run(args []string) error {
	path := cmd.Flags.String("p", "", "path to the import CSV file")
	cmd.Flags.Parse(args)

	cfg, err := cmd.LoadConfig()
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
	err = tools.PromptPasskey(&passkey)
	if err != nil {
		return err
	}

	for _, password := range passwords {
		fmt.Printf("%s -> ", NAME_STYLE.Format(password.Meta.Name))
		password.Meta.Passkey = passkey

		err = cmd.savePassword(workflow, password)
		if err != nil {
			fmt.Printf("%s %s\n", style.BoldError.Format("[ERROR]"), err)
			continue
		}

		style.BoldInfo.Println("Imported")
	}

	return nil
}

func (cmd ImportCommand) savePassword(workflow vw.VaultWorkflow, cache *password.Cache) error {
	if workflow.Vault.Index.HasName(cache.Meta.Name) {
		return util.Error(fmt.Sprintf("name \"%s\" already exists", cache.Meta.Name))
	}

	return workflow.Encrypt(cache)
}

func (cmd ImportCommand) loadImportCSV(path string) ([]*password.Cache, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, util.ChainError(err, "error opening import file")
	}
	defer file.Close()

	passwords := []*password.Cache{}

	scanner := bufio.NewScanner(file)
	line := 0

	for scanner.Scan() {
		line++
		cache := cmd.parseLine(strings.TrimSpace(scanner.Text()))

		err = cache.Password.Validate()
		if err != nil {
			return nil, util.ChainErrorF(err, "[line %d] error validating password", line)
		}

		passwords = append(passwords, cache)
	}

	if err := scanner.Err(); err != nil {
		return nil, util.ChainError(err, "error parsing CSV file")
	}

	return passwords, nil
}

func (cmd ImportCommand) parseLine(line string) *password.Cache {
	tokens := strings.SplitN(line, ",", 5)

	pswrd := &password.Password{
		RecoveryCodes: []string{},
	}

	if len(tokens) > 1 {
		pswrd.Password = tokens[1]
	}

	if len(tokens) > 2 {
		pswrd.Username = tokens[2]
	}

	if len(tokens) > 3 {
		pswrd.URL = tokens[3]
	}

	if len(tokens) > 4 {
		codes := strings.Split(tokens[4], ":")

		pswrd.RecoveryCodes = make([]string, len(codes))
		copy(pswrd.RecoveryCodes, codes)
	}

	return &password.Cache{
		Password: pswrd,
		Meta:     password.NewMeta(tokens[0], ""),
	}
}
