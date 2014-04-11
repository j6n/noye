var autojoin = []
var quakenet = {}

var save = function() { core.save("autojoin", JSON.stringify(autojoin)) }

share.init("quakenet", function(data) {
  quakenet = JSON.parse(data)
})

share.init("channels", function(data) {
  autojoin = _.union(autojoin, JSON.parse(data))
  save()
})

respond("!autojoin (?P<method>add|remove) (?P<chan>#.*?)$", function(msg, res) {
  if (!_.contains(noye.auth, msg.From.Nick)) {
    msg.Reply("you're not allowed to do that");
    return
  }

  if (res.method && res.chan) {
    if (res.method == "add" && !_.contains(autojoin, res.chan)) {
      autojoin.push(res.method)
      msg.Reply("added '%s' to my autojoin", res.chan)
      noye.bot.Join(res.chan)
    }
    if (res.method == "remove" && _.contains(autojoin, res.chan)) {
      var i = autojoin.indexOf(res.chan)
      autojoin.splice(i, 1)
      msg.Reply("removed '%s' from my autojoin", res.chan)
      noye.bot.Part(res.chan)
    }
    
    save()
    return
  }

  msg.Reply("usage: !autojoin add|del #channel")
})

listen("001", function(msg) {
  if (qnetRegex.test(msg.Source.Nick)) {
    if (quakenet.user && quakenet.pass) {
      noye.bot.Send("PRIVMSG %s :AUTH %s %s", "Q@CServe.quakenet.org", quakenet.user, quakenet.pass)
      noye.bot.Send("MODE %s +x", msg.Args[0])
      core.wait(3)
    }
  }

  var init = core.load("autojoin")
  if (init) autojoin = _.union(autojoin, JSON.parse(init))

  for (var i in autojoin) {
    noye.bot.Join(autojoin[i])
  }
})

cleanup(function() {
  save()
})

var qnetRegex = new RegExp("quakenet\\.org$")