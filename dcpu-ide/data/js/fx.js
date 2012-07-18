// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

const FxFrameRate = 50;
const FxFrameTime = 1000/FxFrameRate;

// This defines some basic visual effects for UI elements..
var fx = {
	// show makes the given element visible.
	show : function (e) {
		e.style.visibility = 'visible';
		e.style.display = 'block';
		return this;
	},

	// hide makes the given element invisible.
	hide : function (e) {
		e.style.visibility = 'hidden';
		e.style.display = 'none';
		return this;
	},

	// isVisible returns true if the given element is currently visible.
	isVisible : function (e) {
		return e.style.visibility == 'visible';
	},

	// metrics returns the element's pixel coordinates and dimensions.
	metrics : function (e) {
		var t = parseInt(e.style.top) || parseInt(e.clientTop) || 0;
		var l = parseInt(e.style.left) || parseInt(e.clientLeft) || 0;
		var w = parseInt(e.clientWidth) || 0;
		var h = parseInt(e.clientHeight) || 0;
		var b = parseInt(e.style.bottom) || t + h;
		var r = parseInt(e.style.right) || l + w;
		return {top: t, left: l, bottom: b, right: r, width: w, height: h};
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
	slideTo : function (cfg)
	{
		if (!cfg.node || cfg.node._fx_busy) {
			return;
		}

		// Ensure we don't get stuck in more than one animation.
		cfg.node._fx_busy = true;

		if (!cfg.unit) {
			cfg.unit = 'px';
		}

		// Number of steps we can fit in to @duration with
		// the established framerate.
		var steps = Math.ceil((parseInt(cfg.duration) || 100) / FxFrameRate);
		var style = cfg.node.style;
		var m = fx.metrics(cfg.node);
		var dx = 0, dy = 0;
		var sx, sy;

		// Find distance between points on x axis.
		if (cfg.left != undefined) {
			dx = Math.floor((cfg.left - m.left) / steps);
			sx = 'left';
		} else if (cfg.right != undefined) {
			dx = Math.floor((cfg.right - m.right) / steps);
			sx = 'right';
		}

		// Find distance between points on y axis.
		if (cfg.top != undefined) {
			dy = Math.floor((cfg.top - m.top) / steps);
			sy = 'top';
		} else if (cfg.bottom != undefined) {
			dy = Math.floor((cfg.bottom - m.bottom) / steps);
			sy = 'bottom';
		}

		// Nothing to do?
		if (!sx && !sy) {
			return;
		}

		// Perform incremental move to new area.
		var interval = setInterval(function() {
			if (steps <= 0) {
				clearInterval(interval);
				cfg.node._fx_busy = false;

				if (cfg.onFinish) {
					cfg.onFinish();
				}

				return;
			}

			if (sx) {
				style[sx] = (parseInt(style[sx]) || 0) + dx + cfg.unit;
			}

			if (sy) {
				style[sy] = (parseInt(style[sy]) || 0) + dy + cfg.unit;
			}

			steps--;
		}, FxFrameTime);

		return this;
	},
};
