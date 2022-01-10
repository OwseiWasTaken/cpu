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
		// make out.txt file (if ArgType was true, LDS would be set to out.txt)
		Op{O_OPEN, false, "out.txt"},

		// set stream to out.txt
		Op{O_LDS, false, 2},

		// write "poggers\n" to out.txt stream
		Op{O_WRITE, false, "poggers\n"},

		// flush out.txt stream
		Op{O_FLUSH, false, 0},

		Op{O_EXIT, false, 0}, // -1
	}
	printf("%d ops\n", OP_LEN)

	for {
		CPU.NextTick()
	}
	exit(0)
}
