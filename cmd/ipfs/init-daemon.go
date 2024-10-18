package main

import (
	oldcmds "github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core/commands"
	fsrepo "github.com/ipfs/go-ipfs/repo/fsrepo"

	cmds "github.com/ipfs/go-ipfs-cmds"
)

var initDaemonCmd = &cmds.Command{
	Helptext: cmds.HelpText{
		Tagline: "Initializes ipfs config file and run daemon.",
		ShortDescription: `
Initializes ipfs configuration files and generates a new keypair.

If you are going to run IPFS in server environment, you may want to
initialize it using 'server' profile.

For the list of available profiles see 'ipfs config profile --help'

ipfs uses a repository in the local file system. By default, the repo is
located at ~/.ipfs. To change the repo location, set the $IPFS_PATH
environment variable:

    export IPFS_PATH=/path/to/ipfsrepo
`,
	},
	Arguments: []cmds.Argument{
		cmds.FileArg("default-config", false, false, "Initialize with the given configuration.").EnableStdin(),
	},
	Options: []cmds.Option{
		cmds.StringOption(algorithmOptionName, "a", "Cryptographic algorithm to use for key generation.").WithDefault(algorithmDefault),
		cmds.IntOption(bitsOptionName, "b", "Number of bits to use in the generated RSA private key."),
		cmds.BoolOption(emptyRepoOptionName, "e", "Don't add and pin help files to the local storage."),
		cmds.StringOption(profileOptionName, "p", "Apply profile settings to config. Multiple profiles can be separated by ','"),
		cmds.StringOption(announceAddressName, "n", "Onion service address to announce."),
		cmds.StringOption(bootStrapAddressName, "t", "Bootstrap peer address."),
		cmds.StringOption(ppChannelUrlName, "u", "PPChannel URL.").WithDefault("http://localhost:28080"),
		cmds.IntOption(commandPortName, "m", "PPChannel command callback listen port").WithDefault(30500),
		cmds.StringOption(torPathName, "r", "Tor executable path."),
		cmds.StringOption(torDataDirName, "d", "Tor data dir."),
		cmds.StringOption(torConfigPathName, "o", "Tor configuration path."),
		cmds.StringOption(routingTypeName, "doute", "Routing type."),
		cmds.StringOption(routingOptionKwd, "Overrides the routing option").WithDefault(routingOptionDefaultKwd),
		cmds.BoolOption(mountKwd, "Mounts IPFS to the filesystem"),
		cmds.BoolOption(writableKwd, "Enable writing objects (with POST, PUT and DELETE)"),
		cmds.StringOption(ipfsMountKwd, "Path to the mountpoint for IPFS (if using --mount). Defaults to config setting."),
		cmds.StringOption(ipnsMountKwd, "Path to the mountpoint for IPNS (if using --mount). Defaults to config setting."),
		cmds.BoolOption(unrestrictedApiAccessKwd, "Allow API access to unlisted hashes"),
		cmds.BoolOption(unencryptTransportKwd, "Disable transport encryption (for debugging protocols)"),
		cmds.BoolOption(enableGCKwd, "Enable automatic periodic repo garbage collection"),
		cmds.BoolOption(adjustFDLimitKwd, "Check and raise file descriptor limits if needed").WithDefault(true),
		cmds.BoolOption(migrateKwd, "If true, assume yes at the migrate prompt. If false, assume no."),
		cmds.BoolOption(enablePubSubKwd, "Instantiate the ipfs daemon with the experimental pubsub feature enabled."),
		cmds.BoolOption(enableIPNSPubSubKwd, "Enable IPNS record distribution through pubsub; enables pubsub."),
		cmds.BoolOption(enableMultiplexKwd, "DEPRECATED"),
		cmds.StringOption(homeRootFolder, "hrf", "Current folder(relative by user home)."),

		// TODO need to decide whether to expose the override as a file or a
		// directory. That is: should we allow the user to also specify the
		// name of the file?
		// TODO cmds.StringOption("event-logs", "l", "Location for machine-readable event logs."),
	},
	NoRemote: true,
	Extra:    commands.CreateCmdExtras(commands.SetDoesNotUseRepo(true), commands.SetDoesNotUseConfigAsInput(true)),
	PreRun:   commands.DaemonNotRunning,
	Run: func(req *cmds.Request, res cmds.ResponseEmitter, env cmds.Environment) error {
		cctx := env.(*oldcmds.Context)
		if !fsrepo.IsInitialized(cctx.ConfigRoot) {
			initFunc(req, res, env)
		}
		return daemonFunc(req, res, env)
	},
}
