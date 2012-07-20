// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Workspace contains our code editor content.
function Workspace ()
{
	this.node = document.createElement('div');
}

// init initializes the workspace and its UI elements.
Workspace.prototype.init = function ()
{
	this.node.id = 'workspace';
	document.body.appendChild(this.node);
	return true;
}
