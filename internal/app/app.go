package app

import (
	"messenger/internal/dbclient"
	"messenger/internal/httpserver"
)

func Run() {
	dbclient.New()
	httpserver.Start()
}
