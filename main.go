package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	vault "github.com/gozaddy/secret.ly"
)

var (
	v vault.Vault
)

func init() {
	var err error
	v, err = vault.FileVault("secretly-cli", vault.FileVaultOptions{CreateNew: true})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

type setCommand struct {
	fs            *flag.FlagSet
	keyname       string
	keyvalue      string
	encryptionKey string
}

type getCommand struct {
	fs            *flag.FlagSet
	keyname       string
	encryptionKey string
}

type runner interface {
	init([]string) error
	run() error
	name() string
}

func newSetCommand() *setCommand {
	sc := &setCommand{
		fs: flag.NewFlagSet("set", flag.ContinueOnError),
	}

	sc.fs.StringVar(&sc.keyname, "name", "", "the name (keyname) of the secret to be stored")
	sc.fs.StringVar(&sc.keyvalue, "value", "", "secret to be stored")
	sc.fs.StringVar(&sc.encryptionKey, "key", "", "encryption key used for encrypting/decrypting secrets")

	return sc
}

func newGetCommand() *getCommand {
	gc := &getCommand{
		fs: flag.NewFlagSet("get", flag.ContinueOnError),
	}

	gc.fs.StringVar(&gc.keyname, "name", "", "the name (keyname) of the secret to be retrieved")
	gc.fs.StringVar(&gc.encryptionKey, "key", "", "encryption key used for encrypting/decrypting secrets")

	return gc
}

func (sc *setCommand) init(args []string) error {
	return sc.fs.Parse(args)
}

func (sc *setCommand) name() string {
	return sc.fs.Name()
}

func (sc *setCommand) run() error {

	if sc.keyname == "" || sc.keyvalue == "" || sc.encryptionKey == "" {
		return errors.New("keyname, keyvalue or encryptionKey should not be empty")
	}
	err := v.Set(sc.keyname, sc.keyvalue, sc.encryptionKey)
	if err != nil {
		return err
	}
	fmt.Println("Secret stored!")
	return nil

}

func (gc *getCommand) init(args []string) error {
	return gc.fs.Parse(args)
}

func (gc *getCommand) name() string {
	return gc.fs.Name()
}

func (gc *getCommand) run() error {

	if gc.keyname == "" || gc.encryptionKey == "" {
		return errors.New("keyname or encryptionKey should not be empty")
	}
	result, err := v.Get(gc.keyname, gc.encryptionKey)
	if err != nil {
		return err
	}
	fmt.Println("Secret Retrieved:")
	fmt.Println(result)
	return nil
}

func root(args []string) error {
	if len(args) < 1 {
		return errors.New("You must pass a sub command")
	}

	cmds := []runner{
		newGetCommand(),
		newSetCommand(),
	}

	subcommand := os.Args[1]

	for _, cmd := range cmds {
		if cmd.name() == subcommand {
			cmd.init(os.Args[2:])
			return cmd.run()
		}
	}

	return fmt.Errorf("Unknown subcommand: %s", subcommand)
}

func main() {
	if err := root(os.Args[1:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
