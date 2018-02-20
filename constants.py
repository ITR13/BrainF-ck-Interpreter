# Note: These files are disorganized and have circular dependencies,
#  I really should clean them up and write them better, but seeing how
#  the main purpose was to compile the compiler, and there isn't any
#  more for me to add to these files, I don't really feel like it

bands = 5

#0 = 0     #Zero Track Band      - Middle
DPBAND = 3 #Data Pointer Band    - Whole
DBAND  = 3 #Data Band            - Whole

IFBAND = 0 #Switch Test Band     - Right
PPBAND = 1 #Program Pointer Band - Right
PBAND  = 2 #Program Band         - Right

LCBAND = 0 #Loop Counter Band - Left
LGBAND = 1 #Loop Goal Band    - Left
CBAND  = 2 #Constant Band     - Left

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

legalPosTypes = [
	"Right",
	"Left",
	"Zero",
	"Program",
	"Loop",
	"Lost",
]

pos_rules = {
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
		pos_rules[i][j](0)

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