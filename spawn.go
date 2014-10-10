package chromeapp

import (
	"fmt"
	"os"
	"os/exec"
)

func spawn(executable string, argv []string, options ...map[string]interface{}) error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Could not get work directory")
	}
	var env []string
	if len(options) == 1 {
		opts := options[0]
		if opts["Dir"] != nil {
			wd = opts["Dir"].(string)
		}
		if opts["Env"] != nil {
			env = opts["Env"].([]string)
		}
	}

	cmd := exec.Command(executable, argv...)
	cmd.Dir = wd
	if len(env) > 0 {
		cmd.Env = env
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	go func() {
		err = cmd.Start()
		if err != nil {
			panic(fmt.Errorf("Could not start process %s\n", executable))
		}
		c := make(chan error, 1)
		c <- cmd.Wait()
	}()
	return nil
}
