#! /usr/bin/python3.10
#imports
from util import *
import subprocess

class _Enum:
	def __init__(this):
		this.I= -1
	@property
	def iota(this):
		this.I+=1
		return this.I
Enum = _Enum()

OPS = {
	# stack (+LDA)
	"POP"	 :	Enum.iota,
	"PUSH"	 :	Enum.iota,
	# IOStream (+LDS)
	"OPEN"	 :	Enum.iota,
	"CLOSE"  :	Enum.iota,
	# regs (+ACC)
	"LRI"	 :	Enum.iota,
	"LDA"	 :	Enum.iota,
	"LDB"	 :	Enum.iota,
	"LDS"	 :	Enum.iota,
	"LDX"	 :	Enum.iota,
	"LDH"	 :	Enum.iota,
	"ACC2LDA":	Enum.iota,
	# syscall (+LDS)
	"WRITE"  :	Enum.iota,
	"WRITE_LRI":Enum.iota,
	"FLUSH"  :	Enum.iota,
	"EXIT"	 :	Enum.iota,
	# math (+LDB)
	"ADD"	 :	Enum.iota,
	"SUB"	 :	Enum.iota,
	"MUL"	 :	Enum.iota,
	"DIV"	 :	Enum.iota,
	"MOD"	 :	Enum.iota,
	"INC"	 :	Enum.iota,
	"RAND"	 :	Enum.iota,
	"RRAND"  :	Enum.iota,
	# mem (+LDA +LRA)
	"LDA2MEM":	Enum.iota,
	"LDMEM"  :	Enum.iota,
	"LRA2MEM":	Enum.iota,
	"MEM"	 :	Enum.iota,
	"MEMDEL" :	Enum.iota,
	# branch
	"CMP"	 :	Enum.iota,
	"GOTO_LDX": Enum.iota,
	"GOTO_LDH": Enum.iota,
	"GOTO"	 :	Enum.iota,
	"RGOTO"  :	Enum.iota,
	"JTL"	 :	Enum.iota,
	"JIT"	 :	Enum.iota,
	"JTL_lDX":	Enum.iota,
	"JTL_lDH":	Enum.iota,
	# debug
	"DBGPRT" :	Enum.iota,
	# &*
	"INT2PRT" : Enum.iota,
	# reader
	"I_MakeLabel":Enum.iota,
	# extra
	"_LEN" : Enum.iota,
}
# TODO make REGS op when (ldx)
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

include "gutil"
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
exit(0)
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
