package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/sunny0826/dingtalk"
	"log"
	"os"
	"os/exec"
	"time"
)

// CheckCommand for check cmd
type CheckCommand struct {
	baseCommand
}

const localPath = "https://github.com/sunny0826/kustomize.git"
const upstreamPath = "https://github.com/kubernetes-sigs/kustomize.git"
const dirName = "kustomize"
const logFileName = "updateLog"

var dingtoken = os.Getenv("DING_TOKEN")

func (cc *CheckCommand) Init() {
	cc.command = &cobra.Command{
		Use:     "check",
		Short:   "Check docs change",
		Long:    "Check docs change",
		Aliases: []string{"v"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cc.runCheck(cmd, args)
		},
		Example: checkExample(),
	}
}

func (cc *CheckCommand) runCheck(command *cobra.Command, args []string) error {
	git := commandGit()
	cloneCmd := cc.commandClone(git)
	cloneOut, err := cloneCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(cloneOut))
	}

	remote := cc.commandRemoteAdd(git)
	remote.Dir = dirName
	remoteOut, err := remote.CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(remoteOut))
	}
	fetchCmd := cc.commandFetchUpstream(git)
	fetchCmd.Dir = dirName
	fetchOut, err := fetchCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(fetchOut))
	}
	diffCmd := cc.commandDiffTree(git)
	diffCmd.Dir = dirName
	diffOut, err := diffCmd.Output()
	loc, _:= time.LoadLocation("Asia/Shanghai")
	update_time := time.Now().In(loc).Format("2006-01-02 15:04:05")
	if string(diffOut) != "" {
		err := cc.sendDingtalk(string(diffOut))
		if err != nil {
			return err
		}
	}
	content := fmt.Sprintf("%s\n%s", update_time, string(diffOut))
	writeFile(logFileName, content)
	return nil
}

// commandGit git command bin path
func commandGit() string {
	gitProgram, err := exec.LookPath("git")
	if err != nil {
		log.Fatal("no 'git' program on path")
	}
	return gitProgram
}

// commandClone git clone
func (cc *CheckCommand) commandClone(gitCmd string) *exec.Cmd {
	cc.command.Println("start clone")
	return exec.Command(
		gitCmd,
		"clone",
		"--depth=1",
		localPath,
	)
}

// commandFetchUpstream git fetch upstream
func (cc *CheckCommand) commandFetchUpstream(gitCmd string) *exec.Cmd {
	cc.command.Println("fetch upstream")
	return exec.Command(
		gitCmd,
		"fetch",
		"upstream",
		"master",
	)
}

// commandDiffTree git diff-tree master upstream/master
func (cc *CheckCommand) commandDiffTree(gitCmd string) *exec.Cmd {
	cc.command.Println("diff tree")
	return exec.Command(
		"bash", "-c",
		fmt.Sprintf("%s diff --name-only master upstream/master | egrep '^examples|^docs'", gitCmd),
	)
}

// commandRemoteAdd git remote add upstream
func (cc *CheckCommand) commandRemoteAdd(gitCmd string) *exec.Cmd {
	cc.command.Println("remote add")
	return exec.Command(
		gitCmd,
		"remote",
		"add",
		"upstream",
		upstreamPath,
	)
}

func (cc *CheckCommand) commandDiffDosc(gitCmd string) *exec.Cmd {
	cc.command.Println("diff dosc")
	return exec.Command(
		gitCmd,
		"remote",
		"add",
		"upstream",
		upstreamPath,
	)
}

func (cc *CheckCommand) sendDingtalk(msg string) error {
	webHook := dingtalk.NewWebHook(dingtoken, "")

	//test send text message
	err := webHook.SendMarkdownMsg("kustomize change", fmt.Sprintf("# Kustomize Change \n\n [%s](%s) \n\n >%s \n\n", localPath, localPath, msg), false, "15235111699")
	if nil != err {
		return err
	}
	cc.command.Println("send dingtalk successful!")
	return nil
}

func writeFile(fileName string, content string) {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(fileName, err)
		return
	}
	defer f.Close()
	f.WriteString(content)
}

func checkExample() string {
	return `kust-check check`
}
