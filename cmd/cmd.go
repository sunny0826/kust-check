package cmd

func CmdInit() *baseCommand {
	cli := NewCli()
	baseCmd := &baseCommand{
		command: cli.rootCmd,
	}
	// add version command
	baseCmd.AddCommand(&VersionCommand{})
	// add check command
	baseCmd.AddCommand(&CheckCommand{})

	return baseCmd
}
