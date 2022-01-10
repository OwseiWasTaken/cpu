
package main

include "gutil"
include "cpu"

// OWEN worker

func main() {
	InitRand()
	InitGetCh()
	var Code []Op = []Op{
Op{O_MEM, false, 1},// 0
Op{O_INC, true, 0},// 1
Op{O_WRITE, true, 0},// 2
Op{O_WRITE, false, "\n"},// 3
Op{O_GOTO, false, 1},// 4 -> (Main)=1

}
CPU.CODE = Code

var CL int = len(Code)

for ;CPU.ADDR<CL;{
CPU.NextTick()
}
exit(0)
}
