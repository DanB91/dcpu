// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// StateTracker monitors the applications connection to the server.
function StateTracker ()
{
	this.node         = null;
	this.pingInterval = 5000;
	this.isConnected  = true;
}

// init initializes the tracker and its associated components.
StateTracker.prototype.init = function ()
{
	this.node = document.getElementById('statetracker');
	if (!this.node) {
		return false;
	}

	this.node.title = "Not connected to server.";

	this.toggle();
	this.poll();
	return true;
}

// poll issues periodic keep-alive pings and monitors
// if the application is still connected to a server.
StateTracker.prototype.poll = function ()
{
	var me = this;
	var interval = setInterval(function () {
		api.request({
			url: '/api/ping',
			onData : function (data) {
				if (!me.isConnected) {
					me.isConnected = true;
					me.toggle();
				}
			},
			onError : function (msg, status) {
				if (me.isConnected) {
					me.isConnected = false;
					me.toggle();
				}
			},
		});
	}, this.pingInterval);
}

// toggle shows or hides the state tracker using a sliding animation.
StateTracker.prototype.toggle = function ()
{
	var m = fx.metrics(this.node);
	var hide = this.isConnected;
	var node = this.node;

	fx.show(node)
	  .slideTo({
		node:     node,
		bottom:   hide ? -m.height : 0,
		duration: 1000,
		unit:     'px',
		onFinish: function() {
			if (hide) {
				fx.hide(node);
			}
		},
	});
}
