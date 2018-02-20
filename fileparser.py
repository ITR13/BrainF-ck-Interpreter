from typing import List as List
#0 = 0     #Zero Track Band      - Middle
DPBAND = 3 #Data Pointer Band    - Whole
DBAND  = 3 #Data Band            - Whole

IFBAND = 0 #Switch Test Band     - Right
PPBAND = 1 #Program Pointer Band - Right
PBAND  = 2 #Program Band         - Right

LCBAND = 0 #Loop Counter Band - Left
LGBAND = 1 #Loop Goal Band    - Left
CBAND  = 2 #Constant Band     - Left


legalPosTypes = [
	"Right",
	"Left",
	"Zero",
	"Program",
	"Loop",
	"Lost",
]

posRules = {
	"Right": {
		"Right": lambda x: x>=0,
		"Left": lambda _: False,
		"Zero": lambda x: x>0,
		"Program": lambda x: x+2>=0,
		"Loop": lambda x: False,
		"Lost": lambda _: False,
	},
	"Left": {
		"Right": lambda _: False,
		"Left": lambda x: x<=0,
		"Zero": lambda x: x<0,
		"Program": lambda x: False,
		"Loop": lambda x: x-2<=0,
		"Lost": lambda _: False,
	},
	"Zero": {
		"Right": lambda _: False,
		"Left": lambda _: False,
		"Zero": lambda x: x==0,
		"Program": lambda _: False,
		"Loop": lambda _: False,
		"Lost": lambda _: False,
	},
	"Program": {
		"Right": lambda _: False,
		"Left": lambda _: False,
		"Zero": lambda _: False,
		"Program": lambda x: x==0 or x==-1,
		"Loop": lambda _: False,
		"Lost": lambda _: False,
	},
	"Loop": {
		"Right": lambda _: False,
		"Left": lambda _: False,
		"Zero": lambda _: False,
		"Program": lambda _: False,
		"Loop": lambda x: x==0 or x==1,
		"Lost": lambda _: False,
	},
	"Lost": {
		"Right": lambda _: False,
		"Left": lambda _: False,
		"Zero": lambda _: False,
		"Program": lambda _: False,
		"Loop": lambda _: False,
		"Lost": lambda _: False,
	},
}

# Test that all posType combinations have a method
for i in legalPosTypes:
	for j in legalPosTypes:
		posRules[i][j](0)

bands = 5

class Graph:
	def __init__(
			self,
			name: str,
			inStatus: 'Status', 
			outStatus: 'Status' = None,
			subGraphs: List['Graph'] = [],
			specials: List[str] = [],
			code: List[str] = [],
	):
		print(f'Creating {name} - {inStatus}')
		self.name = name
		self.inStatus = inStatus
		self.outStatus = outStatus
		self.subGraphs = subGraphs
		self.specials = specials
		self.code = code


		if len(code)+len(specials) == 0 \
		or len(code)+len(subGraphs) == 0 \
		or len(specials)+len(subGraphs) == 0:
			pass
		else:
			raise ValueError("Illegal Graph")

		if len(specials) > 0:
			for i in specials:
				inStatus = parseSpecial(i, inStatus)
				print(f'\t{i} - {inStatus}')
			self.outStatus = inStatus

	# Checks a graph for conflict between itself and the first and last
	#  sub-graph, as well as inbetween each sub-graph
	# Returns "" on success, otherwise an error describing the conflict
	def Check(self) -> str:
		if len(self.specials)>0:
			return ""

		if len(self.code)>0:
			brackets = 0
			for i in self.code:
				for j in i:
					if j=='[':
						brackets+=1
					elif j==']':
						brackets-=1
					if brackets<0:
						return f'More ] than [ in {self.name}'
			if brackets > 0:
						return f'More [ than ] in {self.name}'
			return ""
		print(
			f"Checking {self.name} - "
			f"{self.inStatus} -> ... -> "
			f"{self.outStatus}"
		)
		outStatus = self.inStatus
		for i in self.subGraphs:
			print(f'\t{i.name} - {i.inStatus} -> {i.outStatus}')
			if not outStatus.isEqual(i.inStatus):
				return (
					f'Error in file {self.name} '
					f'with out {i.name}:\n'
					f'{outStatus} != {i.inStatus}'
				)
			outStatus = i.outStatus

		if not outStatus.isEqual(self.outStatus):
			return (
				f'Error with out in {self.name}:\n'
				f'{outStatus} != {self.outStatus}'
			)

		for i in self.subGraphs:
			err = i.Check()
			if err != "":
				return err
		return ""

	def Compile(self, comment: bool):
		if comment:
			return (
				f'({self.name} {self.inStatus})'
				f' {self.__compile__(comment)} '
				f'(/{self.name} {self.outStatus})'
			)
		return self.__compile__(comment)

	def __compile__(self,comment) -> str:
		if len(self.specials) != 0:
			return "".join([
				specials_code[i] for i in self.specials
			])
		if len(self.subGraphs) != 0:
			return "".join([
				i.Compile(comment) for i in self.subGraphs
			])
		return "".join(self.code)
		
		
