// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

const AppTitle = 'dcpu-ide &#x2af9;&#x2afa;';

var dialogs   = [];
var workspace = null;
var dashboard = null;
var socket    = null;
var config    = null;
var project   = null;

window.onload = function ()
{
	// Hook some events.
	document.onkeydown = onKeyDown;
	document.onkeyup = onKeyUp;
	document.onmousemove = onMouseMove;
	document.onmousedown = onMouseDown;
	document.onmouseup = onMouseUp;
	document.onmousewheel = onMouseWheel;
	document.onselectstart = function () { return false; };

	dashboard = new Dashboard();
	if (!dashboard.init()) {
		(new ErrorDialog())
			.content('Failed to initialize dashboard.')
			.open();
		return;
	}

	workspace = new Workspace();
	if (!workspace.init()) {
		(new ErrorDialog())
			.content('Failed to initialize workspace.')
			.open();
		return;
	}

	// initialize socket.
	socket = new Socket({
		onopen: function ()
		{
			apiHandshake(socket);
		},
		onclose: function ()
		{
			socket = null;
		}
	});

	if (!socket.init()) {
		(new ErrorDialog())
			.content('Unable to open a websocket.<br />'+
					 'This application requires websockets to work properly.')
			.open();
		return;
	}
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

function onMouseUp (e)
{
	
}

function onMouseDown (e)
{
	
}

function onMouseWheel (e)
{
	
}

String.prototype.trim = function()
{
	return this.replace(/(^\W+)|(\W+$)/, '')
}

