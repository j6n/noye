var autojoin = []
var append = function(data) { autojoin = _.union(autojoin, JSON.parse(data)) }

share.init("channels", append);
share.init("autojoin", append);

respond("!autojoin (?P<method>add|remove) (?P<chan>#.*?)$", function(msg, res) {
  if (!_.contains(noye.auth, msg.From.Nick)) {
    msg.Reply("you're not allowed to do that")
    return
  }

  if (res.method && res.chan) {
    if (res.method == "add" && !_.contains(autojoin, res.chan)) {
      autojoin.push(res.chan)
      msg.Reply("added '%s' to my autojoin", res.chan)
      noye.bot.Join(res.chan)
    }
    if (res.method == "remove" && _.contains(autojoin, res.chan)) {
      var i = autojoin.indexOf(res.chan)
      autojoin.splice(i, 1)
      msg.Reply("removed '%s' from my autojoin", res.chan)
      noye.bot.Part(res.chan)
    }
    return
  }

  msg.Reply("usage: !autojoin add|del #channel")
})

listen("001", function(msg) {
  log("got a 001: %s", autojoin)
  for (var i in autojoin) {
    noye.bot.Join(autojoin[i])
  }
})

cleanup(function(){
  core.save("autojoin", JSON.stringify(autojoin))
})