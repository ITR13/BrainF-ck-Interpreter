# Note: These files are disorganized and have circular dependencies,
#  I really should clean them up and write them better, but seeing how
#  the main purpose was to compile the compiler, and there isn't any
#  more for me to add to these files, I don't really feel like it
from typing import List as List

def parse_name(name: List[str], inStatus: 'Status') -> 'Graph':
	if name[0] in specials_code.keys():
		return Graph(name[0], inStatus, None, [], [name[0]], [])
	
	elif name[0] == "Switch":
		return Switch(name[0], inStatus, name[1::2], name[2::2])

	elif name[0] == "Move":
		return move_to(name[1], inStatus)

	elif name[0] == "Loop":
		return Loop(inStatus, [read_file(name[1], inStatus)])

	return read_file(name[0], inStatus)

def parse_special(special: str, status: 'Status'):
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
			if pos_rules["Loop"][status.posType](status.posIndex):
				return Status(status.band,"Zero",-1)
			return Status(0,"Lost",0)
		return Status(0,"Lost",0)

	if special=="SL":
		if status.band==LGBAND or status.band==LCBAND:
			if status.posType!="Zero" or status.posIndex!=-2:
				return Status(0,"Lost",0)
			return Status(status.band, "Loop", -1)
		if status.band==PPBAND:
			if pos_rules["Program"][status.posType](status.posIndex):
				return Status(PPBAND, "ZERO", 1)
			return Status(0,"Lost",0)
		return Status(0,"Lost",0)

	raise ValueError(f'Illegal special {special}')

def parse_status(name: str, band: str, posType: str, posIndex: str):
	if not posType in legalPosTypes:
		raise ValueError(f'Illegal posType {posType} in {name}') 
	return Status(int(band),posType,int(posIndex))

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

def read_file(name: str, inStatus: 'Status') -> 'Graph':
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
	inStatus: 'Status' = parse_status(name, inS[0], inS[1], inS[2])
	outStatus: 'Status' = parse_status(name, outS[0], outS[1], outS[2])
	
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

from classes import *