function ()
{
	var f = new Form({
		id: 'frmNewProject',
		handler: createProject,
	});

	f.add({
		type:     'text',
		label:    'Name',
		id:       'tName',
		validate: function ()
		{
			return (this.value.trim().length > 0);
		},
	});
}

