// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// This is a modal dialog extension specifically meant to display Confirms.
function ConfirmDialog(e)
{
	Dialog.call(this);
	this.button(ButtonNo, e.noHandler)
	    .button(ButtonYes, e.okHandler)
	    .button(ButtonCancel, e.cancelHandler, 'left')
	    .title('Confirm');

	var d = document.createElement('div');
	d.className = 'left icon icon-question';
	d.innerHTML = IconQuestion;
	this.body.appendChild(d);

	d = document.createElement('div');
	this.body.appendChild(d);

	d = document.createElement('div');
	d.className = 'clear';
	this.body.appendChild(d);
}

ConfirmDialog.prototype = new Dialog();
ConfirmDialog.prototype.constructor = ConfirmDialog;

ConfirmDialog.prototype.content = function(data)
{
	this.body.childNodes[1].innerHTML = data;
	return this;
}
