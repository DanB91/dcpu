// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

var themes = ["ambiance",  "blackboard", "cobalt", "default", "eclipse", 
	"elegant", "erlang-dark", "lesser-dark", "monokai", "neat", "night", 
	"rubyblue", "vibrant-ink", "xq-dark"];
var editor = null;
var editor_config = {
	mode: "dasm",
	theme: themes[4],
	indentUnit: 3,
	tabSize: 3,
	smartIndent: true,
	indentWithTabs: false,
	electricChars: true,
	autoClearEmptyLines: false,
	lineWrapping: false,
	lineNumbers: true,
	firstLineNumber: 1,
	gutter: true,
	fixedGutter: true,
	matchBrackets: false,
	tabindex: 1,
	value: "; device_index device_detect( device_id )\n" +
";\n" +
"; This finds a specific hardware device index\n" +
"; based on the device ID you specify in registers A and B.\n" +
";\n" +
"; It returns the device index in A.\n" +
"; A will be -1 if the device was not found.\n" +
";\n" +
";\n" +
"; ## Example usage:\n" +
";\n" +
";    ...\n" +
";    set a, 0xf615\n" +
";    set b, 0x7349\n" +
";    jsr device_detect\n" +
";    \n" +
";    ifu a, 0\n" +
";      set pc, device_not_found\n" +
";    ...\n" +
";\n" +
";\n" +
"; ## Version History:\n" +
";   0.1.1: Fix check for negative (signed) loop counter.\n" +
";   0.1.0: Initial implementation for spec 1.7.\n" +
";\n" +
":device_detect\n" +
"   set i, a\n" +
"   set j, b\n" +
"   hwn z\n" +
"   sub z, 1\n" +
"\n" +
":device_detect_loop\n" +
"   hwq z\n" +
"   ife a, i\n" +
"      ife b, j\n" +
"         set pc, device_detect_ret\n" +
"\n" +
"   sub z, 1\n" +
"   ifa z, 0xffff\n" +
"      set pc, device_detect_loop\n" +
"\n" +
":device_detect_ret\n" +
"   set a, z\n" +
"   set pc, pop\n"
};


window.onload = function()
{
	editor = CodeMirror(document.body, editor_config);

	if (!editor) {
		console.error("Failed to find code editor element.");
		return;
	}
}
