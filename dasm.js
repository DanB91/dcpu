CodeMirror.defineMode("dasm", function(config, parserConfig) {
	var isOperatorChar = /[+\-*&^%=<>!|\/]/;
	var indentUnit = config.indentUnit;

	var keywords = {
		"set":0, "add":0, "sub":0, "mul":0, "mli":0, "div":0,
		"dvi":0, "mod":0, "mdi":0, "and":0, "bor":0, "xor":0,
		"shr":0, "asr":0, "shl":0, "ifb":0, "ifc":0, "ife":0,
		"ifn":0, "ifg":0, "ifa":0, "ifl":0, "ifu":0, "adx":0,
		"sbx":0, "sti":0, "std":0, "jsr":0, "int":0, "iag":0,
		"ias":0, "rfi":0, "iaq":0, "hwn":0, "hwq":0, "hwi":0,
		"dat":0, "panic":0, "exit":0, "def":0, "end":0, "return":0
	};

	var atoms = {
		"a":0, "b":0, "c":0, "x":0, "y":0, "z":0, "i":0,
		"j":0, "ia":0,  "ex":0, "peek":0, "push":0, "pop":0,
		"pc":0, "sp":0
	};

	var branches = {
		"ifb":0, "ifc":0, "ife":0, "ifn":0,
		"ifg":0, "ifa":0, "ifl":0, "ifu":0
	};

	function tokenBase(stream, state) {
		var ch = stream.next();

		if (ch == '"' || ch == "'") {
			state.tokenize = tokenString(ch);
			return state.tokenize(stream, state);
		}

		if (ch == ";") {
			stream.skipToEnd();
			return "comment";
		}

		if (/[\d\.]/.test(ch)) {
			if (ch == ".") {
				stream.match(/^[0-9]+([eE][\-+]?[0-9]+)?/);
			} else if (ch == "0") {
				stream.match(/^[xX][0-9a-fA-F]+/) ||
				stream.match(/^[bB][01]+/) ||
				stream.match(/^[0-7]+/);
			} else {
				stream.match(/^[0-9]*\.?[0-9]*([eE][\-+]?[0-9]+)?/);
			}

			return "number";
		}

		if (/[\[\],\:]/.test(ch)) {
			curPunc = ch;
			return null
		}

		if (isOperatorChar.test(ch)) {
			stream.eatWhile(isOperatorChar);
			return "operator";
		}

		stream.eatWhile(/[\w\$_]/);

		var cur = stream.current();
		if (keywords.propertyIsEnumerable(cur)) {
			if (branches.propertyIsEnumerable(cur)) {
				state.branchDepth++;
			} else {
				state.branchDepth = 1;
			}
			return "keyword";
		}

		if (atoms.propertyIsEnumerable(cur)) {
			return "atom";
		}

		return "word";
	}

	function tokenString(quote) {
		return function(stream, state) {
			var escaped = false, next;

			while ((next = stream.next()) != null) {
				if (next == quote && !escaped) {
					state.tokenize = tokenBase;
					break;
				}
				
				escaped = !escaped && next == "\\";
			}

			return "string";
		};
	}

	return {
		startState: function(basecolumn) {
			return {
				branchDepth: 0,
				indented: indentUnit,
			};
		},

		token: function(stream, state) {
			if (stream.eatSpace()) {
				return null;
			}

			return (state.tokenize || tokenBase)(stream, state);
		},

		indent: function(state, textAfter) {
			var ch = textAfter.charAt(0);
			if (ch == ':' || ch == ';') {
				return 0;
			}

			var n = state.branchDepth * indentUnit;
			if (state.indented != n) {
				state.indented = n;
			}

			return state.indented;
		},

		electricChars: ":"
	};
});

CodeMirror.defineMIME("text/x-dasm", "dasm");
