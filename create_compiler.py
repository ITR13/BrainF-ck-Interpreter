from constants import *
from classes import *
from methods import *
import shortener

main_graph: 'Graph' = read_file("main", Status(0,"Zero",0))
graph_error: str = main_graph.check()

if graph_error:
	print(graph_error)
else:
	print("Successfully validated graph!")
	
	program: str = main_graph.Compile(False)
	with open("./compiled.bf", "w") as f:
		f.write(program)
	
	program: str = shortener.remove_redundant(program)
	program: str = shortener.remove_trailing(program)
	print(program)
	with open("./shortened.bf", "w") as f:
		f.write(program)
	
	program: str = main_graph.Compile(True)
	with open("./commented.bf", "w") as f:
		f.write(program)
	