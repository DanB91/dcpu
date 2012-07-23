// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Project represents a single, open DCPU code project.
function Project ()
{
	this.name = '';
	this.path = '';
	this.files = [];
}

// create creates a new project by the given name.
function createProject (e)
{
	console.log('new project: ', e);

	/*
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
	*/
}

// load loads project data by name.
Project.prototype.load = function (name)
{
	
}

// save saves project data to name.
Project.prototype.save = function (name)
{
	
}
