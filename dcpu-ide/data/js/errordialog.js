// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// This is a modal dialog extension specifically meant to display errors.
function ErrorDialog()
{
	Dialog.call(this);
	this.button(ButtonClose)
	    .title('Error');

	var d = document.createElement('div');
	d.className = 'left icon error';
	d.innerHTML = IconError;
	this.body.appendChild(d);

	d = document.createElement('div');
	this.body.appendChild(d);

	d = document.createElement('div');
	d.className = 'clear';
	this.body.appendChild(d);
}

ErrorDialog.prototype = new Dialog();
ErrorDialog.prototype.constructor = ErrorDialog;

ErrorDialog.prototype.content = function(data)
{
	this.body.childNodes[1].innerHTML = data;
	return this;
}
