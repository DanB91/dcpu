// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

var workspace = null;
var dashboard = null;

window.onload = function ()
{
	// Find our UI elements.
	if ((workspace = document.getElementById('workspace')) == null) {
		console.error("Failed to acquire workspace element.");
		return;
	}

	dashboard = new Dashboard();
	if (!dashboard.init()) {
		console.error("Failed to initialize dashboard.");
		return;
	}

	// Hook some events.
	document.onkeydown = onKeyDown;
	document.onkeyup = onKeyUp;
	document.onmousemove = onMouseMove;
	document.onmousedown = onMouseDown;
	document.onmouseup = onMouseUp;
	document.onmousewheel = onMouseWheel;
}

function onKeyDown (e)
{
	var key = (e.which != 0) ? e.which : e.keyCode;
	
	switch (key) {
	case 192: // ~
		dashboard.toggle();
		break;
	}
}

function onKeyUp (e)
{
	var key = (e.which != 0) ? e.which : e.keyCode;
}

function onMouseMove (e)
{
	
}

function onMouseDown (e)
{
	
}

function onMouseUp (e)
{
	
}

function onMouseWheel (e)
{
	
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
