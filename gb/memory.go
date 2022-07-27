package gb

import (
	"log"
)

type Memory struct {
	MainMemory [0x10000]byte
	dirty      bool
}

func (core *Core) initMemory() {
	log.Println("[Core] Start to initialize memory...")

	log.Println("[Memory] Load first 32KByte of rom data into memory")
	//Load first 32KB of ROM into 0000-7FFF
	for i := 0x0000; i < core.Cartridge.Props.ROMLength && i < 0x8000; i++ {
		core.Memory.MainMemory[i] = core.Cartridge.MBC.ReadRom(uint16(i))
	}

	core.Memory.MainMemory[0xFF05] = 0x00
	core.Memory.MainMemory[0xFF06] = 0x00
	core.Memory.MainMemory[0xFF07] = 0x00
	core.Memory.MainMemory[0xFF0F] = 0xE1
	core.Memory.MainMemory[0xFF10] = 0x80
	core.Memory.MainMemory[0xFF11] = 0xBF
	core.Memory.MainMemory[0xFF12] = 0xF3
	core.Memory.MainMemory[0xFF14] = 0xBF
	core.Memory.MainMemory[0xFF16] = 0x3F
	core.Memory.MainMemory[0xFF17] = 0x00
	core.Memory.MainMemory[0xFF19] = 0xBF
	core.Memory.MainMemory[0xFF1A] = 0x7F
	core.Memory.MainMemory[0xFF1B] = 0xFF
	core.Memory.MainMemory[0xFF1C] = 0x9F
	core.Memory.MainMemory[0xFF1E] = 0xBF
	core.Memory.MainMemory[0xFF20] = 0xFF
	core.Memory.MainMemory[0xFF21] = 0x00
	core.Memory.MainMemory[0xFF22] = 0x00
	core.Memory.MainMemory[0xFF23] = 0xBF
	core.Memory.MainMemory[0xFF24] = 0x77
	core.Memory.MainMemory[0xFF25] = 0xF3
	core.Memory.MainMemory[0xFF26] = 0xF1
	core.Memory.MainMemory[0xFF40] = 0x91
	core.Memory.MainMemory[0xFF42] = 0x00
	core.Memory.MainMemory[0xFF43] = 0x00
	core.Memory.MainMemory[0xFF45] = 0x00
	core.Memory.MainMemory[0xFF47] = 0xFC
	core.Memory.MainMemory[0xFF48] = 0xFF
	core.Memory.MainMemory[0xFF49] = 0xFF
	core.Memory.MainMemory[0xFF4A] = 0x00
	core.Memory.MainMemory[0xFF4B] = 0x00
	core.Memory.MainMemory[0xFFFF] = 0x00

	core.setupSaveLoop()
}
