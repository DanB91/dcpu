function ()
{
	var f = new Form('frmNewProject', "POST", '/api/newproject', 'Create project');
	f.add({
		type:     'text',
		label:    'Location',
		id:       'tLocation',
		validate: function ()
		{
			return (this.value.length > 0);
		},
	});

	f.onData = function (data)
	{
		
	};
}
