package web

func (app *App) handleCall() {
	for call := range app.queue {
		app.infoLog.Print(call.Id)
	}
}
