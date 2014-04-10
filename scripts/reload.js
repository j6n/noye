respond("!reload (?P<script>.*?\.js)", function(msg, res) {
  if (!_.contains(noye.auth, msg.From.Nick)) {
    msg.Reply("you're not allowed to do that");
    return
  }

  if (!res.script) {
    msg.Reply("no script provided");
    return
  }

  var scripts = core.scripts();
  if (_.contains(scripts.Names, res.script)) {
    msg.Reply("reloading: %s", res.script);
    var err = core.manager.Reload(res.script);
    if (err) msg.Reply("error reloading '%s': %+v", res.script, err);
  }
})
