package gb

import (
	"bufio"
	"log"
	"os"
)

/*
0x00  ROM ONLY                 0x19  MBC5
0x01  MBC1                     0x1A  MBC5+RAM
0x02  MBC1+RAM                 0x1B  MBC5+RAM+BATTERY
0x03  MBC1+RAM+BATTERY         0x1C  MBC5+RUMBLE
0x05  MBC2                     0x1D  MBC5+RUMBLE+RAM
0x06  MBC2+BATTERY             0x1E  MBC5+RUMBLE+RAM+BATTERY
0x08  ROM+RAM                  0x20  MBC6
0x09  ROM+RAM+BATTERY          0x22  MBC7+SENSOR+RUMBLE+RAM+BATTERY
0x0B  MMM01
0x0C  MMM01+RAM
0x0D  MMM01+RAM+BATTERY
0x0F  MBC3+TIMER+BATTERY
0x10  MBC3+TIMER+RAM+BATTERY   0xFC  POCKET CAMERA
0x11  MBC3                     0xFD  BANDAI TAMA5
0x12  MBC3+RAM                 0xFE  HuC3
0x13  MBC3+RAM+BATTERY         0xFF  HuC1+RAM+BATTERY
*/
var cartridgeTypeMap = map[byte]string{
	byte(0x00): "ROM ONLY",
	byte(0x01): "MBC1",
	byte(0x02): "MBC1+RAM",
	byte(0x03): "MBC1+RAM+BATTERY",
	byte(0x05): "MBC2",
	byte(0x06): "MBC2+BATTERY",
	byte(0x08): "ROM+RAM",
	byte(0x09): "ROM+RAM+BATTERY",
	byte(0x0B): "MMM01",
	byte(0x0C): "MMM01+RAM",
	byte(0x0D): "MMM01+RAM+BATTERY",
	byte(0x0F): "MBC3+TIMER+BATTERY",
	byte(0x10): "MBC3+TIMER+RAM+BATTERY",
	byte(0x11): "MBC3",
	byte(0x12): "MBC3+RAM",
	byte(0x13): "MBC3+RAM+BATTERY",
	byte(0x15): "MBC4",
	byte(0x16): "MBC4+RAM",
	byte(0x17): "MBC4+RAM+BATTERY",
	byte(0x19): "MBC5",
	byte(0x1A): "MBC5+RAM",
	byte(0x1B): "MBC5+RAM+BATTERY",
	byte(0x1C): "MBC5+RUMBLE",
	byte(0x1D): "MBC5+RUMBLE+RAM",
	byte(0x1E): "MBC5+RUMBLE+RAM+BATTERY",
	byte(0xFC): "POCKET CAMERA",
	byte(0xFD): "BANDAI TAMA5",
	byte(0xFE): "HuC3",
	byte(0xFF): "HuC1+RAM+BATTERY",
}

/*
	ROM bank number is linked to the ROM Size byte (0148).
		1 bank = 16 KBytes
	0x00 means no bank required.
*/
var RomBankMap = map[byte]uint8{
	byte(0x00): 2,
	byte(0x01): 4,
	byte(0x02): 8,
	byte(0x03): 16,
	byte(0x04): 32,
	byte(0x05): 64,
	byte(0x06): 128,
	byte(0x52): 72,
	byte(0x53): 80,
	byte(0x54): 96,
}

/*
	RAM bank number is linked to the RAM Size byte (0149).
		1 bank = 8 KBytes
	0x00 means no bank required.
*/
var RamBankMap = map[byte]uint8{
	byte(0x00): 0,
	byte(0x01): 1,
	byte(0x02): 1,
	byte(0x03): 4,
}

type CartridgeProps struct {
	MBCType string
	ROMBank uint8
	RAMBank uint8
}

type Cartridge struct {
	Props CartridgeProps
	MBC   MBC
}

type MBC interface {
	ReadRom(uint16) byte
	ReadRomBank(uint16) byte
	ReadRamBank(uint16) byte
	WriteRamBank(uint16, byte)
	SaveRamBank(string)
}

type MBCRom struct {
	// ROM data
	rom []byte
	// Current ROM-Bank number
	CurrentROMBank byte
	RAMBank        [0x8000]byte
	// Current RAM-Bank number
	CurrentRAMBank byte
	EnableRAM      bool
}

