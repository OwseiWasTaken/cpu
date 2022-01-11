#! /usr/bin/python3.10
#imports
from util import *
import subprocess

I = 0
def iota():
	global I
	I+=1
	return I

OPS = {
	# stack (+LDA)
	"POP_LDA": iota(),
	"POP_LDB": iota(),
	"POP" : iota(),
	"PUSH" : iota(),
	"PUSH_LRA": iota(),

	# IOStream (+LDS)
	"OPEN"	 :	iota(),
	"CLOSE"  :	iota(),

	# regs (+ACC)
	"LRI"	 :	iota(),
	"LDA"	 :	iota(),
	"LDB"	 :	iota(),
	"LDS"	 :	iota(),
	"LDX"	 :	iota(),
	"LDH"	 :	iota(),
	"ACC2LDA":	iota(),

	# syscall (+LDS)
	"WRITE"  :	iota(),
	"WRITE_LRI":iota(),
	"WRITE_LRA":iota(),
	"FLUSH"  :	iota(),
	"EXIT"	 :	iota(),

	# math (+LDB)
	"RADD"	 :	iota(),
	"RSUB"	 :	iota(),
	"RMUL"	 :	iota(),
	"RDIV"	 :	iota(),
	"RMOD"	 :	iota(),
	"ADD"	 :	iota(),
	"SUB"	 :	iota(),
	"MUL"	 :	iota(),
	"DIV"	 :	iota(),
	"MOD"	 :	iota(),
	"INC"	 :	iota(),
	"RAND"	 :	iota(),
	"RRAND"  :	iota(),

	# mem (+LDA +LRA)
	"LRA2MEM":	iota(),
	"MEM"	 :	iota(),
	"MEMDEL" :	iota(),

	# branch
	"CMP"	 :	iota(),
	"GOTO_LDX": iota(),
	"GOTO_LDH": iota(),
	"GOTO"	 :	iota(),
	"RGOTO"  :	iota(),
	"JTL"	 :	iota(),
	"JIT"	 :	iota(),
	"JTL_lDX":	iota(),
	"JTL_lDH":	iota(),

	# debug
	"DBGPRT" :	iota(),

	# *&
	"INT2PRT" : iota(),
	"PRT2INT" : iota(),

	# test helper
	"I_MakeLabel":iota(),

	# extra
	"LDC" : iota(),
	"CONV": iota(),
	"NOP" : iota(),
	"OP_LEN" : iota(),
}

# TODO make REGS op when {lda}?

#Op Index To Name
OITN = {v:k for k, v in OPS.items()}

class Op:
	def __init__(this, OpName:str, ArgType:bool, Arg:Any, comment:str):
		this.Comment:str = comment
		this.OpId:int	 = OPS[OpName.upper()]
		this.ArgType:str = str(ArgType).lower()
		this.Arg:Any	 = Arg # normally int or string
	def __str__(this):
		return f"Op{{O_{OITN[this.OpId]}, {this.ArgType}, {this.Arg}}},//{this.Comment}"


#main
def Main() -> int:
	filename = get("").last
	if isinstance(filename, str) and exists(filename):
		with open(filename, 'r') as f:
			FILE = list(map(lambda x : TrimSpaces(x).replace("\t", "").replace('\n', ""), f.readlines()))
	else:
		fprintf(stderr, "reader.py didn't recieve any file to translate!\n")
		exit(1)

	TW = []

	# TODO
	labels:dict[str,int] = {}
	ops:list[Op] = []
	addr = 0
	for line in FILE:
		if not line:
			continue
		if line[-1] == ':':
			labels[line[:-1]] = addr+1
			#print(line, addr)
			continue
		addr += 1
	addr = 0
	for line in FILE:
		if not len(line):
			continue

		if line[0] == '/':
			continue

		ln = line.split()
		if not ln:
			continue

		#print(ln)
		commd = ln[0]

		if commd[-1] == ':':
			labels[commd] = addr
			continue

		ArgType = False
		Arg = 0
		comment = f" {addr}"

		if (len(ln)>1):
			if ln[1][0] == '[' and ln[1][-1] == ']':
				ArgType = true # mem
			elif ln[1][0] == '(' and ln[1][-1] == ')':
				ArgType = false
				Arg = labels[ln[1][1:-1]]-1
				comment+= f" -> ({ln[1][1:-1]})={Arg}"
			else:
				Arg = ln[1]
		ops.append(Op(commd, ArgType, Arg, comment))
		addr += 1

	stream = stdout
	ToFile = False
	if get('-o').exists:
		stream = open(get('-o').first, 'w')
		ToFile = True

	# write!
	stream.write(gotext % ('\n'.join([str(op) for op in ops])))
	#for op in ops:
	#	stream.write(str(op)+"\n")
	#stream.write(goend)


	if ToFile:
		stream.close()

	return 0

gotext= """
package main

include "cpu"

// OWEN worker

func main() {
InitRand()
InitGetCh()
CPU.CODE = []Op{
%s
}

for {
CPU.NextTick()
}
}
"""


#start
if __name__ == '__main__':
	start = tm()
	ExitCode = Main()

	if get('--debug').exists:
		if not ExitCode:
			printl("%scode successfully exited in " % COLOR.green)
		else:
			printl("%scode exited with error %d in " % (COLOR.red,ExitCode))
		print("%.3f seconds%s" % (tm()-start,COLOR.nc))
	exit(ExitCode)
