package easydeploy

import "testing"

func TestBuildSimpleCommand(t *testing.T) {
	cmd := BuildSimpleCommand(func(deploy Deploy, config *ServerConfig, args ...interface{}) error {
		t.Log("args0:", args[0])
		t.Log("args0:", args[1])
		return nil
	}, 1, 2)
	cmd.Run(nil, nil)

	cmd = BuildSimpleCommand(func(deploy Deploy, config *ServerConfig, args ...interface{}) error {
		t.Log("args:", args)
		return nil
	})
	cmd.Run(nil, nil)
}
