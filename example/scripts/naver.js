//http://tvcast.naver.com/v/121116

respond("!naver (?:\\s*|.*Id=)(?P<id>\\d+)(?:&|\\s*$)", function(msg, res) {
  if (res.id) {
    music(res.id, msg)
    return
  }

  msg.Reply("usage: !naver id|url")  
})

var music = function(id, msg) {
  var results = []

  var doc = html.new("http://music.naver.com/promotion/specialContent.nhn?articleId=" + id)
  if (!doc) {
    msg.Reply("error getting '%s' :(", id)
    return
  }

  var ids = doc.Find("li[id^='videoContent']", "id", "ids")
  var popup = doc.Find("div.thumb > a[href='#'][onclick^='nhn.music.openMusicVideoPopup']", "onclick", "popup")
  var iframe = doc.Find("div > p > iframe", "src", "iframe")
  var res = doc.Results()

  _.each(res["ids"], function(id) { results.push({Id: _.last(id.split("_"))}) })
  _.each(res["popup"], function(id) { var id = vidRegex.exec(id)[1]; id ? results.push({Id: match}) : null })
  _.each(res["iframe"], function(id) { results.push({Token: new Token(id)}) })

  _.forEach(results, function(r) {
    if (!_.has(r, "Id")) {
      // we already have a token
      return
    }

    http.get("http://music.naver.com/video/popupPlayerVideo.nhn?videoId=" + r.Id, function(ok, res) {
      var matches = mvpRegex.exec(res)
      var token = { Vid: matches[1], Key: matches[2], Dir: "inKey" }
      var url = urlTmpl.replace("_vid_", token.Vid).replace("_dir_", token.Dir).replace("_key_", token.Key)
      http.get(url, function(ok, res) {
        var j = JSON.parse(res)
        var sorted = _.max(j.videoList, function(vid) { return parseFloat(vid.videoBitrate) })
        http.follow(sorted.playUrl, function(ok, res) {
          http.shorten(res, function(ok, res) {
            msg.Reply("[%s] %s | %s", sorted.encodingOptionName, j.postSubject, res)
          })
        })
      }, {"Referer": "http://serviceapi.rmcnmv.naver.com/flash/shareServicePlayer.nhn", "User-Agent": mobileUA})
    }, {"Referer": "http://music.naver.com"})
  })
}

function Token(url) {
  var res = {}
  _(url.split("&"))
    .map(function(e) { var a = e.split("="); return {'k': a[0], 'v': a[1]} })
    .filter(function(e) { if (e.k) return _.contains(["vid", "outKey", "inKey"], e.k); return false })
    .forEach(function(e) { res[e.k] = e.v });

  this.Key = res["outKey"] || res["inKey"]
  this.Dir = _.has(res, "outKey") ? "outKey" : "inKey"
  this.Vid = res["vid"]
}

var mobileUA = "Mozilla/5.0 (iPhone; U; CPU iPhone OS 4_0 like Mac OS X; en-us) AppleWebKit/532.9 (KHTML, like Gecko) Version/4.0.5 Mobile/8A293 Safari/6531.22.7"
var urlTmpl  = "http://serviceapi.rmcnmv.naver.com/mobile/getMobileVideoInfo.nhn?videoId=_vid_&_dir_=_key_&protocol=http"
var mvpRegex = new RegExp("\"musicVideoPlayer\",\\s*\"(.+?)\",\\s*\"(.+?)\"")
var vidRegex = new RegExp("\\((\\d+),")