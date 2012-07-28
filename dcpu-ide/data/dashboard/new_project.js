function ()
{
	var f = new Form({
		id: 'frmNewProject',
		handler: apiCreateProject,
	});

	f.add({
		type:     'text',
		label:    'Name',
		id:       'name',
		validate: function ()
		{
			return (this.value.trim().length > 0);
		},
	});
}

