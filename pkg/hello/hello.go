package hello

func SayHello(name string) string {
	str := "Hello world!"
	if name != "" {
		str = "Hello, " + name + "."
	}
	return str
}
