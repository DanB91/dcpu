// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

const AppTitle = 'dcpu-ide &#x2af9;&#x2afa;';

var stateTracker = null;
var workspace    = null;
var dashboard    = null;

window.onbeforeunload = function() {
	//TODO: implement clean shutdown.
};

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

	stateTracker = new StateTracker();
	if (!stateTracker.init()) {
		console.error("Failed to initialize state tracker.");
		return;
	}

	// Hook some events.
	document.onkeydown = onKeyDown;
	document.onkeyup = onKeyUp;
	document.onmousemove = onMouseMove;
	document.onmousedown = onMouseDown;
	document.onmouseup = onMouseUp;
	document.onmousewheel = onMouseWheel;
};

function onKeyDown (e)
{
	dashboard.onKey(e);
}

function onKeyUp (e)
{
	
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
