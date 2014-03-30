respond("!hello", function(msg) {
	noye.reply(msg, "hi!");
});

listen("001", function(msg) {
	var channels = ["#museun","#nanashin","#noye"];
	channels.map(noye.bot().Join);	
})