class Switch(Graph):
	def __init__(
		self,
		name: str, 
		inStatus: 'Status', 
		keys: List[str], 
		values: List[str],
	):
		self.name = name
		self.inStatus = inStatus
		self.outStatus = Status(0, 'Zero', 0)
		self.keys = [
			ord(i[1]) if i[0] == "'" else int(i)
			for i in keys
		]
		self.values = [ 
				readfile(i, Status(0,"Right", 0))
				for i in values
		]

	def Check(self) -> bool:
		if self.inStatus.posType == "Zero":
			if self.inStatus.band == 0 \
			and self.inStatus.posIndex >= 0:
				return (
					f'Illegal in status for '
					f'switch: {self.inStatus}'
				)
		elif not self.inStatus.isEqual(Status(
			self.inStatus.band,
			"Right",
			0
		)):
			return (
				f'Illegal in status for switch: {self.inStatus}'
			)
			

		for i in self.values:
			if not Status(0,"Right",2).isEqual(i.inStatus):
				return (
					f'Error in switch with in {i.name}\n'
					f'(b0 Right:2) != {i.inStatus})'
				)
			if not i.outStatus.isEqual(Status(0,"Right",1)):
				return (
					f'Error in switch with out {i.name}\n'
					f'{i.outStatus} != (b0 Right:1)'
				)

			err = i.Check()
			if err != "":
				return err

	def __compile__(self, comment) -> str:
		moveDist = self.inStatus.band
		down = "<"*moveDist
		up = ">"*moveDist
		right = ">"*bands
		left = "<"*bands
		run = ""
		prev = 0
		for i,j in zip(self.keys,self.values):
			run = run + (
				#Decrease to value we want to test, and set unset bit
				#Also set helper bit to help with aligning
				f'{"-"*((i-prev)%256)}{right}+{right*2}+{left*3}'
				#If not zero, set unset bit to zero, then move right
				# twice, so that if zero is on the empty bit, and
				# if not zero is on the helper bit. Then scan left so
				# that if not zero lands on the empty bit and if zero
				# stands still. Then move right and unset the helper
				# bit and move left twice to land on the unset bit
				f'[{right}-]{right*2}[{left}]{right}-{left*2}'
				#If the unset bit has not been set then run the
				# corresponding program. Move one left to be back on
				# the counter
				f'[-{j.Compile(comment)}]{left}'
			)
			prev = i

		if self.inStatus.posType=="Zero":
			xDiff = 2-self.inStatus.posIndex
			if xDiff<0:
				down += "<"*bands*(-xDiff)
				up += ">"*bands*(-xDiff)
			else:
				down += ">"*bands*xDiff
				up += "<"*bands*xDiff

		return (
			# Copy to diamond band twice
			f'[-{down}+{right}+{left+up}]'
			# Copy back to original pos using the duplicate
			f'{down+right}[-{left+up}+{down+right}]'
			# Go to copy, apply test
			f'{left+run}'
			# Increase to zero then search back
			f'[+]-[+{"<"*5}-]+'
		)

class Loop(Graph):
	def __init__(
		self,
		inStatus: 'Status',
		subGraphs: List[Graph],
	):
		self.name: str = "Loop"
		self.inStatus = inStatus
		self.outStatus = inStatus
		self.subGraphs = subGraphs
		self.code = []
		self.specials = []

	def __compile__(self, comment) -> str:
		return f'[{super().__compile__(comment)}]'

class Status:
	def __init__(self, band: int, posType: str, posIndex: int):
		self.band = band
		self.posType = posType
		self.posIndex = posIndex
	
	def __repr__(self):
		return (
			f'(b{self.band} {self.posType}:{self.posIndex})'
		)

	def isEqual(self, other: 'Status') -> bool:
		if self.band != other.band:
			return False
		if  self.posType == other.posType \
		and self.posIndex == other.posIndex:
			return True
		return posRules[other.posType][self.posType](
			self.posIndex-other.posIndex
		)

	def movePos(self, offset):
		return Status(self.band, self.posType, self.posIndex+offset)

	def moveBand(self, offset):
		return Status(self.band+offset, self.posType, self.posIndex)

specials_code = {
	"Right": ">"*bands,
	"Left": "<"*bands,
	"Add": "+",
	"Sub": "-",
	"Up": ">",
	"Down": "<",
	"Inv": "+"*128,
	"SR": "["+">"*bands+"]",
	"SL": "["+"<"*bands+"]",
	"SZL": "-[+"+"<"*5+"-]+",
	"SZR": "-[+"+">"*5+"-]+",
	"Write": ".",
	"Read": ",",
	"Clear": "[-]"
}

