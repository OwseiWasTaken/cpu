
// jump example
var Code []Op = []Op{
	Op{O_RAND, false, 3},		// 0	(3)
	Op{O_LRA2MEM, false, 0},	// 1
	Op{O_LDB, true, 0},         // 2
	Op{O_ADD, false, 1},        // 3
	Op{O_MEMDEL, false, 0},     // 4
	Op{O_LRA2MEM, false, 0},	// 5
	Op{O_RGOTO, true, 0},		// 6 -> +(0..3)
	Op{O_WRITE, false, "1!\n"}, // 7  (skip?)
	Op{O_WRITE, false, "2!\n"}, // 8  (skip?)
	Op{O_WRITE, false, "3!\n"}, // 8  (skip?)
	Op{O_EXIT, false, 0},		// 9
}
