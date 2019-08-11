package nes

import (
	"errors"
	"image"
)

type Nes struct {
	cassette Ines
	cpu      *Cpu
	ppu      *Ppu
	bus      *Bus
}

func NewNes(cassette Ines) *Nes {
	wram := NewRam(0x800)
	renderer := NewRenderer()
	controller := NewController()
	bus := NewBus(wram, cassette.PrgRom())
	cpu := NewCpu(bus)
	ppu := NewPpu(bus, cassette.ChrRom(), renderer)
	bus.cpu = cpu
	bus.ppu = ppu
	bus.controller = controller
	return &Nes{
		cassette: cassette,
		cpu:      cpu,
		ppu:      ppu,
		bus:      bus,
	}
}

func (n *Nes) isSetCassette() bool {
	return n.cassette != nil
}

func (n *Nes) Init() error {
	if !n.isSetCassette() {
		return errors.New("cassette must be set")
	}

	n.cpu.PC = 0x8000
	return nil
}

// run 60fps
func (n *Nes) Run(){

	for !n.step(){
	}
}

func (n *Nes) step() bool {

	// check interrupt
	if n.cpu.intrrupt != nil && !n.cpu.isIrqForbitten(){
		// fmt.Println("========= Interrupt! =========")
		n.cpu.intrrupt()
	}
	n.cpu.intrrupt = nil

	pc := n.cpu.PC

	// decode
	b := n.bus.Load(pc)
	inst := n.cpu.decode(b)
	cycle := inst.cycle
	addr := n.cpu.solveAddrMode(inst.addrMode)

	// for debug
	n.cpu.dump(b, addr, inst.mnemonic, inst.addrMode)

	n.cpu.advance(inst.addrMode)
	n.cpu.execute(inst, addr)

	n.cpu.cycle += cycle

	if n.ppu.run(cycle * 3){
		n.ppu.renderer.render()
		return true
	}

	return false
}

func (n *Nes) Buffer() *image.RGBA {
	return n.ppu.renderer.Buffer()
}

func (n *Nes) PushButton(b [8]bool) {
	n.bus.controller.SetButton(b)
}