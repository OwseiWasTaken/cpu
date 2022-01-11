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
		Op{O_LDA, false, 0},// 0
		Op{O_LDB, false, 10},// 1
		Op{O_NOP, false, 0},// 2
		Op{O_DBGPRT, false, 0},
		Op{O_RRAND, false, 0},// 3
		Op{O_WRITE_LRA, false, 0},// 4
		Op{O_WRITE, false, "\n"},// 5
		Op{O_GOTO, false, 2},// 6 -> (Main)=2
		Op{O_EXIT, false, 0},// 7
	}
	printf("%d ops\n", OP_LEN)

	for {
		CPU.NextTick()
	}
	exit(0)
}
