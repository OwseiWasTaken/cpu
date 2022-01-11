package main

include "gutil"
include "./cpu"

// OWEN worker


func main() {
	Init()
	InitGetCh()
	CPU.CODE = []Op{
		CPU.MakeOp(O_ADD, 5, false, 7),
		CPU.MakeOp(O_WRITE_LRA),
		CPU.MakeOp(O_WRITENL),
		CPU.MakeOp(O_WRITE, "poggers\n"),
		CPU.MakeOp(),
	}
	printf("%d ops\n", OP_LEN)

	for {
		CPU.NextTick()
	}
}
