package Parser

// Original Grammar:
// 	Source = Statement { newline Statement } EOF .
// 	Statement = (UnaryInstruction | BinaryInstruction | LoadImmediate) right_arrow register .
// 	Operand = register | BinaryInstruction .
//	 	UnaryInstruction = Operand unary_operator .
// 	BinaryInstruction = Operand Operand binary_operator .
// 	LoadImmediate = op_paren immediate cl_paren .

// Simplified Grammar:
// 	Source = Source1 EOF .
// 	Source1 = Statement .
// 	Source1 = Statement newline Source1 .
// 	Statement = UnaryInstruction right_arrow register .
// 	Statement = BinaryInstruction right_arrow register .
// 	Statement = LoadImmediate right_arrow register .
// 	Operand = register .
// 	Operand = BinaryInstruction .
// 	UnaryInstruction = Operand unary_operator .
// 	BinaryInstruction = Operand Operand binary_operator .
// 	LoadImmediate = op_paren immediate cl_paren .

// Reverse Grammar:
// 	EOF Source1 -> Source
// 	Statement -> Source1
// 	Source1 newline Statement -> Source1
// 	register right_arrow UnaryInstruction -> Statement
// 	register right_arrow BinaryInstruction -> Statement
// 	register right_arrow LoadImmediate -> Statement
// 	register -> Operand
// 	BinaryInstruction -> Operand
// 	unary_operator Operand -> UnaryInstruction
// 	binary_operator Operand Operand -> BinaryInstruction
// 	cl_paren immediate op_paren -> LoadImmediate

// Shift-Reduce Grammar:
// 	[EOF] Source1 -> Source : Reduce
//
//	 	EOF [Source1] -> Source : Shift
// 	[Source1] newline Statement -> Source1 : Reduce
//
// 	[Statement] -> Source1 : Reduce
// 	Source1 newline [Statement] -> Source1 : Shift
//
// 	Source1 [newline] Statement -> Source1 : Shift
//
// 	[register] right_arrow UnaryInstruction -> Statement : Reduce
// 	[register] right_arrow BinaryInstruction -> Statement : Reduce
// 	[register] right_arrow LoadImmediate -> Statement : Reduce
// 	[register] -> Operand : Reduce
//
// 	register [right_arrow] UnaryInstruction -> Statement : Shift
// 	register [right_arrow] BinaryInstruction -> Statement : Shift
// 	register [right_arrow] LoadImmediate -> Statement : Shift
//
// 	register right_arrow [UnaryInstruction] -> Statement : Shift
//
// 	register right_arrow [BinaryInstruction] -> Statement : Shift
// 	[BinaryInstruction] -> Operand : Reduce
//
// 	register right_arrow [LoadImmediate] -> Statement : Shift
//
// 	[unary_operator] Operand -> UnaryInstruction : Reduce
//
// 	unary_operator [Operand] -> UnaryInstruction : Shift
// 	binary_operator [Operand] Operand -> BinaryInstruction : Shift
// 	binary_operator Operand [Operand] -> BinaryInstruction : Shift
//
// 	[binary_operator] Operand Operand -> BinaryInstruction : Reduce
//
// 	[cl_paren] immediate op_paren -> LoadImmediate : Reduce
//
// 	cl_paren [immediate] op_paren -> LoadImmediate : Shift
//
// 	cl_paren immediate [op_paren] -> LoadImmediate : Shift
