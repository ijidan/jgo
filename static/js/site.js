
$(document).ready(function(){
	//设置ifram的高度
	$("iframe").load(function(){
		var $this=$(this);
		var content=$this.contents();
		debugger;
		var contentHeight=$this[0].document.body.scrollHeight;
		$this.height(contentHeight);

	});
	// $iframe2=$("#iframe2");
	// $iframe2.load(function(){
	// 	var ifm= $iframe2[0];
	// 	var subWeb = document.frames ? document.frames[id].document :ifm.contentDocument;
	// 	if(ifm && subWeb) {
	// 		console.info(ifm, subWeb.body.scrollHeight);
	//
	// 		$iframe2.height(subWeb.body.scrollHeight);
	// 	}
	// });
});