/**
Read a byte from RAM bank.
In ROM only cartridge, RAM is not supported.
*/
func (mbc *MBCRom) ReadRamBank(address uint16) byte {
	return byte(0x00)
}

// WriteRamBank /**
func (mbc *MBCRom) WriteRamBank(address uint16, data byte) {

}

/**
Read a byte from ROM bank.
In ROM only cartridge, ROM banking is not supported.
*/
func (mbc *MBCRom) ReadRomBank(address uint16) byte {
	return mbc.rom[address]
}

/**
Read a byte from raw rom via address
*/
func (mbc *MBCRom) ReadRom(address uint16) byte {
	return mbc.rom[address]
}

func (mbc *MBCRom) HandleBanking(address uint16, val byte) {
}

func (mbc *MBCRom) SaveRam(path string) {
}

/*
	Read cartridge data from file
*/
func (core *Core) readRomFile(romPath string) []byte {
	return readDataFile(romPath, false)
}

type MBC1 struct {
	rom            []byte
	RAMBank        []byte
	CurrentRAMBank byte
	EnableRAM      bool
	ROMBankingMode bool
}

func (mbc *MBC1) ReadRoomBank(address uint16) byte {
	newAddress := uint32(address - 0x4000)
	return mbc.rom[newAddress+(uint32(mbc.CurrentRAMBank)*0x4000)]
}

func (mbc *MBC1) WriteRoomBank(address uint16, date byte) {
	if !mbc.EnableRAM {
		return
	}
	newAddress := uint32(address - 0xA000)
	mbc.rom[newAddress+(uint32(mbc.CurrentRAMBank)*0x2000)] = date
}

func (mbc *MBC1) ReadRom(address uint16) byte {
	return mbc.rom[address]
}

func (mbc *MBC1) HandleBanking(address uint16, val byte) {
	if address < 0x2000 {
		/**
		0000-1FFF - RAM Enable (Write Only)
		  00h  Disable RAM (default)
		  0Ah  Enable RAM
		*/
		mbc.DoRamBankEnable(val)
	}

	if address >= 0x4000 && address < 0x6000 {
		if mbc.ROMBankingMode {
			mbc.DoChangeHiRomBank(val)
		} else {
			mbc.DoRAMBankChange(val)
		}
	}
}

//
func (mbc *MBC1) DoRamBankEnable(val byte) {
	testData := val & 0xF
	if testData == 0xA {
		mbc.EnableRAM = true
	}

	if testData == 0x0 {
		mbc.EnableRAM = false
	}
}

func (mbc *MBC1) DoChangeHiRomBank(val byte) {
	// turn off the upper 3 bits of the current rom
	mbc.CurrentRAMBank &= 31

	// turn off the lower 5 bits of the data
	val &= 224
	mbc.CurrentRAMBank |= val
	if mbc.CurrentRAMBank == 0 {
		mbc.CurrentRAMBank++
	}
}

func (mbc *MBC1) DoRAMBankChange(val byte) {
	mbc.CurrentRAMBank = val & 0x3
}

func (mbc *MBC1) SaveRam(path string) {
	writeRamFile(path, mbc.RAMBank)
}

func readDataFile(path string, ram bool) []byte {
	name := "rom"
	if ram {
		name = "ram"
	}
	log.Println("[Core] Loading", name, "file...")
	romFile, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) && ram {
			return nil
		}

		log.Fatal(err)
	}
	defer romFile.Close()

	state, statsErr := romFile.Stat()
	if statsErr != nil {
		log.Fatal(statsErr)
	}
	var size int64 = state.Size()
	bytes := make([]byte, size)

	bufReader := bufio.NewReader(romFile)
	_, err = bufReader.Read(bytes)
	log.Println("[Core]", size, "Bytes", name, "loaded")
	return bytes
}

func writeRamFile(ramPath string, data []byte) {
	ramFile, err := os.Create(ramPath)
	if err != nil {
		log.Fatal(err)
	}
	defer ramFile.Close()

	bufWriter := bufio.NewWriter(ramFile)
	size, err := bufWriter.Write(data)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("[Core] %d Bytes ram written\n", size)
}
