// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

const AppTitle = 'dcpu-ide &#x2af9;&#x2afa;';

var dialogs      = [];
var stateTracker = null;
var workspace    = null;
var dashboard    = null;
var config       = null;
var project      = null;

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
	try {
		config = api.request({
			url: '/api/config',
			type: "json",
			async: false,
		});
	} catch(e) {
		console.error('window.onload: ', e.message);	
	}

	// Hook some events.
	document.onkeydown = onKeyDown;
	document.onkeyup = onKeyUp;
	document.onmousemove = onMouseMove;
	document.onmousedown = onMouseDown;
	document.onmouseup = onMouseUp;
	document.onmousewheel = onMouseWheel;
};

// lockApplication locks various nodes to prevent tab input to them,
// This is called by the given modal dialog when it is opened.
function lockApplication(dlg)
{
	var nodelists = [];

	if (dialogs.length == 0) {
		nodelists = nodelists.concat(
			findLockableNodes(dashboard.node),
			findLockableNodes(workspace.node)
		)
	} else {
		nodelists = nodelists.concat(
			findLockableNodes(dialogs[dialogs.length-1].node)
		)
	}

	for (var i = 0; i < nodelists.length; i++) {
		for (var j = 0; j < nodelists[i].length; j++) {
			lock(nodelists[i][j]);
		}
	}

	// Update a list of open dialogs.
	dialogs.push(dlg);
}

// unlockApplication unlocks various nodes to re-enable tab input to them,
// This is called by the given modal dialog when it is closed.
function unlockApplication()
{
	// Update a list of open dialogs.
	dialogs.pop();

	var nodelists = [];

	if (dialogs.length == 0) {
		nodelists = [].concat(
			findLockableNodes(dashboard.node),
			findLockableNodes(workspace.node)
		)
	} else {
		nodelists = nodelists.concat(
			findLockableNodes(dialogs[dialogs.length-1].node)
		)
	}

	for (var i = 0; i < nodelists.length; i++) {
		for (var j = 0; j < nodelists[i].length; j++) {
			unlock(nodelists[i][j]);
		}
	}
}

// findLockableNodes finds all children of the given node
// which we consider relevant for locking/unlocking.
function findLockableNodes(e)
{
	return [].concat(
		e.getElementsByTagName('a'),
		e.getElementsByTagName('select'),
		e.getElementsByTagName('input'),
		e.getElementsByTagName('button'),
		e.getElementsByTagName('textarea')
	);
}

// lock locks the given control.
// It is disabled to prevent it from getting key input.
function lock (e)
{
	if (e == undefined || e.attributes['disabled'] != undefined) {
		return; // We don't want to unlock it later when it shouldn't be.
	}

	e.setAttribute('disabled', 'disabled');
	e._locked = true;
}

// unlock unlocks the given control.
// It is (re-)enabled to allow key input.
function unlock (e)
{
	if (e != undefined && e._locked == true) {
		e.removeAttribute('disabled');
		e._locked = false;
	}
}

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

String.prototype.ltrim = function()
{
	return this.replace(/^\W+/, '')
}

String.prototype.rtrim = function()
{
	return this.replace(/\W+$/, '')
}

String.prototype.trim = function()
{
	return this.replace(/(^\W+)|(\W+$)/, '')
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
