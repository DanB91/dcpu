// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

if (typeof XMLHttpRequest == "undefined") {
	XMLHttpRequest = function ()
	{
		try {
			return new ActiveXObject("Msxml2.XMLHTTP.6.0");
		} catch (e) { }
		try {
			return new ActiveXObject("Msxml2.XMLHTTP.3.0");
		} catch (e) { }
		try {
			return new ActiveXObject("Microsoft.XMLHTTP");
		} catch (e) { }

		//Microsoft.XMLHTTP points to Msxml2.XMLHTTP and is redundant
		throw new Error("This browser does not support XMLHttpRequest.");
	}
}

function XhrError(code, argv)
{
	this.code = code || 0;
	this.message = ErrorStrings[code];
	this.args = argv || [];
}

var xhr = {
	// request fetches data from a remote host.
	// The supplied object can contain any of these fields:
	//
	// Required fields:
	//  - url:
	//    Target url to fetch.
	// 
	// Optional fields:
	//  - ondata: 
	//    Handler we call when data has been fetched. It gets one parameter
	//    holding the actual data. This can be omitted in synchronous calls.
	//    In which case, request() returns the result directly.
	//  
	//  - method:
	//    GET, POST, HEAD.
	//  
	//  - data:
	//    Data to send to target.
	//  
	//  - onerror:
	//    A handler we call when an error occurred. It gets two parameters
	//    holding the error message and http status code.
	//  
	//  - refresh:
	//    If this is true, we force a new fetch from the server instead of
	//    relying on cached data. Defaults to false.
	//  
	//  - type:
	//    This determines in what format the data parameter comes in the
	//    onData handler. Possible type are 'json' and 'text'.
	//    This defaults to text. 'json' type will attempt to parse the
	//    returned data as a json encoded object and return it.
	//  
	//  - async:
	//    Determines if we should fetch content asynchronously or not.
	//    Defaults to true.
	//  
	//  - headers:
	//    A key/value set of HTTP headers we wish to include.
	request : function (e)
	{
		if (e.method == undefined) {
			e.method = 'GET';
		}

		if(e.method == 'POST') {
			if (e.headers == undefined) {
				e.headers = {};
			}

			e.headers['Content-Type'] = 'application/x-www-form-urlencoded';
		}

		if (e.async == undefined) {
			e.async = true;
		}

		if (e.refresh) {
			var t = new Date();
			e.url += '?' + t.getYear() + t.getMonth() +
				t.getDay() + t.getHours() + t.getMinutes() + t.getSeconds();
		}

		var req = new XMLHttpRequest();

		if (e.async) {
			req.onreadystatechange = function ()
			{
				if (req.readyState == 4) {
					xhr.handleResponse(e, req);
				}
			}
		}

		req.open(e.method, e.url, e.async);

		if (e.headers) {
			for (var key in e.headers) {
				req.setRequestHeader(key, e.headers[key]);
			}
		}

		req.send(e.data);

		if (!e.async) {
			return xhr.handleResponse(e, req);
		}
	},

	handleResponse : function (e, req)
	{
		var d = req.responseText;

		if (e.type == 'json') {
			if (d.length == 0) {
				d = null;
			}

			try {
				eval('d = ' + d);
			} catch (ex) {
				console.error(e.url, ex.toString());
			}
		}

		if (req.status != 200) {
			if (d == null) {
				d = new XhrError(ErrUnknown);
			} else {
				d = new XhrError(d.Code, d.Args);
			}

			if (e.onerror) {
				e.onerror(req.status, d);
			} else {
				throw d;
			}
		} else if (e.ondata) {
			e.ondata(d);
		}

		return d;
	}
};
