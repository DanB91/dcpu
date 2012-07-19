// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// A Form is a dynamic version of a simple HTML form.
// It takes care of hooking up proper field validation and
// generally serves to reduce the boilerplate we have
// to deal with when writing forms manually.
//
// Note that this does not create an actual <form> element.
// Just a list of elements we track for content. 
function Form (id, method, submitLabel)
{
	this.method = method || "POST";
	this.list = document.createElement('ul');

	var me = this;
	var submit = document.createElement('button');
	submit.innerHTML = submitLabel || 'Submit';
	submit.onclick = function ()
	{
		me.submit();
	}
	this._add(null, submit, true);

	var node = document.getElementById(id);
	node.className = 'form';
	node.appendChild(this.list);
}

// add takes an object describing the field to add to the form.
// It returns the form itself, so we can chain calls.
Form.prototype.add = function (e)
{
	switch (e.type) {
	case 'text':
		var l = null;

		if (e.label) {
			l = document.createElement('label');
			l.setAttribute('for', e.id);
			l.innerHTML = e.label ? e.label + ":" : '';
		}

		var n = document.createElement('input');
		n.type = 'text';
		n.value = e.value || '';
		n.id = e.id;
		n.name = n.id;

		this._add(l, n, false);
		break;
	}

	this.validate();
	return this;
}

// _add adds the given nodes to the form list as a single entry.
Form.prototype._add = function (label, node, append)
{
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
}

// validate returns true if all form fields meet their requirements.
// It also disables the submit button for as long as this is not the case.
Form.prototype.validate = function ()
{
	return false;
}

// submit submits the form.
Form.prototype.submit = function ()
{
	console.log('submitting ', this);
}
