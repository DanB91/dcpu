// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// This is a modal dialog extension specifically meant to display information.
function InfoDialog()
{
	Dialog.call(this);
	this.button(ButtonClose)
	    .title('Info');

	var d = document.createElement('div');
	d.className = 'left icon icon-info';
	d.innerHTML = IconOk;
	this.body.appendChild(d);

	d = document.createElement('div');
	this.body.appendChild(d);

	d = document.createElement('div');
	d.className = 'clear';
	this.body.appendChild(d);
}

InfoDialog.prototype = new Dialog();
InfoDialog.prototype.constructor = InfoDialog;

InfoDialog.prototype.content = function(data)
{
	this.body.childNodes[1].innerHTML = data;
	return this;
}
