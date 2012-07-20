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

var api = {
	// request fetches data from the backend.
	// The supplied object can contain any of these fields:
	//
	// url:      Target url to fetch.
	// onData:   (optional) Handler we call when data has been fetched.
	//                      It gets one parameter holding the actual data.
	//                      This can be omitted in synchronous calls.
	//                      In which case, request() returns the result directly.
	// method:   (optional) GET, POST, HEAD.
	// data:     (optional) Data to send to target.
	// onError:  (optional) A handler we call when an error occurred.
	//                      It gets two parameters holding the error message
	//                      and http status code.
	// refresh:  (optional) If this is true, we force a new fetch from the
	//                      server instead of relying on cached data. Defaults
	//                      to false.
	// type:     (optional) This determines in what format the data parameter
	//                      comes in the onData handler. Possible type are:
	//                      'json' and 'text'. This defaults to text. 'json'
	//                      type will attempt to parse the returned
	//                      data as a json encoded object and return it.
	// async:    (optional) Determines if we should fetch content asynchronously
	//                      or not. Defaults to true.
	request : function (e)
	{
		if (e.method == undefined) {
			e.method = 'GET';
		}

		if (e.async == undefined) {
			e.async = true;
		}

		if (e.refresh) {
			var t = new Date();
			e.url += '?' + t.getYear() + t.getMonth() +
				t.getDay() + t.getHours() + t.getMinutes() + t.getSeconds();
		}

		var xhr = new XMLHttpRequest();

		if (e.async) {
			xhr.onreadystatechange = function ()
			{
				if (xhr.readyState == 4) {
					api.handleResponse(e, xhr);
				}
			}
		}

		xhr.open(e.method, e.url, e.async);
		xhr.send(e.data);

		if (!e.async) {
			return api.handleResponse(e, xhr);
		}
	},

	handleResponse : function (e, xhr)
	{
		var d = xhr.responseText;
		if (e.type == 'json') {
			if (d.length == 0) {
				d = null;
			}

			eval('d = ' + d);
		}

		if (xhr.status != 200) {
			if (e.onError) {
				e.onError(d, xhr.status);
			}
			return null;
		}

		if (!e.onData) {
			return d;
		}

		e.onData(d);
		return null;
	}
};
