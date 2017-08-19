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

function validateLoginForm(email, password, passwordVerify){
    var checker = ((validateEmail(email))&&(validatePassword(password, passwordVerify)))
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
        },
        error: function(error){
            console.log(error.responseText)
        }
    })
}

$(document).ready(function(){

  // Add smooth scrolling to all links
  $("#signUpErrMessage").hide()
  $("#signUpErrMessage").text("")

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
        window.location.replace("app.html")
    }else{
        $("#loginErrorMessage").hide()
        $("#loginErrorMessage").text("")

        $("#LoginModal").modal("show")
    }
  })

  $("#signUpButton").on('click', function(event){
    event.preventDefault()
    $("#signUpErrMessage").hide()
    $("#signUpErrMessage").text("")
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
                api.createCookie("accessToken",data,7)
                window.location.replace("app.html")
                },
            error: function(error){
                $("#loginErrorMessage").show()
                var msg = error.responseText
                switch(msg){
                    case "MissingEmailField":
                        $("#loginErrorMessage").html("<h5>"+"Pleae Enter an Email Address"+"</h5>")
                        break
                    case "MissingPasswordField":
                        $("#loginErrorMessage").html("<h5>"+"Please Enter a Password"+"</h5>")
                        break
                    case "NonexistentUser":
                        $("#loginErrorMessage").html("<h5>"+"User Doesn't Exist With This Email"+"</h5>")
                        break
                    case "InvalidPassword":
                        $("#loginErrorMessage").html("<h5>"+"Incorrect Password"+"</h5>")
                        break
                    default:
                        $("#loginErrorMessage").hide()
                        $("#loginErrorMessage").text("")
                        console.log(msg)
                }

            }
        })

  })

  $("#modalSignUpButton").on('click', function(event){


    var hostname = window.location.href

    event.preventDefault()
    api.eraseCookie("accessToken")
    var email = $("#email").val()
    var name = $("#name").val()
    var password = $("#pwd").val()
    var passVerify = $("#verifyPassword").val()

    var signUpData = {
        "email": email,
        "password": password,
        "passwordVerify": passVerify,
        "name" : name
    }

//    if (validateLoginForm(email,password,passVerify)){
    $.ajax({
            url: hostname+"/api/signup",
            type: "POST",
            data: {"p":JSON.stringify(signUpData)},
            success: function(resp){

                //replace this log with a redirect link to proper user interface
            api.createCookie("accessToken",resp,7)
            window.location.replace("setup.html")
                //window.location.replace("http://www.google.com")
            },
            error: function(error){
                $("#signUpErrMessage").show()
                var msg = error.responseText
                switch(msg){
                    case "InvalidEmailAddress":
                        $("#signUpErrMessage").html("<h5>"+"Invalid Email Address"+"</h5>")
                        break
                    case "PasswordTooShort":
                        $("#signUpErrMessage").html("<h5>"+"Password is Too Short"+"</h5>")
                        break
                    case "EmailAlreadyExists":
                        $("#signUpErrMessage").html("<h5>"+"User Already Exists With This Email"+"</h5>")
                        break
                    case "PasswordsDontMatch":
                        $("#signUpErrMessage").html("<h5>"+"Provided Password does not match Verify Password"+"</h5>")
                        break
                    default:
                        $("#signUpErrMessage").hide()
                        $("#signUpErrMessage").html("")
                        console.log(msg)


                }

            }
        })
//    }else{
//        console.log("Error: One or more forms invalid")
//    }

  })
});