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
    var innerText;
    innerText = $.trim($("#texxt").val());
    if (innerText !== "") {
      $(".messages").append(
            "<li class=\"i\"><div class=\"head\"><span class=\"time\">" +
            (new Date().getHours()) + ":" + (new Date().getMinutes()) +
            ", Today</span><span class=\"name\"> Me</span></div><div class=\"message\">" +
            innerText + "</div></li>"
      );
      claerResizeScroll();
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
      insertI();
      return false;
    });
  });

}).call(this);