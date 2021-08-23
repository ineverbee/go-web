 /*
	Indus by TEMPLATE STOCK
	templatestock.co @templatestock
	Released for free under the Creative Commons Attribution 3.0 license (templated.co/license)
*/

/* ------------------------------------------------------------------------------
 This is jquery module for main page
 ------------------------------------------------------------------------------ */

 /* Global constants */

 /*global jQuery */
 jQuery(function ($) {
  'use strict';



  var App = {
    $options: {},
    $loader: $(".loader"),
    $animationload: $(".animationload"),
    $countdown: $('#countdown_dashboard'),

    bindEvents: function() {
      //binding events
      $(window).on('load', this.load.bind(this));
      $(document).on('ready', this.docReady.bind(this));
    },
    load: function() {
      /*===============================================
      1.Page Preloader
      ===============================================*/
      this.$loader.delay(300).fadeOut();
      this.$animationload.delay(600).fadeOut("slow");
    },
    docReady: function() {
      
      /*===============================================
      NiceScroll
      ===============================================*/
      $("html").niceScroll({
        scrollspeed: 50,
        mousescrollstep: 38,
        cursorwidth: 7,
        cursorborder: 0,
        autohidemode: true,
        zindex: 9999999,
        horizrailenabled: false,
        cursorborderradius: 0
      });

      /*===============================================
      Parallax
      ===============================================*/
      $(window).stellar({
        horizontalScrolling: false,
        responsive: true,
        scrollProperty: 'scroll',
        parallaxElements: false,
        horizontalOffset: 0,
        verticalOffset: 0
      });

      /*===============================================
      When mobile, search popup close, click outside
      ===============================================*/

      $(document).mouseup(function(e) 
      {
          var container = $(".popup");

          // if the target of the click isn't the container nor a descendant of the container
          if (!container.is(e.target) && container.has(e.target).length === 0) 
          {
              $("#search_popup").hide("slow");
              $( "#popup-chat" ).hide( "slow" );
          }
      });

      /*===============================================
      When resize switch between search button and form
      ===============================================*/

      $(window).resize(function() {
          if (window.innerWidth < 560) {

              $('#find-post').replaceWith('<button id="find-post" class="btn post-search-btn"><span><i class="fa fa-search"></i></span></button>');
              
              $( "#find-post" ).click(function() {
                  $( "#search_popup" ).show( "slow" );
              });     

          } else if (window.innerWidth > 560) {

              $('#find-post').replaceWith('<form id="find-post" name="search"><div class="input-group"><input class="form-control" type="text" name="search_text" id="search_text" placeholder="Search" /><span class="input-group-btn"><button class="btn btn-custom" id="" name="search_btn"><span><i class="fa fa-search"></i></span></button></span></div></form>');
              
              $( "#search_popup" ).hide();
          }
      }).resize();

      /*===============================================
      Open-close adding post section and chat
      ===============================================*/

      $( "#add-post-btn" ).click(function() {
          $( "#add-post" ).show( "slow" );
          $("html, body").animate({
              scrollTop: $("#add-post").offset().top
          }, 1000);
      });

      $( "#close-post-btn" ).click(function() {
          $( "#add-post" ).hide( "slow" );
      });


      if (window.innerWidth < 560) {

          $( "#chat-open-btn" ).click(function() {
              $( "#popup-chat" ).show( "slow" );
          });
          $('#popup-chat').replaceWith('<section id="popup-chat" class="b-popup"><div id="chat" class="b-popup-content popup"><div class="row">'
            +'<div id="chat-messages" class="card-content"></div></div>'
            +'<form id="msgForm"><div class="row input-group"><input type="text" id="chat-message" class="form-control chat-input" placeholder="Type msg" style="width: auto; border-radius: 25px; margin-left: 12px;">' 
            +'<button class="btn btn-custom chat-button" id="chat-btn"><i class="fa fa-paper-plane"></i></button></div></form></div></section>');
          $('#chat').replaceWith('<div id="chat"></div>');    

      } else if (window.innerWidth > 560) {

          $( "#chat-open-btn" ).click(function() {
              if ($("#chat").is(":visible")) {
                $( "#chat" ).hide( "slow" );
              } else {
                $( "#chat" ).show( "slow" );
              }
          });
          $('#popup-chat').replaceWith('<span id="popup-chat"></span>');
          $('#chat').replaceWith('<div id="chat" class="d-flex flex-column">'
            +'<div id="chat-messages" class="card-content"></div>'
            +'<form id="msgForm"><div class="row input-group"><input type="text" id="chat-message" class="form-control chat-input" placeholder="Type msg" style="width: auto; border-radius: 25px; margin-left: 12px;">' 
            +'<button class="btn btn-custom chat-button" id="chat-btn"><i class="fa fa-paper-plane"></i></button></div></form></div>');
          
      }
      

      /*===============================================
      Adding file in add-post section
      ===============================================*/

      (function() {
         
        'use strict';
       
        $('.input-file').each(function() {
          var $input = $(this),
              $label = $input.next('.js-labelFile'),
              labelVal = $label.html();
           
         $input.on('change', function(element) {
            var fileName = '';
            if (element.target.value) fileName = element.target.value.split('\\').pop();
            fileName ? $label.addClass('has-file').find('.js-fileName').html(fileName) : $label.removeClass('has-file').html(labelVal);
         });
        });
       
      })();

      /*===============================================
      Reversed order of posts
      ===============================================*/

      $.fn.reverseChildren = function() {
        return this.each(function(){
          var $this = $(this);
          $this.children().each(function(){ $this.prepend(this) });
        });
      };
      $('.reversed').reverseChildren();

      /*===============================================
      Show-hide search suggestions list
      ===============================================*/

      var my_timer;
      $("#search_text").on("focus",
          function () {
              clearTimeout(my_timer);
              $("#search_list").show();
          }                      
      );

      $("#search_text").on("focusout",
          function () {
              var $this = $("#search_list");
              my_timer = setTimeout(function () {
                  $this.hide();
              }, 200);
          } 
      );

      $("#search_text").on("keyup", function() {
      var  value = $(this).val().toLowerCase();
      $("#search_list li").filter(function() {
         $(this).toggle($(this).text().toLowerCase().indexOf(value) > -1)
      });
      });

      $("#mobile_search_text").on("focus",
          function () {
              clearTimeout(my_timer);
              $("#mobile_search_list").show();
          }                      
      );

      $("#mobile_search_text").on("focusout",
          function () {
              var $this = $("#mobile_search_list");
              my_timer = setTimeout(function () {
                  $this.hide();
                  $( "#search_popup" ).hide( "slow" );
              }, 200);
          } 
      );

      $("#mobile_search_text").on("keyup", function() {
      var  value = $(this).val().toLowerCase();
      $("#mobile_search_list li").filter(function() {
         $(this).toggle($(this).text().toLowerCase().indexOf(value) > -1)
      });
      });

      /*===============================================
      Chat app
      ===============================================*/

      var ws = new WebSocket('ws://' + window.location.host + '/ws');
      ws.addEventListener('message', function(e) {
          var msg = JSON.parse(e.data);
          $("#chat-messages").append('<div><div class="chip"><strong>' + msg.username + '</strong><br/>' + msg.message + '</div></div>');
          var element = document.getElementById('chat-messages');
          element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
      });

      $("#msgForm").on("submit", function(e){
          e.preventDefault();
          var form = $(this);
          var val  = $(this).find("input[type=text]").val();
          if(val != ""){
              ws.send(
                JSON.stringify({
                      username: $("#menu-btn").val(),
                      message: val // Strip out html
                  }
              ));
              form[0].reset();
          }
      });

      /*===============================================
      Hiding everything that needs to be hide
      ===============================================*/

      $(function() {
          $( "#add-post" ).hide();
          $( "#search_list" ).hide();
          $( "#mobile_search_list" ).hide();
          $( "#search_popup" ).hide();
          $( "#chat" ).hide()
          $( "#popup-chat" ).hide()
      });

    },
    init: function (_options) {
      $.extend(this.$options, _options);
      this.bindEvents();
    }
  }

  //Initializing the app
  //passing the countdown date
  App.init({});
});