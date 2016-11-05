package vacuum_server

import "github.com/xiaonanln/vacuum/vlog"

func setupLog(loglevel string) {
	lvl, err := vlog.ParseLevel(loglevel)
	if err != nil {
		panic(err)
	}

	vlog.SetLevel(lvl)
}
