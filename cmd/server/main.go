package main

// @title Swagger ui-assignment API
// @version 1.0
// @description This is a sample ui-assignment server.
// @termsOfService http://swagger.io/terms/

// @contact.name Ben Chuang
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8000
// @BasePath /
func main() {

	ui := Ubiquiti{}
	ui.Initialize()

	ui.Run()
}
