// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Available dialog buttons.
const ButtonOk       = 0;
const ButtonCancel   = 1;
const ButtonClose    = 2;
const ButtonPrevious = 3;
const ButtonNext     = 4;
const ButtonYes      = 5;
const ButtonNo       = 6;
const ButtonAbort    = 7;
const ButtonRetry    = 8;
const ButtonIgnore   = 9;
const ButtonForward  = 10;
const ButtonBack     = 11;

// The order of these should match the constants above.
var buttonLabels = [
	"Ok", "Cancel", "Close", "Prev", "Next",
	"Yes", "No", "Abort", "Retry", "Ignore",
	"Forward", "Back"
];

// The order of these should match the constants above.
var buttonTitles = [
	"Ok", "Cancel", "Close", "Previous", "Next",
	"Yes", "No", "Abort", "Retry", "Ignore",
	"Forward", "Back"
];

// Dialog is a simple 'popup window' displaying arbitrary content
// along with some buttons.
//
// In most cases, you do not want to use this directly, but go for some
// subclassed dialog implementation. This class supplies the simplest
// possible imlpementation of a dialog.
//
// Buttons can be specified in arbitrary arrangements. The possible
// button types are listed as constants above. By default, a dialog
// has no buttons at all.
//
// This also means it can not be closed by the user (only programatically).
// The precense of at least one button will allow it to be closed by the user;
// either by clicking a button or by hitting the escape key (in some cases).
function Dialog ()
{
	this.canClose = false;

	// Screen-filling, transparent background.
	// This prevents the user from getting mouse input to any
	// controls below this div. Together with redirection of
	// keyboard input, this makes it a modal dialog.
	this.node = document.createElement('div');
	this.node.className = 'dialog';

	// Actual dialog frame with its contents.
	this.frame = document.createElement('div');
	this.frame.className = 'frame';
	this.frame.style.left = -10000;
	this.node.appendChild(this.frame);

	// Title bar
	this.titlebar = document.createElement('div');
	this.titlebar.className = 'titlebar';
	this.frame.appendChild(this.titlebar);

	// This holds the actual title text.
	var node = document.createElement('div');
	node.className = 'title';
	this.titlebar.appendChild(node);

	// This holds the dialog content.
	this.body = document.createElement('div');
	this.body.className = 'body';
	this.frame.appendChild(this.body);

	// Button containers.
	this.buttonbar = document.createElement('div');
	this.buttonbar.className = 'buttons';
	this.frame.appendChild(this.buttonbar);

	node = document.createElement('div');
	node.className = 'left';
	this.buttonbar.appendChild(node);

	node = document.createElement('div');
	node.className = 'right';
	this.buttonbar.appendChild(node);

	node = document.createElement('div');
	node.className = 'clear';
	this.buttonbar.appendChild(node);

	this.getValue = function () { return null; }
}

// title sets the dialog title.
Dialog.prototype.title = function (title)
{
	var node = this.titlebar.childNodes[0];
	if (node == undefined) {
		return this;
	}

	node.innerHTML = title;
	return this;
}

// content sets the dialog contents.
Dialog.prototype.content = function (data)
{
	this.body.innerHTML = data;
	return this;
}

// button adds a single button definition to the given
// side of the dialog (left or right).
Dialog.prototype.button = function (btn, click, side)
{
	var node = document.createElement('button');

	switch (btn) {
	case ButtonClose:
	case ButtonCancel:
		this.canClose = true;

		var me = this;
		node.onclick = function ()
		{
			me.close();
		}

		break;
	}

	if (node.onclick == null) {
		node.onclick = click;
	}

	node.title = buttonTitles[btn];
	node.innerHTML = buttonLabels[btn];

	var side = side == "left" ? 0 : 1;
	this.buttonbar.childNodes[side].appendChild(node);
}

// buttons sets the dialog buttons.
// The argument is a list of objects with the following fields:
//
// - type
//   An integer determining the type of button.
//   This should be any of the predefined ButtonXXX constants.
//
// - click
//   The onclick handler for the button.
Dialog.prototype.buttons = function (e)
{
	if (e == undefined || e.length == 0) {
		return this;
	}

	for (var n = 0; n < e.length; n++) {
		var side = 'left';

		switch (e[n].type) {
		case ButtonOk:
		case ButtonNext:
		case ButtonForward:
		case ButtonYes:
		case ButtonNo:
			side = 'right';
			break;
		}

		this.button(e[n].type, e[n].click, side);
	}
	
	return this;
}

// onKey handles keyboard input.
Dialog.prototype.onKey = function (e)
{
	var key = (e.which != 0) ? e.which : e.keyCode;

	switch (key) {
	case 27: // escape
		if (this.canClose) {
			this.close();
		}
	}
}

// open displays the dialog.
Dialog.prototype.open = function ()
{
	lockApplication(this);

	document.body.appendChild(this.node);
	fx.show(this.node);

	var node = this.frame;
	fx.show(node);

	var sm = fx.metrics();
	var nm = fx.metrics(node);

	fx.move({
		node: node,
		left: -nm.width,
		top:  (sm.height / 2) - (nm.height / 2),
	}).slideTo({
		node:     node,
		left:     (sm.width / 2) - (nm.width / 2),
		duration: 800,
	});

	return this;
}

// close hides the dialog and removes its elements from the DOM.
Dialog.prototype.close = function ()
{
	var n = this.node;
	var m = fx.metrics();

	fx.slideTo({
		node:     this.frame,
		left:     m.width,
		duration: 800,
		onFinish: function ()
		{
			document.body.removeChild(n);
			dialogs.pop();
		}
	});

	unlockApplication(this);
	return this;
}
