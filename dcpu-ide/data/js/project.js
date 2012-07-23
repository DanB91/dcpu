// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Project represents a single, open DCPU code project.
function Project (name, path, files)
{
	this.name = name || '';
	this.path = path || '';
	this.files = files || [];
	this.hasChanges = true;
}

// create creates a new project by the given name.
function createProject (e)
{
	if (project != null && project.hasChanges) {
		var dlg = new ConfirmDialog({
			yesHandler: function ()
			{
				dlg.close();
				_createProject(e);
			}
		});

		dlg.content('There are unsaved changes to the current project. ' + 
			'Are you sure you want to open a new one? All unsaved ' +
			'progress will be lost.').open();
		return;
	}

	_createProject(e);
}

function _createProject (e)
{
	project = null;

	var query = '';
	for (var k in e) {
		query += k + '=' + encodeURIComponent(e[k]) + '&';
	}

	try {
		var data = api.request({
			url:    '/api/newproject',
			method: 'POST',
			type:   'json',
			async:  false,
			data:   query,
		});

		project = new Project(data.Name, data.Path, data.Files);
	} catch (err) {
		(new ErrorDialog())
			.content('Project creation failed: <br />' + err.msg)
			.open();
	}
}

// load loads project data by name.
Project.prototype.load = function (name)
{
	
}

// save saves project data to name.
Project.prototype.save = function (name)
{
	
}
