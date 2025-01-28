package commands

import (
	"fmt"
	"io"
	"runtime"

	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"
	plusVersion "github.com/ipfs/go-ipfs/version"

	cmds "github.com/ipfs/go-ipfs-cmds"
)

type VersionPlusOutput struct {
	Version     string
	BuildDate   string
	Repo        string
	System      string
	Golang      string
	BuildNumber string
	CommitHash  string
}

const (
	versionPlusNumberOptionName = "number"
	versionPlusDateOptionName   = "date"
	versionPlusRepoOptionName   = "repo"
	versionPlusAllOptionName    = "all"
)

var VersionPlusCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline:          "Show ipfs version information.",
		ShortDescription: "Returns the current version of ipfs and exits.",
	},
	Subcommands: map[string]*cmds.Command{
		"deps": depsVersionCommand,
	},

	Options: []cmds.Option{
		cmds.BoolOption(versionNumberOptionName, "n", "Only show the version number."),
		cmds.BoolOption(versionCommitOptionName, "Show the commit hash."),
		cmds.BoolOption(versionRepoOptionName, "Show repo version."),
		cmds.BoolOption(versionAllOptionName, "Show all version information"),
	},
	// must be permitted to run before init
	Extra: CreateCmdExtras(SetDoesNotUseRepo(true), SetDoesNotUseConfigAsInput(true)),
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		return cmds.EmitOnce(res, &VersionPlusOutput{
			Version:     plusVersion.Version(),
			BuildDate:   plusVersion.BuildDate(),
			Repo:        fmt.Sprint(fsrepo.RepoVersion),
			System:      runtime.GOARCH + "/" + runtime.GOOS, //TODO: Precise version here
			Golang:      runtime.Version(),
			BuildNumber: plusVersion.BuildNumber(),
			CommitHash:  plusVersion.CommitHash(),
		})
	},
	Encoders: cmds.EncoderMap{
		cmds.Text: cmds.MakeTypedEncoder(func(req *cmds.Request, w io.Writer, version *VersionPlusOutput) error {
			all, _ := req.Options[versionAllOptionName].(bool)
			if all {
				ver := version.Version

				out := fmt.Sprintf("go-ipfs version: %s\n"+
					"repo version: %s\nsystem version: %s\ngolang version: %s\nbuild date:%s\nbuild number:%s\nbuild commit: %s\n",
					ver, version.Repo, version.System, version.Golang, version.BuildDate, version.BuildNumber, version.CommitHash)
				fmt.Fprint(w, out)
				return nil
			}

			commit, _ := req.Options[versionPlusDateOptionName].(bool)
			commitTxt := ""
			if commit && version.BuildDate != "" {
				commitTxt = " " + version.BuildDate
			}

			repo, _ := req.Options[versionRepoOptionName].(bool)
			if repo {
				fmt.Fprintln(w, version.Repo)
				return nil
			}

			number, _ := req.Options[versionNumberOptionName].(bool)
			if number {
				fmt.Fprintln(w, version.Version+commitTxt)
				return nil
			}

			fmt.Fprint(w, fmt.Sprintf("ipfs version %s%s\n", version.Version, commitTxt))
			return nil
		}),
	},
	Type: VersionPlusOutput{},
}
