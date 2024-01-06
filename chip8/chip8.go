package chip8

type chip8 struct {
	memory    []byte
	registers []byte
	timers    []byte
	stack     []uint16
	pc        uint16
	sp        uint8
	i         uint16
}

func NewChip8() *chip8 {
	chip := &chip8{
		memory:    make([]byte, 4096),
		registers: make([]byte, 16),
		timers:    make([]byte, 2),
		stack:     make([]uint16, 16),
		pc:        0x200,
		sp:        0,
		i:         0,
	}

	return chip
}
