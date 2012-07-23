// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// A Form is a dynamic version of a simple HTML form.
// It takes care of hooking up proper field validation and
// generally serves to reduce the boilerplate we have
// to deal with when writing forms manually.
//
// Note that this does not create an actual <form> element.
// Just a list of elements we track for content and changes.
// Sending a form to a remote host, is the responsibility of the
// supplied handler function. it os only called whenever the
// form passes all validation steps on submission.
//
// The parameter is an object with the following fields:
//
// - id: The form's id.
// - handler: A function used to handle the form submission.
//            It accepts one parameter wich is an object with all
//            form field names and values.
function Form (e)
{
	if (e == undefined) {
		return;
	}

	if (e.handler == undefined) {
		console.error('Form has no submission handler.');
		return;
	}

	this.handler = e.handler;
	this.list = document.createElement('ul');
	this.controls = [];

	var me = this;
	var sb = document.createElement('button');

	sb.innerHTML = 'Submit';
	sb.onclick = function ()
	{
		me.submit();
	}

	this._add(null, sb, true);

	var node = document.getElementById(e.id);
	node.className = 'form';
	node.appendChild(this.list);
}

// add takes an object describing the field to add to the form.
// It returns the form itself, so we can chain calls.
Form.prototype.add = function (e)
{
	var me = this, l, n;

	switch (e.type) {
	case 'text':
		n = document.createElement('input');
		n.type = 'text';
		n.value = e.value || '';
		n.addEventListener('keyup', function ()
		{
			me.validate();
		}, false);
		n.addEventListener('blur', function ()
		{
			me.validate();
		}, false);
		break;
		
	default:
		return this;
	}

	if (e.label != undefined) {
		l = document.createElement('label');
		l.setAttribute('for', e.id);
		l.innerHTML = e.label ? e.label + ":" : '';
	}

	n.id = e.id;
	n.validate = e.validate || null;
	n.getValue = e.getValue || function ()
	{
		return this.value;
	};
	n.name = n.id;

	this._add(l, n, false);
	return this;
}

// _add adds the given nodes to the form list as a single entry.
Form.prototype._add = function (label, node, append)
{
	// Add control to our list of tracked nodes.
	this.controls.push(node);

	var li = document.createElement('li');

	if (append) {
		this.list.appendChild(li);
	} else {
		var last = this.list.childNodes[this.list.childNodes.length-1];
		this.list.insertBefore(li, last);
	}

	var div = document.createElement('div');
	div.className = "left";
	li.appendChild(div);

	if (label) {
		div.appendChild(label);
	}

	div = document.createElement('div');
	div.className = "right";
	li.appendChild(div);

	if (node) {
		div.appendChild(node);
	}

	// Final, empty div for layout purposes.
	div = document.createElement('div');
	div.className = "clear";
	li.appendChild(div);

	this.list.childNodes[0].childNodes[1].childNodes[0].focus();
	this.validate();
}

// validate returns true if all form fields meet their requirements.
// It also disables the submit button for as long as this is not the case.
Form.prototype.validate = function ()
{
	if (this.onValidate != undefined && !this.onValidate()) {
		return false;
	}

	for (var n = 1; n < this.controls.length; n++) {
		if (!this.controls[n].validate) {
			continue;
		}

		if (!this.controls[n].validate(this.controls[n])) {
			this.disable();
			return false;
		}
	}

	this.enable();
	return true;
}

// enable enables the submit button; optionally on a timer.
Form.prototype.enable = function (timeout)
{
	var c = this.controls[0];
	timeout = timeout || 0;

	if (timeout == 0) {
		c.removeAttribute('disabled');
		return;
	}

	setTimeout(function() {
		c.removeAttribute('disabled');
	}, timeout);
}

// disable disables the submit button; optionally on a timer.
Form.prototype.disable = function (timeout)
{
	var c = this.controls[0];
	timeout = timeout || 0;

	if (timeout == 0) {
		c.setAttribute('disabled', 'disabled')
		return;
	}

	setTimeout(function() {
		c.setAttribute('disabled', 'disabled')
	}, timeout);
}

// submit submits the form.
Form.prototype.submit = function ()
{
	var data = {};
	for (var n = 1; n < this.controls.length; n++) {
		data[this.controls[n].id] = this.controls[n].getValue();
	}

	this.handler(data);
}
