// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

const ApiPingInterval = 10000;

// api_sendId is a convenience wrapper for calls which
// send a message id and no parameters.
function apiSendId(s, id)
{
	if (s == null) {
		return false;
	}

	var b = new Uint8Array(1);
	b[0] = id;
	return s.send(b.buffer);
}

// apiSocketReceive is called whenever we receive a message from the server.
function apiSocketReceive (data)
{
	var arr = new Uint8Array(data);
	var id = arr[0];
	arr = arr.subarray(1);

	console.log(id, MessageStrings[id]);

	switch (id) {
	case ErrUnknown:
		break;
	case ErrMissingName:
		break;
	case ErrDuplicateProject:
		break;
	case ErrTemplateFailure:
		break;
	case ErrInvalidPath:
		break;
	case ErrNotDirectory:
		break;
	case ErrNotFile:
		break;
	case ErrPathNotExist:
		break;
	case ErrFileRead:
		break;
	case ApiReadFile:
		break;
	case ApiDirList:
		break;
	case ApiNewProject:
		break;
	}
}

// apiHandshake initiates our connection.
function apiHandshake(s)
{
	apiSendId(s, ApiHello);
	apiPing(s);
}

// Ping sends periodic ping requests to keep the connection alive.
function apiPing(s)
{
	setTimeout(function()
	{
		apiSendId(s, ApiPing);
		apiPing(socket);
	}, ApiPingInterval);
}

// apiCreateProject creates a new 
function apiCreateProject(s, name)
{
	
}

