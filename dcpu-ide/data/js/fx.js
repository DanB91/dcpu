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
		return {top: t, left: l, bottom: t+h, right: l+w, width: w, height: h};
	},

	// slideTo moves the element to a specific location using an animation.
	// It accepts a number of options as part of the input object:
	//
	// - node:     The target node.
	// - top:      (optional) top location to move to.
	// - left:     (optional) left location to move to.
	// - duration: (optional) Number of milliseconds the animation should take.
	// - unit:     (optional) The coordinate unit. Defaults to 'px'.
	// - onFinish: (optional) An event handler which is fired when the
	//                        animation is done.
	slideTo : function (cfg)
	{
		if (!cfg.node || cfg.node._fx_busy) {
			return;
		}

		// Ensure we don't get stuck in more than one animation.
		cfg.node._fx_busy = true;
		
		var m = fx.metrics(cfg.node);
		if (cfg.left == m.left && cfg.top == m.top) {
			return; // Nothing to do.
		}

		if (!cfg.unit) {
			cfg.unit = 'px';
		}

		// Number of steps we can fit in to @duration with
		// the established framerate.
		var steps = Math.ceil((parseInt(cfg.duration) || 100) / FxFrameRate);

		// Distance per step to new point.
		var dx = Math.floor((cfg.left - m.left) / steps);
		var dy = Math.floor((cfg.top - m.top) / steps);

		var style = cfg.node.style;

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

			style.left = (parseInt(style.left) || 0) + dx + cfg.unit;
			style.top = (parseInt(style.top) || 0) + dy + cfg.unit;
			steps--;
		}, FxFrameTime);

		return this;
	},
};
