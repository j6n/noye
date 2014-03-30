respond("!hello", function(msg) {
	noye.reply(msg, "hi!");
});

listen("001", function(msg) {
	var channels = ["#noye"];
	for (var i in channels) {
		noye.bot.Join(channels[i]);
	}
})

respond("!ip", function(msg) {
	var data = core.http("http://ifconfig.me/ip").get();
	noye.reply(msg, data);
})