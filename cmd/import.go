package cmd

import (
	"bufio"
	"fmt"
	"os"
	"pvault/data/config"
	"pvault/data/password"
	"pvault/data/vault"
	"pvault/tools"
	vw "pvault/workflows/vault"
	"strings"

	"github.com/binary-soup/go-commando/alert"
	"github.com/binary-soup/go-commando/command"
	"github.com/binary-soup/go-commando/style"
)

type ImportCommand struct {
	command.ConfigCommandBase[config.Config]
}

func NewImportCommand() ImportCommand {
	return ImportCommand{
		ConfigCommandBase: command.NewConfigCommandBase[config.Config]("import", "import many passwords from CSV [name|password|username|url]. All items will use the same passkey"),
	}
}

func (cmd ImportCommand) Run(args []string) error {
	path := cmd.Flags.String("p", "", "path to the import CSV file")
	cmd.Flags.Parse(args)

	cfg, err := cmd.LoadConfig(DATA_DIR)
	if err != nil {
		return err
	}

	if *path == "" {
		return alert.Error("(p)ath missing or invalid")
	}

	workflow := vw.NewVaultWorkflow(cfg.Vault)
	defer cfg.Vault.Close()

	passwords, err := cmd.loadImportCSV(*path, cfg.Vault.Index)
	if err != nil {
		return err
	}

	style.Success.Println("CSV file is valid")

	var passkey string
	err = tools.PromptPasskey(&passkey)
	if err != nil {
		return err
	}

	for _, password := range passwords {
		fmt.Printf("%s -> ", NAME_STYLE.Format(password.Meta.Name))
		password.Meta.Passkey = passkey

		err = workflow.Encrypt(password)
		if err != nil {
			fmt.Printf("%s %s\n", style.BoldError.Format("[ERROR]"), err)
			continue
		}

		style.BoldInfo.Println("Imported")
	}

	return nil
}

func (cmd ImportCommand) loadImportCSV(path string, index *vault.Index) ([]*password.Cache, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, alert.ChainError(err, "error opening import file")
	}
	defer file.Close()

	passwords := []*password.Cache{}
	errors := []string{"Invalid CSV File"}

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		cache := cmd.parseLine(line)

		if index.HasName(cache.Meta.Name) {
			errors = append(errors, cmd.lineError(lineNum, alert.ErrorF("name \"%s\" already exists", cache.Meta.Name)))
			continue
		}

		err = cache.Password.Validate()
		if err != nil {
			errors = append(errors, cmd.lineError(lineNum, err))
			continue
		}

		passwords = append(passwords, cache)
	}

	if err := scanner.Err(); err != nil {
		return nil, alert.ChainError(err, "error parsing CSV file")
	}

	if len(errors) > 1 {
		return nil, alert.Error(strings.Join(errors, "\n  "))
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

func (cmd ImportCommand) lineError(line int, err error) string {
	return fmt.Sprintf("%s %s", style.Bolded.FormatF("[line %d]", line), err.Error())
}
