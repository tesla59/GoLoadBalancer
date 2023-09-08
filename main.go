package main

func init() {
	BuildImage("./worker", "worker-image")
}

func main() {
	RunImage("worker-image")
}
