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
	// onData:   Handler we call when data has been fetched.
	//           It gets one parameter holding the actual data.
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
	request : function (e)
	{
		if (!e.method) {
			e.method = 'GET';		
		}

		var xhr = new XMLHttpRequest();
		xhr.onreadystatechange = function ()
		{
			if (xhr.readyState != 4) {
				return;
			}

			switch (xhr.status) {
			case 200:
				if (!e.onData) {
					break;
				}
	
				switch (e.type) {
				case 'json':
					var _data = {};
					eval('data = ' + xhr.responseText);
					e.onData(data);
					break;

				default:
					e.onData(xhr.responseText);
					break;
				}

				break;
				
			default:
				if (e.onError) {
					e.onError(xhr.statusText, xhr.status);
				}
			}
		}

		// This is a hack to avoid caching of the requested url.
		// It comes in the form <url>?<timestamp>. The caching
		// mechanism will consider this a unique url and will therefor
		// force a refetch from the server.
		//
		// The server will still refer to the same page and simply
		// ignore the querystring component.
		if (e.refresh) {
			var t = new Date();
			e.url += '?' + t.getYear() + t.getMonth() +
				t.getDay() + t.getHours() + t.getMinutes() + t.getSeconds();
		}

		xhr.open(e.method, e.url, true);
		xhr.send(e.data);
	}
};
