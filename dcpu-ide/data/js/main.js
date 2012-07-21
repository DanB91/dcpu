// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

const AppTitle = 'dcpu-ide &#x2af9;&#x2afa;';

var stateTracker = null;
var workspace    = null;
var dashboard    = null;
var config       = null;

window.onbeforeunload = function()
{
	//TODO: implement clean shutdown.
};

window.onload = function ()
{
	// Find our UI elements.
	workspace = new Workspace();
	if (!workspace.init()) {
		console.error("Failed to initialize workspace.");
		return;
	}

	stateTracker = new StateTracker();
	if (!stateTracker.init()) {
		console.error("Failed to initialize state tracker.");
		return;
	}

	dashboard = new Dashboard();
	if (!dashboard.init()) {
		console.error("Failed to initialize dashboard.");
		return;
	}

	// Load configuration data.
	api.request({
		url: '/api/config',
		type: "json",
		onData : function (data)
		{
			config = data;
		},
		onError : function (msg, status)
		{
			console.error("Failed to load configuration data.");
		},
	});

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
	// Forward input to top-most modal dialog if need be.
	if (dialogs.length > 0) {
		dialogs[dialogs.length-1].onKey(e);
		return;
	}

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
