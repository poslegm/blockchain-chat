(function() {
  var claerResizeScroll, conf, getRandomInt, insertI, lol;

  conf = {
    cursorcolor: "#696c75",
    cursorwidth: "4px",
    cursorborder: "none"
  };

  lol = {
    cursorcolor: "#cdd2d6",
    cursorwidth: "4px",
    cursorborder: "none"
  };

  getRandomInt = function(min, max) {
    return Math.floor(Math.random() * (max - min + 1)) + min;
  };

  claerResizeScroll = function() {
    $("#texxt").val("");
    $(".messages").getNiceScroll(0).resize();
    return $(".messages").getNiceScroll(0).doScrollTop(999999, 999);
  };

  insertI = function() {
    var innerText, otvet;
    innerText = $.trim($("#texxt").val());
    if (innerText !== "") {
      $(".messages").append(
            "<li class=\"i\"><div class=\"head\"><span class=\"time\">" +
            (new Date().getHours()) + ":" + (new Date().getMinutes()) +
            ", Today</span><span class=\"name\"> Букер</span></div><div class=\"message\">" +
            innerText + "</div></li>"
      );
      claerResizeScroll();

      return otvet = setInterval(function() {
        $(".messages").append(
            "<li class=\"friend-with-a-SVAGina\"><div class=\"head\"><span class=\"name\">VS94SKI  </span><span class=\"time\">" +
            (new Date().getHours()) + ":" + (new Date().getMinutes()) +
            ", Today</span></div><div class=\"message\">" +
            "ЧПЕК" + "</div></li>"
        );
        claerResizeScroll();
        return clearInterval(otvet);
      }, getRandomInt(2500, 500));
    }
  };

  $(document).ready(function() {
    $(".list-friends").niceScroll(conf);
    $(".messages").niceScroll(lol);
    $("#texxt").keypress(function(e) {
      if (e.keyCode === 13) {
        sendMessage();
        insertI();
        return false;
      }
    });
    return $(".send").click(function() {
      sendMessage();
      return insertI();
    });
  });

}).call(this);