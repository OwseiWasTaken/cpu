package main

include "gutil"
include "./cpu"

// OWEN worker
// test.go is used to test cpu.go!
// if you want to execute .gr files:
//	   $ reader.py (FILE.gr) -o (OUT.go)
//	   $ gc OUT.go
//	   $ ./OUT


func main() {
	InitRand()
	InitGetCh()
	CPU.CODE = []Op{
		Op{O_MEM, false, 10},	   // 0
		Op{O_INT2PRT, true, 0},    // 1
		Op{O_DBGPRT, false, 0},    // 2
		Op{O_INC, true, 0},		   // 3
		Op{O_DBGPRT, false, 0},    // 4
		Op{O_EXIT, false, 0},	   // 5
	}
	printf("%d ops\n", OP_LEN)

	for {
		CPU.NextTick()
	}
	exit(0)
}
