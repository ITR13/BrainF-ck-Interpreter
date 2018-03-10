empties = {
	'+': {'+': 1, '-': -1},
	'-': {'-': 1, '+': -1},
	'>': {'>': 1, '<': -1},
	'<': {'<': 1, '>': -1},
}

def remove_redundant( program: str ) -> str:
	out = ""
	collect = ""
	count = 0
	for i in program:
		if collect == "":
			collect = i
			count = 1
		elif collect in empties:
			v = empties[collect]
			if i in v:
				count += v[i]
				if count == 0:
					collect = ""
			else:
				out += collect*count
				collect = i
				count = 1
		else:
			out += collect*count
			collect = i
			count = 1
	
	out += collect*count
	return out
	
def remove_trailing( program: str ) -> str:
	out = ""
	collect = ""
	for i in program:
		if i in empties:
			collect += i
		else:
			out += collect+i
			collect = ""
	
	return out