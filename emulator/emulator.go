package emulator

type Screen []byte

func ClearScreen(screen *Screen) *Screen {
	for i := range *screen {
		(*screen)[i] = 0
	}
	return screen
}
