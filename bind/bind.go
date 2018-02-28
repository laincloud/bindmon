package bind

import (
	"github.com/facebookgo/pidfile"
	"log"
	"syscall"
)

func Reload(path string) error {
	log.Println("reload named " + path)
	pidfile.SetPidfilePath(path)
	pid, err := pidfile.Read()
	if err != nil {
		log.Println(err)
		return err
	}
	err = syscall.Kill(pid, syscall.SIGHUP)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}
