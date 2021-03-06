// lib for test.go and main.py -> out.go

// if you want to execute .gr files:
//	   $ main.py (FILE.gr) -o (OUT.go)
//	   $ gc OUT.go (gc so it includes this file and gutil)
//	   $ ./OUT

include "gutil"

type Flags struct {
	CMP bool
	//RMDMflag
}

type Regs struct {
	//RMDMregs
	LDA int // command
	LDB int // math i
	LRA int // math r
	LDX int // jump
	LDH int // jump back
	LDS int // stream
	LRI string // string
	LDC interface{} // converter
}

type Cpu struct {
	STACK []interface{}
	MEM   []interface{}
	CODE  []Op
	ADDR  int
	ACC   int
	TICK  uint
	FLAGS Flags
	REGS  Regs
	STACK_LEN int
	// stack top
	LABELS map[int]int
	// id:line
}

func MakeCpu() *Cpu {
	InitConv()
	return &Cpu{
		[]interface{}{},
		[]interface{}{},
		[]Op{},
		0,
		0,
		uint(0),
		//RMDMflag
		Flags{false},
		//RMDMregs
		Regs{0, 0, 0, 0, 0, 0, "", nil},
		0,
		map[int]int{},
	}
}

var CPU *Cpu = MakeCpu()

func (c *Cpu) Print() {
	print("\n")
	if (c == CPU) {
		printf("Main Cpu:\n")
	} else {
		printf("Extra Cpu:\n")
	}

	printf("slen: %d\n", c.STACK_LEN)
	if (c.STACK_LEN > 51) {
		printf("stack: %v, ...\n", c.STACK[:50])
	} else {
		printf("stack: %v\n", c.STACK)
	}

	printf("mlen: %d\n", len(c.MEM))
	if (len(c.MEM) > 51) {
		printf("mem: %v, ...\n", c.MEM[:50])
	} else {
		printf("mem: %v\n", c.MEM)
	}

	printf("@: %d\n", c.ADDR)
	printf("t: %d\n", c.TICK)
	printf("ac: %d\n", c.ACC)
	if ( len(c.LABELS)>0 ) {
		print("labels{\n")
		for id, line := range c.LABELS {
			printf("  %d:%d\n", id, line)
		}
		print("}\n")
	}
	if ( streamcount>0 ) {
		print("streams{\n")
		for id, name:= range sstreams{
			printf("  %d:%s\n", id, name)
		}
	}
	print("}\n")
	print("flags{\n")
	printf("  CMP:%t\n", c.FLAGS.CMP)
	//RMDMflag update flags
	print("}\n")
	print("regs{\n")
	printf("  LDA:%d\n", c.REGS.LDA)
	printf("  LDB:%d\n", c.REGS.LDB)
	printf("  LRA:%d\n", c.REGS.LRA)
	printf("  LDX:%d\n", c.REGS.LDX)
	printf("  LDH:%d\n", c.REGS.LDH)
	printf("  LDS:%d\n", c.REGS.LDS)
	printf("  LDC:`%v`\n", c.REGS.LDC)
	printf("  LRI:\"%s\"\n", c.REGS.LRI)
	print("}\n\n##########\n")
}

func (c *Cpu) NextTick() {
	c.RunCode()
	c.TICK++
}

func (c *Cpu) RunCode() {
	c.RunAsmCode()
	c.ADDR++
}

type Op struct {
	Op int
	Arg interface{} // default 0
	ArgType bool // immd/addr
}

func (c *Cpu) MakeOp(args... interface{}) Op {
	if (len(args) == 3) {
		return Op{(args[0]).(int), args[1], (args[2]).(bool)}
	} else if (len(args) == 2) {
		return Op{(args[0]).(int), args[1], false}
	} else if (len(args) == 1) {
		return Op{(args[0]).(int), 0, false}
	} else if (len(args) == 4) {
		c.REGS.LDA = (args[3]).(int)
		return Op{(args[0]).(int), args[1], (args[2]).(bool)}
	} else if (len(args) == 5) {
		c.REGS.LDA = (args[3]).(int)
		c.REGS.LDB = (args[4]).(int)
		return Op{(args[0]).(int), args[1], (args[2]).(bool)}
	} else {
		return Op{O_EXIT, 0, false}
	}
}

include "ops"
include "errs"

var sstreams map[int]string = map[int]string{
	0:"stdout",
	1:"stderr",
}

var streams map[int]*bufio.Writer = map[int]*bufio.Writer{
	0:stdout,
	1:stderr,
}

var streamcount = 1

var fd *FILE
var err error
var arg interface{}
var OP Op

const (
	C_INT_FLOAT = iota
	C_FLOAT_INT = iota
	C_ANY_STRING= iota
)

var converter = make(map[int]func(interface{}) interface{})


func InitConv(){
	converter[C_FLOAT_INT] = func (c interface{}) (interface{}) {
		return interface{}(int((c).(float64)))
	}

	converter[C_ANY_STRING] = func (c interface{}) (interface{}) {
		return interface{}(fs("%v", c))
	}

	converter[C_INT_FLOAT] = func (c interface{}) (interface{}) {
		return interface{}(float64((c).(int)))
	}
}

func (c *Cpu) AddLabel( line int ) {
	c.LABELS[len(c.LABELS)] = line
}

