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

// load loads project data by name.
Project.prototype.load = function (name)
{
	
}

// save saves project data to name.
Project.prototype.save = function (name)
{
	
}

// createProject creates a new project
function createProject(e)
{
	if (project != null && project.hasChanges) {
		var dlg = new ConfirmDialog({
			yesHandler: function ()
			{
				dlg.close();
				apiCreateProject(socket, e["name"]);
			}
		});

		dlg.content('There are unsaved changes to the current project. ' + 
			'Are you sure you want to open a new one? All unsaved ' +
			'progress will be lost.').open();
		return;
	}

	apiCreateProject(socket, e["name"]);
}

