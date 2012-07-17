// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

var menu      = null;
var workspace = null;
var editor    = null;

window.onload = function ()
{
	// Find our UI elements.
	if ((menu = document.getElementById('menu')) == null) {
		console.error("Failed to acquire menu element.");
		return;
	}

	if ((workspace = document.getElementById('workspace')) == null) {
		console.error("Failed to acquire workspace element.");
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
		toggleMenu();
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

function toggleMenu()
{
	var m = fx.metrics(menu);
	var hide = m.top == 0;

	fx.show(menu)
	  .slideTo({
		node:     menu,
		top:      hide ? -m.height : 0,
		duration: 500,
		onFinish: function() {
			if (hide) {
				fx.hide(menu);
			}
		},
	});
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