def parseSpecial(special: str, status: 'Status'):
	if special=="Right":
		return status.movePos(1)
	if special=="Left":
		return status.movePos(-1)
	if special=="Add"   \
	or special=="Sub"   \
	or special=="Inv"   \
	or special=="Read"  \
	or special=="Print" \
	or special=="Clear":
		return status
	if special=="Up":
		return status.moveBand(1)
	if special=="Down":
		return status.moveBand(-1)
	if special=="SZL":
		if status.isEqual(Status(0,"Right",0)):
			return Status(0,"Zero",0)
		return Status(0,"Lost",0)
	if special=="SZR":
		if status.isEqual(Status(0,"Left",0)):
			return Status(0,"Zero",0)
		return Status(0,"Lost",0)

	if special=="SR":
		if status.band==PBAND:
			if status.posType!="Zero" or status.posIndex<1:
				return Status(0,"Lost",0)
			return Status(PBAND,"Right",0)
		if status.band==PPBAND:
			if status.posType!="Zero" or status.posIndex!=2:
				return Status(0,"Lost",0)
			return Status(PPBAND, "Program", 1)
		if status.band==LGBAND or status.band==LCBAND:
			if posRules["Loop"][status.posType](status.posIndex):
				return Status(status.band,"Zero",-1)
			return Status(0,"Lost",0)
		return Status(0,"Lost",0)

	if special=="SL":
		if status.band==LGBAND or status.band==LCBAND:
			if status.posType!="Zero" or status.posIndex!=-2:
				return Status(0,"Lost",0)
			return Status(status.band, "Loop", -1)
		if status.band==PPBAND:
			if posRules["Program"][status.posType](status.posIndex):
				return Status(PPBAND, "ZERO", 1)
			return Status(0,"Lost",0)
		return Status(0,"Lost",0)

	raise ValueError(f'Illegal special {special}')

def parseStatus(name: str, band: str, posType: str, posIndex: str):
	if not posType in legalPosTypes:
		raise ValueError(f'Illegal posType {posType} in {name}') 
	return Status(int(band),posType,int(posIndex))

def parse_name(name: List[str], inStatus: 'Status') -> 'Graph':
	if name[0] in specials_code.keys():
		return Graph(name[0], inStatus, None, [], [name[0]], [])
	
	elif name[0] == "Switch":
		return Switch(name[0], inStatus, name[1::2], name[2::2])

	elif name[0] == "Move":
		return move_to(name[1], inStatus)

	elif name[0] == "Loop":
		return Loop(inStatus, [readfile(name[1], inStatus)])

	return readfile(name[0], inStatus)


fixedPoints = {
	"0": (0,0),
	"LoopBit": (CBAND, -1),
	"ErrorBit": (CBAND, -2),
	"DataBit": (CBAND, -3),
	"ProgramPointer": (PPBAND, 2),
	"Program": (PBAND, 2), #2 to avoid crashing switch and zero
	"DataPointer": (DPBAND, 0),
	"Data": (DBAND, 0),
	"LoopGoal": (LGBAND, -2), #-2 to avoid crashing loopcount and zero
	"LoopCount": (LCBAND, -2),
}

def move_to(location: str, inStatus: 'Status') -> 'Graph':
	if inStatus.posType != "Zero":
		return Graph("Lost", Status(0,"Lost",0),  Status(0,"Lost",0))
	y, x = fixedPoints[location]
	x = x-inStatus.posIndex
	y = y-inStatus.band

	path: List[str] = []
	if x < 0:
		path.extend(["Left"]*(-x))
	else:
		path.extend(["Right"]*x)
	
	if y < 0:
		path.extend(["Down"]*(-y))
	else:
		path.extend(["Up"]*y)
	g = Graph("Move", inStatus, inStatus, [], path, [])
	return g

def readfile(name: str, inStatus: 'Status') -> Graph:
	print(f'Reading {name} - {inStatus}')

	path: str = f'./Graph/{name}.ng'

	lines: List[str] = []
	with open(path) as f:
		lines = f.read().splitlines()	
	nodeType: str = lines[0]
	if nodeType == "Special":
		return Graph(name, inStatus, inStatus, [], lines[1:], [])
	
	inS: List[str] = lines[1].split()
	outS: List[str] = lines[2].split()
	inStatus: 'Status' = parseStatus(name, inS[0], inS[1], inS[2])
	outStatus: 'Status' = parseStatus(name, outS[0], outS[1], outS[2])
	
	if nodeType == "Code":
		return Graph(name, inStatus, outStatus, [], [], lines[3:])
	if nodeType == "Sequence":
		graphs: List[Graph] = []
		tempInStatus: 'Status' = inStatus
		for i in lines[3:]:
			if i != "":
				graph = parse_name(i.split(), tempInStatus)
				tempInStatus = graph.outStatus
				graphs.append(graph)
		return Graph(name, inStatus, outStatus, graphs, [], [])
	else:
		raise ValueError(f'Illegal nodeType in {name}')

main_graph: 'Graph' = readfile("main", Status(0,"Zero",0))
graph_error: str = main_graph.Check()

if graph_error:
	print(graph_error)
else:
	print("Successfully validated graph!")
	program: str = main_graph.Compile(False)
	print(program)
	with open("./compiled.bf", "w") as f:
		f.write(program)
	program: str = main_graph.Compile(True)
	with open("./commented.bf", "w") as f:
		f.write(program)


