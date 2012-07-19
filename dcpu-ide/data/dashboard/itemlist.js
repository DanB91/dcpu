[
	{
		id:    'diGettingStarted',
		title: 'Getting started',
		src:   '/dashboard/getting_started.html', 
		data:  '',
	},
	{
		id:    'diNewProject',
		title: 'New Project',
		src:   '/dashboard/new_project.html', 
		key:   'N',
		data:  '',
		init:  function ()
		{
			var f = new Form('frmNewProject', "POST", 'Create project');
			f.add({
				type: 'text',
				label: 'Location',
				value: '',
				id: 'tLocation',
				required: true,
			});
		}
	},
	{
		id:    'diOpenProject',
		title: 'Open Project',
		src:   '/dashboard/open_project.html', 
		key:   'O',
		data:  '',
	},
	{
		id:    'diConfig',
		title: 'Configuration',
		src:   '/dashboard/config.html', 
		key:   'C',
		data:  '',
	},
	{
		id:    'diHelp',
		title: 'Help',
		src:   '/dashboard/help.html', 
		key:   'H',
		data:  '',
	},
	{
		id:    'diAbout',
		title: 'About',
		src:   '/dashboard/about.html', 
		data:  '',
	},
];