//TODO:
//	use more than one immd
//	use E_ stuff
func (c *Cpu) RunAsmCode() {
	OP = c.CODE[c.ADDR] // Op struct
	arg = OP.Arg
	if (OP.ArgType) {
		arg = c.MEM[(OP.Arg).(int)]
	}

	switch (OP.Op) {
		// syscall
		case O_WRITENL:
			fprintf(streams[c.REGS.LDS], "\n")
		case O_WRITE_LRI:
			fprintf(streams[c.REGS.LDS], c.REGS.LRI)
		case O_WRITE:
			fprintf(streams[c.REGS.LDS], "%v", arg)
		case O_WRITE_LRA:
			fprintf(streams[c.REGS.LDS], "%d", c.REGS.LRA)
		case O_FLUSH:
			streams[c.REGS.LDS].Flush()
		case O_OPEN:
			streamcount++
			fd, err = fmake((arg).(string))
			if (err != nil) {
				c.Print()
				panic(err)
			}
			streams[streamcount] = fwriter(fd)
			sstreams[streamcount] = (arg).(string)
			if (OP.ArgType) {
				c.REGS.LDS = streamcount
			}
		case O_EXIT:
			exit((arg).(int))

		// stack
		case O_PUSH:
			c.STACK_LEN++
			c.STACK = append([]interface{}{arg}, c.STACK...)
		case O_PUSH_LRA:
			c.STACK_LEN++
			c.STACK = append([]interface{}{interface{}((c.REGS.LRA))}, c.STACK...)
		case O_POP_LDA:
			c.STACK_LEN--
			c.REGS.LDA = (c.STACK[0]).(int)
			c.STACK = c.STACK[1:]
		case O_POP_LDB:
			c.STACK_LEN--
			c.REGS.LDB = (c.STACK[0]).(int)
			c.STACK = c.STACK[1:]
		case O_POP:
			c.STACK_LEN--
			c.MEM = append(c.MEM, (c.STACK[0]).(int))
			c.STACK = c.STACK[1:]

		// regs
		case O_LRI:
			c.REGS.LRI = (arg).(string)
		case O_LDA:
			c.REGS.LDA = (arg).(int)
		case O_LDB:
			c.REGS.LDB = (arg).(int)
		case O_LDS:
			c.REGS.LDS = (arg).(int)
		case O_LDX:
			c.REGS.LDX = (arg).(int)-1
		case O_LDH:
			c.REGS.LDH = (arg).(int)-1
		case O_ACC2LDA:
			c.REGS.LDA = c.ACC

		// math
		case O_RADD:
			c.REGS.LRA = c.REGS.LDB+c.REGS.LDA
		case O_RSUB:
			c.REGS.LRA = c.REGS.LDB-c.REGS.LDA
		case O_RMUL:
			c.REGS.LRA = c.REGS.LDB*c.REGS.LDA
		case O_RDIV:
			c.REGS.LRA = c.REGS.LDB/c.REGS.LDA
		case O_ADD:
			c.REGS.LRA = c.REGS.LDB+(arg).(int)
		case O_SUB:
			c.REGS.LRA = c.REGS.LDB-(arg).(int)
		case O_MUL:
			c.REGS.LRA = c.REGS.LDB*(arg).(int)
		case O_DIV:
			c.REGS.LRA = c.REGS.LDB/(arg).(int)
		case O_RAND:
			c.REGS.LRA = rint(0, (arg).(int))
		case O_RRAND:
			c.REGS.LRA = rint(c.REGS.LDA, c.REGS.LDB)
		case O_INC:
			if (OP.ArgType) { //mem
				c.MEM[(OP.Arg).(int)] = (c.MEM[(OP.Arg).(int)]).(int)+1
			} else {
				if (arg != 0) {
					c.ACC+=(arg).(int)
				} else {
					c.ACC++
				}
			}

		// mem
		case O_MEM:
			c.MEM = append(c.MEM, arg)
		case O_MEMDEL:
			c.MEM = append(c.MEM[:(arg).(int)], c.MEM[(arg).(int)+1:]...)
		case O_LRA2MEM:
			c.MEM = append(c.MEM, c.REGS.LRA)

		// branch
		case O_CMP:
			c.FLAGS.CMP = c.REGS.LDA == (arg).(int)
		case O_GOTO_LDX:
			c.ADDR = c.REGS.LDX // set -1
		case O_GOTO_LDH:
			c.ADDR = c.REGS.LDH // set -1
		case O_GOTO:
			c.ADDR = (arg).(int)-1
		case O_RGOTO:
			c.ADDR += (arg).(int)-1
		case O_JTL:
			c.ADDR = c.LABELS[(arg).(int)]
		case O_JIT: // jump to (label[) immd/ram (]) if flags.cmp
			if (c.FLAGS.CMP) {
				c.ADDR = c.LABELS[(arg).(int)]
			}
		case O_JTL_lDX:
			c.ADDR = c.LABELS[c.REGS.LDX]
		case O_JTL_lDH:
			c.ADDR = c.LABELS[c.REGS.LDH]

		// *&
		case O_INT2PRT:
			if (OP.ArgType) { //mem
				c.MEM = append(c.MEM, &(c.MEM[(OP.Arg).(int)]))
			} else {
				fprintf(stderr, "can't cast immedeate int to pointer\n")
				exit(1)
			}
		case O_PRT2INT:
			if (OP.ArgType) { //mem
				c.MEM = append(c.MEM, *(c.MEM[(OP.Arg).(int)]).(*interface{}))
			} else {
				fprintf(stderr, "can't cast immedeate pointer to int\n")
				exit(1)
			}

		// extra
		case O_NOP:
		case O_LDC:
			c.REGS.LDC = arg
		case O_CONV:
			c.STACK = append([]interface{}{converter[(arg).(int)](c.REGS.LDC)}, c.STACK...)


		// debug
		case O_DBGPRT:
			c.Print()

		// test helper
		case I_MakeLabel:
			c.LABELS[len(c.LABELS)] = (arg).(int)
		default:
			printf("\n%v\n", OP)
	}
}

