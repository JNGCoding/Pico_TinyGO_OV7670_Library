package CORE

type Program struct {
	Setup func()
	Loop  func()

	EXIT_FLAG bool
	EXIT_STR  string
	EXIT_CODE int8
}

func CreateApplication() *Program {
	var Result Program

	Result.EXIT_CODE = 0
	Result.EXIT_STR = "All Good"
	Result.EXIT_FLAG = false

	return &Result
}

func (p *Program) LetSetup(SETUP func()) {
	p.Setup = SETUP
}

func (p *Program) LetLoop(LOOP func()) {
	p.Loop = LOOP
}

func (p *Program) Exit(EXITCODE int8, EXITSTR string) {
	p.EXIT_CODE = EXITCODE
	p.EXIT_STR = EXITSTR
	p.EXIT_FLAG = true
}

func (p *Program) ExitInfo() (int8, string) {
	return p.EXIT_CODE, p.EXIT_STR
}

func (p *Program) Run() {
	p.Setup()
	for !p.EXIT_FLAG {
		p.Loop()
	}
}

func (p *Program) Reset() {
	p.EXIT_FLAG = true
}
