package commands

import "strings"
import "github.com/kamilhark/etcd-console/pathresolver"
import "github.com/kamilhark/etcd-console/etcdclient"
import "github.com/kamilhark/etcd-console/common"

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
	return strings.EqualFold(command, "cd")
}

func (cdCommand *CdCommand) Handle(args []string) {

	switch {
	case (len(args) == 0):
		cdCommand.PathResolver.Clear()
	case (args[0] == ".."):
		cdCommand.PathResolver.RemoveLast()
	case (args[0] == "."):
		return
	default:
		{
			pathElements := strings.Split(args[0], "/")
			for _, element := range pathElements {
				cdCommand.PathResolver.Add(element)
			}
		}
	}
}

func (cdCommand *CdCommand) Verify(args []string) error {
	if len(args) > 1 {
		return common.NewStringError("'cd' command supports only one argument")
	}

	if len(args) == 0 {
		return nil
	}

	if args[0] == ".." {
		return nil
	}

	nextPath := cdCommand.PathResolver.CurrentPath() + args[0]
	response, err := cdCommand.etcdClient.Get(nextPath)
	if err != nil {
		return err
	}

	if !response.Node.Dir {
		return common.NewStringError("not a directory")
	}

	return nil
}
