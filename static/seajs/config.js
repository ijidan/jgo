//log patch
if(!window['console']){
	window['console'] = {
		'info': function(){},
		'log': function(){},
		'error': function(){},
		'warn': function(){}
	};
}
seajs.config({
	alias: {
		"jquery": "jquery/jquery-1.8.3.min.js",
		"jquery-1.11.2": "jquery/jquery-1.11.2.min.js",
		"jquerycolor": "jquery/jquerycolor.js",
		"jquery/ui": "jquery/ui/jquery-ui.min.js",
		"jquery/ui/timepicker": "jquery/ui/jquery-ui-timepicker-addon.js",
		"jquery/ui/tooltip": "jquery/ui/jquery-ui-tooltip-addon.js",
		"jquery-lazyload":"jquery/jquery.lazyload.min.js",
		"swiper": "swiper/swiper.min.js",
		"waterfall": "waterfall/waterfall.js",
		"ueditor":"ueditor/ueditor.all.js",
		"ueditor_admin_config":"ueditor/ueditor.admin.js",
		"jqzoom":"jquery/jquery.jqzoom.js"
	},

	paths: {
		"ywj": "ywj/component",
		"ywjui": "ywj/ui",
		"jquery/ui": "jquery/ui/jquery-ui.min.js",
		"jquery/ui/timepicker": "jquery/ui/jquery-ui-timepicker-addon.js",
		"jquery/ui/tooltip": "jquery/ui/jquery-ui-tooltip-addon.js"
	},

	//全局使用jquery
	preload: [
		!window.jQuery ? 'jquery' : ''
	],
	charset: 'utf-8'
});