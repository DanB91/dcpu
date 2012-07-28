// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// lockApplication locks various nodes to prevent tab input to them,
// This is called by the given modal dialog when it is opened.
function lockApplication(dlg)
{
	var nodelists = [];

	if (dialogs.length == 0) {
		nodelists = nodelists.concat(
			findLockableNodes(workspace.node),
			findLockableNodes(dashboard.node)
		)
	} else {
		nodelists = nodelists.concat(
			findLockableNodes(dialogs[dialogs.length-1].node)
		)
	}

	for (var i = 0; i < nodelists.length; i++) {
		for (var j = 0; j < nodelists[i].length; j++) {
			lock(nodelists[i][j]);
		}
	}

	// Update a list of open dialogs.
	dialogs.push(dlg);
}

// unlockApplication unlocks various nodes to re-enable tab input to them,
// This is called by a modal dialog when it is closed.
function unlockApplication()
{
	// Update a list of open dialogs.
	dialogs.pop();

	var nodelists = [];

	if (dialogs.length == 0) {
		nodelists = nodelists.concat(
			findLockableNodes(workspace.node),
			findLockableNodes(dashboard.node)
		)
	} else {
		nodelists = nodelists.concat(
			findLockableNodes(dialogs[dialogs.length-1].node)
		)
	}

	for (var i = 0; i < nodelists.length; i++) {
		for (var j = 0; j < nodelists[i].length; j++) {
			unlock(nodelists[i][j]);
		}
	}
}

// findLockableNodes finds all children of the given node
// which we consider relevant for locking/unlocking.
function findLockableNodes(e)
{
	return [].concat(
		e.getElementsByTagName('a'),
		e.getElementsByTagName('select'),
		e.getElementsByTagName('input'),
		e.getElementsByTagName('button'),
		e.getElementsByTagName('textarea')
	);
}

// lock locks the given control.
// It is disabled to prevent it from getting key input.
function lock (e)
{
	if (e == undefined || e.attributes['disabled'] != undefined) {
		return; // We don't want to unlock it later when it shouldn't be.
	}

	e.setAttribute('disabled', 'disabled');
	e._locked = true;
}

// unlock unlocks the given control.
// It is (re-)enabled to allow key input.
function unlock (e)
{
	if (e != undefined && e._locked == true) {
		e.removeAttribute('disabled');
		e._locked = false;
	}
}
