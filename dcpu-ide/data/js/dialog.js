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

// Dialog is a simple 'popup window' displaying a message
// along with some buttons. It is either for informative purposes,
// or to request some sort of confirmation from the user.
function Dialog ()
{
	// Screen-filling, transparent background. This prevents the user from
	// getting mouse input to any controls below this div.
	// Together with redirection of keyboard input,
	// this makes it a modal dialog.
	this.node = document.createElement('div');
	this.node.className = 'dialog';

	// Actual dialog frame with its contents.
	this.frame = document.createElement('div');
	this.frame.className = 'frame';
	this.frame.style.left = -10000;
	this.node.appendChild(this.frame);

	// Title bar
	this.title = document.createElement('div');
	this.title.className = 'titlebar';
	this.frame.appendChild(this.title);

	// This holds the actual title text.
	var node = document.createElement('div');
	node.className = 'title';
	this.title.appendChild(node);

	// This holds the dialog content.
	this.body = document.createElement('div');
	this.body.className = 'body';
	this.frame.appendChild(this.body);

	// Button containers.
	this.buttons = document.createElement('div');
	this.buttons.className = 'buttons';
	this.frame.appendChild(this.buttons);

	node = document.createElement('div');
	node.className = 'left';
	this.buttons.appendChild(node);

	node = document.createElement('div');
	node.className = 'right';
	this.buttons.appendChild(node);

	node = document.createElement('div');
	node.className = 'clear';
	this.buttons.appendChild(node);

	this.getValue = function () { return null; }
}

// setTitle sets the dialog title.
Dialog.prototype.setTitle = function (title)
{
	var node = this.title.childNodes[0];
	if (node == undefined) {
		return this;
	}

	node.innerHTML = title;
	return this;
}

// setContent sets the dialog contents.
Dialog.prototype.setContent = function (data)
{
	this.body.innerHTML = data;
	return this;
}

// setButtons sets the dialog buttons.
// The argument is a list of objects with two fields each:
//
// - type
//   An integer determining the type of button.
//   This should be any of the predefined ButtonXXX values.
//
// - click
//   The onclick handler for the button.
Dialog.prototype.setButtons = function (e)
{
	if (e == undefined || e.length == 0) {
		return this;
	}

	var panels = [
		this.buttons.childNodes[0], // left
		this.buttons.childNodes[1], // right
	];

	for (var n = 0; n < e.length; n++) {
		var side = 0;
		var node = document.createElement('button');
		node.onclick = e[n].click;

		switch (e[n].type) {
		case ButtonOk:
			side = 1;
			node.title = 'Ok';
			break;

		case ButtonCancel:
			node.title = 'Cancel';
			break;

		case ButtonClose:
			node.title = 'Close';
			break;

		case ButtonPrevious:
			node.title = 'Prev';
			break;

		case ButtonNext:
			node.title = 'Next';
			side = 1;
			break;

		case ButtonYes:
			side = 1;
			node.title = 'Yes';
			break;

		case ButtonNo:
			node.title = 'No';
			break;

		case ButtonAbort:
			node.title = 'Abort';
			break;

		case ButtonRetry:
			node.title = 'Retry';
			break;

		case ButtonIgnore:
			node.title = 'Ignore';
			break;

		}

		node.innerHTML = node.title;
		panels[side].appendChild(node);
	}
	
	return this;
}

// onKey handles keyboard input.
Dialog.prototype.onKey = function (e)
{
	var key = (e.which != 0) ? e.which : e.keyCode;

	switch (key) {
	case 27: // escape
		this.close();
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
