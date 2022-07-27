package gb

import "time"

type Timer struct {
	TimerCounter    int
	DividerRegister int
	ScanlineCounter int
}

type Core struct {
	Cartridge Cartridge
	Memory    Memory

	/*
	   +++++++++++++++++++++++
	   +        Joypad       +
	   +++++++++++++++++++++++
	*/
	Controller   driver.ControllerDriver
	JoypadStatus byte

	//Frames per-second
	FPS int
	//CPU clock
	Clock int
	//in CBG mode, clock might change to twice as original
	SpeedMultiple int

	Timer   Timer
	RamPath string
}

func (core *Core) Run() {
	ticker := time.NewTicker(time.Second / time.Duration(core.FPS))
	for range ticker.C {
		core.Update()
	}
}

func (core *Core) setupSaveLoop() {
	// each second check if there are new saves (to avoid thousands within a frame)
	saveTimer := time.Tick(time.Second)
	go func() {
		for range saveTimer {
			core.SaveRAM()
		}
	}()
}

func (core *Core) SaveRAM() {
	if core.Memory.dirty {
		core.Memory.dirty = false
		core.Cartridge.MBC.SaveRam(core.RamPath)
	}
}

func (core *Core) ReadMemory(address uint16) byte {
	if (address >= 0x4000) && (address <= 0x7FFF) {
		// are we reading from the rom memory bank?
		return core.Cartridge.MBC.ReadRomBank(address)
	} else if (address >= 0xA000) && (address <= 0xBFFF) {
		// are we reading from ram memory bank?
		return core.Cartridge.MBC.ReadRamBank(address)
	} else if 0xFF00 == address {

	} else if address == 0xFF01 {
		return core.SerialByte
	}
}

func (core *Core) Update() {
	cyclesThisUpdate := 0

	for cyclesThisUpdate < ((core.SpeedMultiple+1)*core.Clock)/core.FPS {
		cycles := 4

	}
}
