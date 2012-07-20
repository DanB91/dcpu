function ()
{
	var f = new Form('frmNewProject', "POST",
		'/api/newproject', 'Create project');
	f.add({
		type:     'text',
		label:    'Name',
		id:       'tName',
		validate: function ()
		{
			return (this.value.length > 0);
		},
	});

	f.onData = function (data)
	{
		
	};

	f.onError = function (status, msg)
	{
		if (msg == null) {
			msg = {Message: 'Unknown error'};
		}
	}
}
