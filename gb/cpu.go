package gb

import "log"

/*
	Registers
	  16bit Hi   Lo   Name/Function
	  AF    A    -    Accumulator & Flags
	  BC    B    C    BC
	  DE    D    E    DE
	  HL    H    L    HL
	  SP    -    -    Stack Pointer
	  PC    -    -    Program Counter/Pointer
*/
type Registers struct {
	A  byte
	B  byte
	C  byte
	D  byte
	E  byte
	F  byte
	HL uint16
	PC uint16
	SP uint16
}

type Flags struct {
	Zero      bool
	Sub       bool
	HalfCarry bool
	Carry     bool

	//	IME - Interrupt Master Enable Flag (Write Only)
	//  	0 - Disable all Interrupts
	//  	1 - Enable all Interrupts that are enabled in IE Register (FFFF)
	InterruptMaster bool

	PendingInterruptEnabled bool
}

type CPU struct {
	Registers Registers
	Flags     Flags
	Halt      bool
}

func (core *Core) initCPU() {
	log.Println("[Core] Initialize CPU flags and registers")

	core.CPU.Flags.Zero = true
	core.CPU.Flags.Sub = false
	core.CPU.Flags.HalfCarry = true
	core.CPU.Flags.Carry = true
	core.CPU.Flags.InterruptMaster = false

	/*
		Initialize register after BIOS
		AF=$01B0
		BC=$0013
		DE=$00D8
		HL=$014D
		Stack Pointer=$FFFE
	*/
	core.CPU.Registers.A = 0x01
	core.CPU.Registers.B = 0x00
	core.CPU.Registers.C = 0x13
	core.CPU.Registers.D = 0x00
	core.CPU.Registers.E = 0xD8
	core.CPU.Registers.F = 0xB0
	core.CPU.Registers.HL = 0x014D
	core.CPU.Registers.PC = 0x0100
	core.CPU.Registers.SP = 0xFFFE
}

func (core *Core) ExecuteNextOPCode() int {
	opcode := core.ReadMemory(core.CPU.Registers.PC)
	core.CPU.Registers.PC++
	return core.ExecuteOPCode(opcode)
}

func (core *Core) ExecuteOPCode(opcode byte) int {

}
