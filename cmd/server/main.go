package main

func main() {

	ui := Ubiquiti{}
	ui.Initialize("user", "password", "")

	ui.Run()
}
