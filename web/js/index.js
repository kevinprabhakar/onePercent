function validateEmail(email) {
  var re = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
  return re.test(email);
}

function validatePassword(password, passwordVerify){
    if ((password != passwordVerify)||(password.length < 6)){
        return false
    }
    return true
}

function validateLoginForm(email, password){
    var checker = ((validateEmail(email))&&(validatePassword(password)))
    return checker
}

function VerifyAccessToken(accessToken){
    $.ajax({
        url:"/api/verifyaccesstoken",
        type: "POST",
        data: {"accessToken":accessToken},
        success: function(resp){
            var response = JSON.parse(resp)
            var uid = response.userId

            //replace this log with a redirect link to proper user interface
            console.log(uid)
            //window.location.replace("http://www.google.com")
        }
    })
}

$(document).ready(function(){
  // Add smooth scrolling to all links
  $("a").on('click', function(event) {

    if (this.hash !== "") {
      event.preventDefault();
      var hash = this.hash;
      $('html, body').animate({
        scrollTop: $(hash).offset().top
      }, 800, function(){
        window.location.hash = hash;
      });
    }
  });

  // Reset all modal forms upon close
  $('.modal').on('hidden.bs.modal', function(){
      $(this).find('form')[0].reset();
  });

  $("#navLogin").on('click',function(event){
    event.preventDefault()
    if (api.readCookie("accessToken")!=null){
        var token = api.readCookie("accessToken")
        console.log(token)
        VerifyAccessToken(token)
    }
  })


  //Collect/Validate Login Info, Send Request, and create Cookie
  $("#modalLoginButton").on('click', function(event){
        event.preventDefault()

        var emailString = $("#emailLogin").val()
        var passString = $("#pwdLogin").val()

        var signInData = {
        "email" : emailString,
        "password": passString
        }

        $.ajax({
        url:"/api/signin",
        type: "POST",
        data: {"p":JSON.stringify(signInData)},
        success: function(data){
            api.createCookie("accessToken",data,1)
            //Redirect to proper user interface
        }
        })

  })
});