approved = [];

share.init("approved", function(data) {
  approved = _.union(approved, JSON.parse(data));
})

respond("!join (?P<chan>#.*?)$", function(msg, res) {
  if (!_.contains(approved, msg.From.Nick) || !_.contains(noye.auth, msg.From.Nick)) {
    msg.Reply("you're not allowed to do that");
    return
  }

  if (res.chan) {
    noye.bot.Join(res.chan)
  }

  msg.Reply("usage: !join #channel")
});

respond("!part\s*(?:$|(?P<chan>#.*?))$", function(msg, res) {
  if (!_.contains(approved, msg.From.Nick) || !_.contains(noye.auth, msg.From.Nick)) {
    msg.Reply("you're not allowed to do that");
    return
  }

  if (res.chan) {
    noye.bot.Part(res.chan)
  } else {
    noye.bot.Part(msg.Target)
  }

  msg.Reply("usage: !part <#channel>")
});
