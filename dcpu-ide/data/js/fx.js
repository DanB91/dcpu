// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

const FxFrameRate       = 50;
const FxFrameTime       = 1000/FxFrameRate;
const FxDefaultDuration = 100;

// This defines some basic visual effects for UI elements..
var fx = {
	// show makes the given element visible.
	show : function (e)
	{
		e.style.visibility = 'visible';
		e.style.display = e.style._fx_display || 'block';
		return this;
	},

	// hide makes the given element invisible.
	hide : function (e)
	{
		e.style.visibility = 'hidden';
		e.style._fx_display = e.style.display;
		e.style.display = 'none';
		return this;
	},

	// isVisible returns true if the given element is currently visible.
	isVisible : function (e)
	{
		return e.style.visibility == 'visible';
	},

	// metrics returns the element's pixel coordinates and dimensions.
	// If e is null, it returns metrics for the current browser window.
	metrics : function (e)
	{
		var t, l, b, r, w, h;

		if (e == undefined) {
			if (window.innerHeight) {
				w = window.innerWidth;
				h = window.innerHeight;
			} else if (document.documentElement && document.documentElement.clientHeight) {
				w = document.documentElement.clientWidth;
				h = document.documentElement.clientHeight;
			} else if (document.body) {
				w = document.body.clientWidth;
				h = document.body.clientHeight;
			}

			r = w;
			b = h;
		} else {
			t = parseInt(e.style.top) ||
				e.clientTop || e.offsetTop || 0;
			l = parseInt(e.style.left) ||
				e.clientLeft || e.offsetLeft || 0;
			w = e.clientWidth || e.offsetWidth|| 0;
			h = e.clientHeight || e.offsetHeight || 0;
			r = parseInt(e.style.right) || l + w;
			b = parseInt(e.style.bottom) || t + h;
		}

		return {top: t, left: l, bottom: b, right: r, width: w, height: h};
	},

	// move puts the given element at the given location.
	// This does not animate. It just moves the object immediately.
	// If you want animation, use fx.slideTo().
	//
	// It accepts a number of options as part of the input object:
	//
	// - node:     The target node.
	// - top:      (optional) top location to move to.
	// - left:     (optional) left location to move to.
	// - bottom:   (optional) bottom location to move to.
	// - right:    (optional) right location to move to.
	// - unit:     (optional) The coordinate unit. Defaults to 'px'.
	move : function (e)
	{
		if (e.node == undefined || e.node._fx_busy) {
			return;
		}

		var n = e.node;
		var m = fx.metrics(n);
	
		n._fx_busy = true;

		if (e.unit == undefined) {
			e.unit = 'px';
		}

		n.style.left = (e.left || m.left) + e.unit;
		n.style.right = (e.right || m.right) + e.unit;
		n.style.top = (e.top || m.top) + e.unit;
		n.style.bottom = (e.bottom || m.bottom) + e.unit;
		n._fx_busy = false;
		return this;
	},

	// fade fades the given element in or out from the current opacity
	// to the destination opacity. It does so incrementally over a 
	// given amount of time.
	//
	// It accepts a number of options as part of the input object:
	//
	// - node:     The target node.
	// - to:       The target opacity in the range 0.0-1.0.
	//             This can be either larger or smaller than the current
	//             node opacity. This effectively fades in or out.
	// - duration: (optional) Number of milliseconds the animation should take.
	// - onFinish: (optional) An event handler which is fired when the
	//                        animation is done.
	fade : function (e)
	{
		if (e.node == undefined || e.to == undefined || e.node._fx_busy) {
			return;
		}

		e.node._fx_busy = true;
		e.to = parseFloat(e.to) || 0.0;

		if (e.duration == undefined) {
			e.duration = DefaultDuration;
		}

		var steps = Math.ceil((parseInt(e.duration) || FxDefaultDuration) / FxFrameRate);
		var style = e.node.style;
		var clr = new Color(style['background-color'] || 'transparent');
		var unit = Math.ceil((parseInt(e.to * 255) - clr.a) / steps);

		var interval = setInterval(function()
		{
			if (steps <= 0) {
				clearInterval(interval);
				e.node._fx_busy = false;

				if (e.onFinish) {
					e.onFinish();
				}

				return;
			}

			clr.a += unit;
			style['background-color'] = clr.toRgba();

			steps--;
		}, FxFrameTime);
	},

	// slideTo moves the element to a specific location using an animation.
	// It accepts a number of options as part of the input object:
	//
	// - node:     The target node.
	// - top:      (optional) top location to move to.
	// - left:     (optional) left location to move to.
	// - bottom:   (optional) bottom location to move to.
	// - right:    (optional) right location to move to.
	// - duration: (optional) Number of milliseconds the animation should take.
	// - unit:     (optional) The coordinate unit. Defaults to 'px'.
	// - onFinish: (optional) An event handler which is fired when the
	//                        animation is done.
	//
	// As far as the target coords go, we require that at least the
	// following are specified: (left or right) and/or (top or bottom).
	slideTo : function (e)
	{
		if (e.node == undefined || e.node._fx_busy) {
			return;
		}

		// Ensure we don't get stuck in more than one animation.
		e.node._fx_busy = true;

		if (e.duration == undefined) {
			e.duration = DefaultDuration;
		}

		if (e.unit == undefined) {
			e.unit = 'px';
		}

		// Number of steps we can fit in to @duration with
		// the established framerate.
		var steps = Math.ceil((parseInt(e.duration) || FxDefaultDuration) / FxFrameRate);
		var style = e.node.style;
		var m = fx.metrics(e.node);
		var dx = 0, dy = 0;
		var sx, sy;

		// Find distance between points on x axis.
		if (e.left != undefined) {
			dx = Math.floor((e.left - m.left) / steps);
			sx = 'left';
		} else if (e.right != undefined) {
			dx = Math.floor((e.right - m.right) / steps);
			sx = 'right';
		}

		// Find distance between points on y axis.
		if (e.top != undefined) {
			dy = Math.floor((e.top - m.top) / steps);
			sy = 'top';
		} else if (e.bottom != undefined) {
			dy = Math.floor((e.bottom - m.bottom) / steps);
			sy = 'bottom';
		}

		// Nothing to do?
		if (sx == undefined && sy == undefined) {
			return;
		}

		// Perform incremental move to new area.
		var interval = setInterval(function()
		{
			if (steps <= 0) {
				clearInterval(interval);
				e.node._fx_busy = false;

				if (e.onFinish != undefined) {
					e.onFinish();
				}

				return;
			}

			if (sx != undefined) {
				style[sx] = (parseInt(style[sx]) || 0) + dx + e.unit;
			}

			if (sy != undefined) {
				style[sy] = (parseInt(style[sy]) || 0) + dy + e.unit;
			}

			steps--;
		}, FxFrameTime);

		return this;
	},
};
