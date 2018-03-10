# Note: These files are disorganized and have circular dependencies,
#  I really should clean them up and write them better, but seeing how
#  the main purpose was to compile the compiler, and there isn't any
#  more for me to add to these files, I don't really feel like it
from typing import List as List

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
				inStatus = parse_special(i, inStatus)
				print(f'\t{i} - {inStatus}')
			self.outStatus = inStatus

	# Checks a graph for conflict between itself and the first and last
	#  sub-graph, as well as inbetween each sub-graph
	# Returns "" on success, otherwise an error describing the conflict
	def check(self) -> str:
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
			err = i.check()
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
				read_file(i, Status(0,"Right", 0))
				for i in values
		]

	def check(self) -> bool:
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

			err = i.check()
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
				f'{"-"*((i-prev)%256)}{right}+{left}'
				#If not zero, set unset bit to zero, then move right,
				# so that if zero is on the unset bit, and if not zero
				# is on the empty bit. Then scan right so that if not 
				# zero lands on the empty bit and if zero stands still. 
				#Then move left to put both on the unset-bit
				f'[{right}-]{right}[{right}]{left}'
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
		return pos_rules[other.posType][self.posType](
			self.posIndex-other.posIndex
		)

	def movePos(self, offset):
		return Status(self.band, self.posType, self.posIndex+offset)

	def moveBand(self, offset):
		return Status(self.band+offset, self.posType, self.posIndex)

from constants import *
from methods import *
