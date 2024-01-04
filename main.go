package main

type NoFilenameError struct{}

func main() {
}

func (n NoFilenameError) Error() string {
	return "no filename was given"
}

func GetFilenameFromCommand(args []string) (string, error) {
	if len(args) < 3 {
		return "", NoFilenameError{}
	}
	return args[2], nil
}
