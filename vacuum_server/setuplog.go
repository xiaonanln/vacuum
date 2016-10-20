package vacuum_server

import log "github.com/Sirupsen/logrus"

func setupLog(loglevel string) {
	lvl, err := log.ParseLevel(loglevel)
	if err != nil {
		panic(err)
	}

	log.SetLevel(lvl)
	log.SetFormatter(&log.TextFormatter{})
}
