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
		console.log(data);
	};

	f.onError = function (status, err)
	{
		(new ErrorDialog())
			.content('Failed to create new project:<br />' +
					  ErrorStrings[err.Code] + '.')
			.open();
	}
}
