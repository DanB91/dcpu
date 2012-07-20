// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

// Color can parse valid CSS colour values (except hsl).
// Its four fields (r, g, b, a) hold the colour components
// as 8 bit integers (0-255).
//
function Color(value)
{
	this.r = 0;
	this.g = 0;
	this.b = 0;
	this.a = 0;

	this.parse(value);
}

// toRgba returns the colour as a rgba(....) string.
Color.prototype.toRgba = function ()
{
	var a = parseFloat(this.a/100)/2.55; // Alpha comes as a float.
	return 'rgba('+this.r+','+this.g+','+this.b+','+a+')';
}

// toRgb returns the colour as a rgb(....) string.
Color.prototype.toRgb = function ()
{
	return 'rgb('+this.r+','+this.g+','+this.b+')';
}

// parse parses the input string as a colour value.
//
// Instead of doing complicated and error-prone regular expression
// parsing on the input string, we are simply going to create a 1x1 pixel
// canvas. Render to it using the input colour and see whatever it created.
//
// What's that? I'm being a lazy ass? You bet ya!
Color.prototype.parse = function(value)
{
	if (!value || value.length == 0) {
		return;
	}

	var w = 12, h = 12;
	var canvas = document.createElement('canvas');
	canvas.width = w;
	canvas.height = h;
	canvas.style.width = w + 'px';
	canvas.style.height = h + 'px';
	document.body.appendChild(canvas);

	var ctx = canvas.getContext('2d');
	ctx.fillStyle = value;
	ctx.fillRect(0, 0, w, h);

	var pixels = ctx.getImageData(0, 0, w, h).data;
	this.r = pixels[0];
	this.g = pixels[1];
	this.b = pixels[2];
	this.a = pixels[3];

	// Remove canvas.
	document.body.removeChild(canvas);
}
