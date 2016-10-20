package vacuum_server

import log "github.com/Sirupsen/logrus"

func setupLog() {
	log.SetFormatter(&log.TextFormatter{})
}
