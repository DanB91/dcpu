// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

function Dashboard ()
{
	this.node     = null;
	this.overview = null;
	this.items    = [
		{
			id:    'diGettingStarted',
			title: 'Getting started',
			src:    '/dashboard/getting_started.html', 
		},
		{
			id:    'diNewProject',
			title: 'New Project',
			src:   '/dashboard/new_project.html', 
			key:   'N',
		},
		{
			id:    'diConfigureKeyboard',
			title: 'Configure keyboard',
			src:   '/dashboard/configure_keyboard.html', 
			key:   'K',
		},
		{
			id:    'diHelp',
			title: 'Help',
			src:   '/dashboard/help.html', 
			key:   'H',
		},
	];
	this.selectedItem = -1;
}

// init initializes the dashboard and its UI elements.
Dashboard.prototype.init = function (id)
{
	this.node = document.getElementById('dashboard');
	if (!this.node) {
		return false;
	}

	fx.show(this.node);

	this.overview = document.getElementById('dashboardOverview');
	if (!this.overview) {
		return false;
	}

	// Create list for item buttons.
	var ul = document.createElement('ul');
	if (!ul) {
		return false;
	}

	// Title of dashboard is in first list element.
	var li = document.createElement('li');
	var h3 = document.createElement('h3');
	h3.innerHTML = AppTitle;
	li.appendChild(h3);
	ul.appendChild(li);

	// Add menu item buttons.
	var me = this;
	for (var n = 0; n < this.items.length; n++) {
		var li = document.createElement('li');
		var btn = document.createElement('button');

		btn.id = this.items[n].id;
		btn.title = this.items[n].title;
		btn.innerHTML = btn.title;

		if (this.items[n].key) {
			btn.title += ' (alt+' + this.items[n].key + ')';
		}

		(function(idx) {
			btn.onclick = function () {
				me.select(idx);
			}
		}(n));

		li.appendChild(btn);
		ul.appendChild(li);
	}

	var items = document.getElementById('dashboardItems');
	if (!items) {
		return false;
	}

	items.appendChild(ul);

	this.select(0);
	return true;
}

// onKey is called whenever a keypress event occurs.
// The parameter holds the key event data.
Dashboard.prototype.onKey = function (e) {
	var key = (e.which != 0) ? e.which : e.keyCode;

	if (!e.altKey) {
		switch (key) {
		case 192: // ~
			this.toggle();
			break;
		}

		return;
	}

	var ch = String.fromCharCode(key);

	for (var n = 0; n < this.items.length; n++) {
		if (!this.items[n].key) {
			continue;
		}

		if (ch != this.items[n].key) {
			continue;
		}

		if (!fx.isVisible(this.node)) {
			this.toggle();
		}

		this.select(n);
		e.cancelBubble = true;
	}
}

// select changes the currently active dashboard item to the given index.
Dashboard.prototype.select = function (index)
{
	if (index < 0 || index >= this.items.length || this.selectedItem == index) {
		return;
	}

	var me = this;
	api.request({
		url: me.items[index].src,
		onData : function (data) {
			me.overview.innerHTML = data;

			for (var n = 0; n < me.items.length; n++) {
				var e = document.getElementById(me.items[n].id);
				e.className = '';
			}

			var e = document.getElementById(me.items[index].id);
			e.className = 'active';
			me.selectedItem = index;
		},
		onError : function (msg, status) {
			console.error('Dashboard.select: ',
				me.items[index].src, status, msg);
		},
	});
}

// toggle shows or hides the dashboard using a sliding animation.
Dashboard.prototype.toggle = function ()
{
	var m = fx.metrics(this.node);
	var hide = m.top == 0;
	var node = this.node;

	fx.show(node)
	  .slideTo({
		node:     node,
		top:      hide ? -m.height : 0,
		duration: 500,
		unit:     'px',
		onFinish: function() {
			if (hide) {
				fx.hide(node);
			}
		},
	});
}
