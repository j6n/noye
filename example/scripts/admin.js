respond("!part(?:\s*$|(?P<chan>#.*?)$)", function(msg, res) {
  if (!_.contains(noye.auth, msg.From.Nick)) {
    msg.Reply("you're not allowed to do that");
    return
  }

  if (res.chan) {
    noye.bot.Part(res.chan)
  } else {
    noye.bot.Part(msg.Target)
  }
})

respond("!join (?P<chan>#.*?)$", function(msg, res) {
  if (!_.contains(noye.auth, msg.From.Nick)) {
    msg.Reply("you're not allowed to do that");
    return  
  }

  if (res.chan) {
    noye.bot.Join(res.chan)
    return
  }

  msg.Reply("usage: !join #chan")
})