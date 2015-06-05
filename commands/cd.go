package commands

import "github.com/kamilhark/etcdsh/pathresolver"
import "github.com/kamilhark/etcdsh/etcdclient"
import "github.com/kamilhark/etcdsh/common"

type CdCommand struct {
	PathResolver *pathresolver.PathResolver
	etcdClient   *etcdclient.EtcdClient
}

func NewCdCommand(pathResolver *pathresolver.PathResolver, etcdClient *etcdclient.EtcdClient) *CdCommand {
	cdCommand := new(CdCommand)
	cdCommand.PathResolver = pathResolver
	cdCommand.etcdClient = etcdClient
	return cdCommand
}

func (cdCommand *CdCommand) Supports(command string) bool {
	return command == "cd"
}

func (cdCommand *CdCommand) Handle(args []string) {
	if len(args) == 1 {
		cdCommand.PathResolver.GoTo(args[0])
	} else {
		cdCommand.PathResolver.GoTo("")
	}
}

func (cdCommand *CdCommand) Verify(args []string) error {
	if len(args) > 1 {
		return common.NewStringError("'cd' command supports only one argument")
	}

	if len(args) == 0 {
		return nil
	}

	nextPath := cdCommand.PathResolver.Resolve(args[0])
	response, err := cdCommand.etcdClient.Get(nextPath)
	if err != nil {
		return err
	}

	if !response.Node.Dir {
		return common.NewStringError("not a directory")
	}

	return nil
}
