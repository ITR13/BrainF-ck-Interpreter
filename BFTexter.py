def write(text: str, start: int) -> str:
	prev: int = start
	out = ""
	for i in text:
		out = out + getinc(ord(i),prev) + "."
		prev = ord(i)
	return out

def getinc(goal: int, prev: int) -> str:
	diff = goal-prev
	if diff > 128:
		diff = diff-256
	elif diff < -128:
		diff = diff+256

	if diff < 0:
		return "-"*(-diff)
	return "+"*diff

