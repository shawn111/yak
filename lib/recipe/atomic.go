package recipe

import (
	"log"
	"os/exec"
	"strings"

	"github.com/goyak/yak/lib/host/ostree"
	"github.com/goyak/yak/lib/utils"
)

type AtomicRecipeConfig struct {
	BaseRecipeConfig
}

func (r AtomicRecipeConfig) IsInstallable() bool {
	return ostree.IsOstreeHost()
}

func (r AtomicRecipeConfig) Install(dryrun bool) bool {
	// backup current local config
	// ostree admin config-diff

	ostree.Backup(r.Repo)
	remoteName := strings.Split(r.Branch, "/")[0]

	addRemoteCmd := utils.Cmd("ostree", "remote", "add", "--if-not-exists", "--no-gpg-verify", remoteName, r.Source)
	utils.DoRun(addRemoteCmd, dryrun)

	pullCmd := exec.Command("ostree", "pull", remoteName, r.Commit)
	utils.DoRun(pullCmd, dryrun)

	deployCmd := exec.Command("rpm-ostree", "deploy", r.Commit)
	utils.DoRun(deployCmd, dryrun)

	// FIXME Prepare to reboot
	log.Printf("Deployed, please reboot. (systemctl reboot)\n")
	return true
}
