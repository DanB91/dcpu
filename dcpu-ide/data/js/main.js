// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

var workspace = null;
var editor    = null;

window.onload = function ()
{
	workspace = document.getElementById('workspace');

	if (!workspace) {
		console.error("Failed to acquire workspace element.");
		return;
	}
}

/*
	editor = CodeMirror(workspace, {
		mode: "dasm",
		theme: "eclipse",
		indentUnit: 3,
		tabSize: 3,
		smartIndent: true,
		indentWithTabs: false,
		electricChars: true,
		autoClearEmptyLines: false,
		lineWrapping: false,
		lineNumbers: true,
		firstLineNumber: 1,
		gutter: true,
		fixedGutter: true,
		matchBrackets: false,
		tabindex: 1,
		value: 	"   jsr main\n" + 
				"   sub pc, 1 ; halt CPU\n" + 
				"\n" + 
				"; entry point\n" + 
				"def main\n" + 
				"   \n" + 
				"end\n"
	});

	if (!editor) {
		console.error("Failed to initialize CodeMirror.");
		return;
	}
*/
