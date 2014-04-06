approved = [];

share.init("approved", function(data) {
  approved = _.union(approved, JSON.parse(data));
})

respond("!join (?P<chan>#.*?)$", function(msg, res) {
  if (!_.contains(approved, msg.From.Nick) && !_.contains(noye.auth, msg.From.Nick)) {
    msg.Reply("you're not allowed to do that");
    return
  }

  if (res.chan) {
    noye.bot.Join(res.chan)
    return
  }

  msg.Reply("usage: !join #channel")
});

respond("!part\s*(?:$|(?P<chan>#.*?)$)", function(msg, res) {
  if (!_.contains(approved, msg.From.Nick) && !_.contains(noye.auth, msg.From.Nick)) {
    msg.Reply("you're not allowed to do that");
    return
  }

  if (res.chan) {
    log("parting c '%s'", res.chan)
    noye.bot.Part(res.chan)
  } else {
    log("parting t '%s'", msg.Target)
    noye.bot.Part(msg.Target)
  }
});
