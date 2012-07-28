// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

if (window.MozWebSocket) {
  window.WebSocket = window.MozWebSocket;
}

// Socket represents a single binary websocket connection.
function Socket (e)
{
	this.conn = null;
	this.onopen = e.onopen || null;
	this.onclose = e.onclose || null;
	this.onerror = e.onerror || null;
	this.address = "{{.SocketAddress}}";
}

// init initializes the connection.
Socket.prototype.init = function (onopen)
{
	if (window.WebSocket === undefined) {
		console.error('window.Websocket not defined.');
		return false;	
	}

	if (this.conn != null && this.conn.readyState < 2) {
		console.error('Connection already open.');
		return false;	
	}

	this.conn = new WebSocket(this.address);
	if (this.conn == null) {
		console.error('new WebSocket() failed.');
		return false;
	}

	var me = this;
	this.conn.binaryType = "arraybuffer";
	
	this.conn.onopen = function ()
	{
		console.log('websocket.onopen');

		if (me.onopen != null) {
			me.onopen();
		}
	};
	
	this.conn.onclose = function (e)
	{
		console.log('websocket.onclose');

		if (me.onclose != null) {
			me.onclose(e);
		}

		me.conn = null;
		me.close();
	};
	
	this.conn.onmessage = function (e)
	{
		if ((e.data instanceof ArrayBuffer) && e.data.byteLength > 0) {
			apiSocketReceive(e.data);
		}
	};
	
	this.conn.onerror = function (e)
	{
		console.error('websocket.onerror', e);
	
		if (me.onerror != null) {
			me.onerror(e);
		}
	};

	return true;
}

// send sends the given message to the server.
Socket.prototype.send = function (data)
{
	if (data === undefined || this.conn == null || this.conn.readyState > 1) {
		return false;
	}
	return this.conn.send(data);
}

Socket.prototype.close = function ()
{
	if (this.conn != null) {
		this.conn.close();
		this.conn = null;
	}

	if (this.onclose != null) {
		this.onclose();
		this.onclose = null;
	}

	this.onopen = null;
	this.onerror = null;

	(new ErrorDialog())
		.content('We have lost our connection to the server. '+
			     'Please ensure it is still running.')
		.open();
}
