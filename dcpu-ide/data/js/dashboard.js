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
			src:   '/dashboard/getting_started.html', 
			data:  '', 
		},
		{
			id:    'diNewProject',
			title: 'New Project',
			src:   '/dashboard/new_project.html', 
			key:   'N',
			data:  '',
		},
		{
			id:    'diOpenProject',
			title: 'Open Project',
			src:   '/dashboard/open_project.html', 
			key:   'O',
			data:  '',
		},
		{
			id:    'diConfig',
			title: 'Configuration',
			src:   '/dashboard/config.html', 
			key:   'C',
			data:  '',
		},
		{
			id:    'diHelp',
			title: 'Help',
			src:   '/dashboard/help.html', 
			key:   'H',
			data:  '',
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
	var items = document.getElementById('dashboardItems');
	if (!items) {
		return false;
	}

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

	// Add menu item buttons and pre-load the menu content,
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

		(function(idx) {
			api.request({
				refresh: true,
				url: me.items[idx].src,
				onData : function (data) {
					me.items[idx].data = data;

					if (idx == 0) {
						// Set the first view as the active one,
						// once it is done loading.
						me.select(0);
					}
				},
				onError : function (msg, status) {
					console.error('Dashboard.init: ',
						me.items[idx].src, status, msg);
				},
			});
		}(n));
	}

	items.appendChild(ul);
	return true;
}

// onKey is called whenever a keypress event occurs.
// The parameter holds the key event data.
Dashboard.prototype.onKey = function (e) {
	var key = (e.which != 0) ? e.which : e.keyCode;

	if (!e.altKey) {
		return;
	}

	switch (key) {
	case 192: // ~
		this.toggle();
		break;
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

	this.overview.innerHTML = this.items[index].data;

	// Clear 'active' state on all buttons.
	for (var n = 0; n < this.items.length; n++) {
		var e = document.getElementById(this.items[n].id);
		e.className = '';
	}

	// Set button for current item to 'active'
	var e = document.getElementById(this.items[index].id);
	e.className = 'active';
	this.selectedItem = n;
}

// toggle shows or hides the dashboard using a sliding animation.
Dashboard.prototype.toggle = function ()
{
	var m = fx.metrics(this.node);
	var hide = m.left == 0;
	var node = this.node;

	fx.show(node)
	  .slideTo({
		node:     node,
		left:      hide ? -m.width : 0,
		duration: 500,
		unit:     'px',
		onFinish: function() {
			if (hide) {
				fx.hide(node);
			}
		},
	});
}
