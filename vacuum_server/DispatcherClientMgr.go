package vacuum_server

var ()

func manageDispatcherClient() {
	dispatcherClient = newDispatcherClient()
	dispatcherClient.connect()

}
