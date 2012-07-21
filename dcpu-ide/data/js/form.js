// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// A Form is a dynamic version of a simple HTML form.
// It takes care of hooking up proper field validation and
// generally serves to reduce the boilerplate we have
// to deal with when writing forms manually.
//
// Note that this does not create an actual <form> element.
// Just a list of elements we track for content and changes. 
function Form (id, method, target, submitLabel)
{
	if (!id) {
		console.eror('Form has no element id.');
		return;
	}

	if (!method) {
		console.eror('Form has no method.');
		return;
	}

	if (!target) {
		console.eror('Form has no submission target.');
		return;
	}

	this.target = target;
	this.method = method;
	this.list = document.createElement('ul');
	this.controls = [];
	this.onData = null;
	this.onError = function (status, msg)
	{
		console.error(status, msg.Message);
	}

	var me = this;
	var sb = document.createElement('button');
	
	sb.innerHTML = submitLabel || 'Submit';
	sb.onclick = function ()
	{
		me.submit();
	}

	this._add(null, sb, true);

	var node = document.getElementById(id);
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
		n.addEventListener('change', function ()
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
	this.validate();
}

// validate returns true if all form fields meet their requirements.
// It also disables the submit button for as long as this is not the case.
Form.prototype.validate = function ()
{
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
	this.disable();

	var data = [];
	for (var n = 1; n < this.controls.length; n++) {
		data.push(this.controls[n].id + '='
			+ encodeURIComponent(this.controls[n].getValue()));
	}

	var query = data.join('&');
	var me = this;

	api.request({
		url:    this.target,
		method: this.method,
		type:   'json',
		data:   query,
		headers: {
			'Content-Type': 'application/x-www-form-urlencoded'
		},
		onError: function (status, msg)
		{
			if (me.onError) {
				me.onError(status, msg);
			}
		},
		onData: function (data)
		{
			if (me.onData) {
				me.onData(data);
			}
		},
	});
}
