
const (
	// stack (+LDA)
	O_POP_LDA = iota
	O_POP_LDB = iota
	O_POP	  = iota
	O_PUSH	  = iota
	O_PUSH_LRA= iota

	// IOStream (+LDS)
	O_OPEN	 =	iota
	O_CLOSE  =	iota

	// regs (+ACC)
	O_LRI	 =	iota
	O_LDA	 =	iota
	O_LDB	 =	iota
	O_LDS	 =	iota
	O_LDX	 =	iota
	O_LDH	 =	iota
	O_ACC2LDA=	iota

	// syscall (+LDS)
	O_WRITE  =	iota
	O_WRITE_LRI=iota
	O_WRITE_LRA=iota
	O_FLUSH  =	iota
	O_EXIT	 =	iota

	// math (+LDB)
	O_RADD	 =	iota
	O_RSUB	 =	iota
	O_RMUL	 =	iota
	O_RDIV	 =	iota
	O_RMOD	 =	iota
	O_ADD	 =	iota
	O_SUB	 =	iota
	O_MUL	 =	iota
	O_DIV	 =	iota
	O_MOD	 =	iota
	O_INC	 =	iota
	O_RAND	 =	iota
	O_RRAND  =	iota

	// mem (+LDA +LRA)
	O_LRA2MEM=	iota
	O_MEM	 =	iota
	O_MEMDEL =	iota

	// branch
	O_CMP	 =	iota
	O_GOTO_LDX= iota
	O_GOTO_LDH= iota
	O_GOTO	 =	iota
	O_RGOTO  =	iota
	O_JTL	 =	iota
	O_JIT	 =	iota
	O_JTL_lDX=	iota
	O_JTL_lDH=	iota

	// debug
	O_DBGPRT =	iota

	// &*
	O_INT2PRT = iota
	O_PRT2INT = iota

	// test helper
	I_MakeLabel=iota

	// extra
	O_NOP  = iota
	OP_LEN = iota
)
