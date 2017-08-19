function GetUser(accessToken){
    $.ajax({
        url:"/api/verifyaccesstoken",
        type: "POST",
        data: {"accessToken":accessToken},
        success: function(resp){
            var response = JSON.parse(resp)
            var uid = response.userId

            var uidList = []
            uidList.push(uid)

            var uidListReal = {
                "idList" : uidList
            }

            $.ajax({
                url:"/api/user",
                type: "POST",
                data: {"accessToken":accessToken,"p":JSON.stringify(uidListReal)},
                success: function(resp){
                    var response = JSON.parse(resp)

                    var userName = response[0].Name
                    var userEmail = response[0].Email
                    var userCheckers = response[0].CheckeeOf


                    $("#currUserName").text(userName)
                    $("#changeEmailButton").before(userEmail + "<br>")
                    for (i = 0; i<userCheckers.length;i++){
                        $("#addPartnerButton").before("<span name='email_" + i.toString()+"'>"+
                                                        userCheckers[i].email+"</span>"+"<span name='name_" + i.toString()+"'>"+" ("+
                                                        userCheckers[i].name+")"+"</span>"+
                                                        "<button class='btn btn-danger' id='email_"+i.toString()+"'>Remove</button>"+"<br>")

                    }
                    $('[id^="email"]').click(function(){
                            var emailNum = $(this).attr('id')
                            var num = emailNum.split("_")[1]
                            var removeEmail = $('[name=' + emailNum + ']').text()
                            var accessToken = api.readCookie("accessToken")

                            if (accessToken == null){
                                window.location.replace('index.html')
                            }

                            $.ajax({
                                    url:"/api/removecheckeeof",
                                    type: "POST",
                                    data: {"accessToken":accessToken, "checkee":removeEmail},
                                    success: function(resp){
                                        $('[name=' + emailNum + ']').css('color', '#a9bad6');
                                        $('[name=' + 'name_' + num + ']').css('color', '#a9bad6');
                                        $("#"+emailNum).remove()
                                        console.log(resp)
                                    },
                                    error: function(error){
                                        console.log(error.responseText)
                                    }
                                })
                        })

                },
                error: function(error){
                    console.log(error.responseText)
                }
            })
        },
        error: function(error){
            console.log(error.responseText)
        }
    })
}

function AddCheckees(){
    var addPartnersList = []

    var emailList = document.getElementsByName('addPartnersEmailForm')
    var nameList = document.getElementsByName('addPartnersNameForm')

    if (emailList.length > 0){
        for (i=0;i<emailList.length;i++){
            if ((emailList[i].value.length == 0)||(nameList[i].value.length == 0)){
                //return error
                return
            }

            var checker = {
                "email": emailList[i].value,
                "name": nameList[i].value
            }
            addPartnersList.push(checker)
        }

        var realAddList = {
            "checkerList" : addPartnersList
        }

        accessToken = api.readCookie("accessToken")

        $.ajax({
                url:"/api/addcheckeeof",
                type: "POST",
                data: {"accessToken":accessToken, "p":JSON.stringify(realAddList)},
                success: function(resp){
                    console.log(resp)
                },
                error: function(error){
                    console.log(error.responseText)
                }
            })
    }
}

function ChangeEmail(){
    var newEmail = document.getElementsByName('changeEmailForm')
    if (newEmail.length != 0){
        var accessToken = api.readCookie("accessToken")
        $.ajax({
            url:"/api/changeemail",
            type: "POST",
            data: {"accessToken":accessToken, "email":newEmail[0].value},
            success: function(resp){
                console.log(resp)
            },
            error: function(error){
                console.log(error.responseText)
            }
        })
    }
}

function ChangePassword(){
    var oldPass = $("#oldPassword").val()
    var newPass = $("#newPassword").val()
    var newPassVerify = $("#newPasswordVerify").val()

    if (newPass != newPassVerify){
        console.log("New Pass and Verify Pass don't match")
        //return error
        return
    }

    console.log(oldPass)

    var params = {
        "oldPassword" : oldPass,
        "newPassword" : newPass
    }

    $.ajax({
        url:"/api/changepassword",
        type: "POST",
        data: {"accessToken":api.readCookie("accessToken"), "p":JSON.stringify(params)},
        success: function(resp){
            if (resp=='{"success":1}'){
                location.reload()
            }
        },
        error: function(error){
            console.log(error.responseText)
        }
    })
}

function DeleteAccount(){
    $.ajax({
            url:"/api/deleteaccount",
            type: "POST",
            data: {"accessToken":api.readCookie("accessToken")},
            success: function(resp){
                api.eraseCookie("accessToken")
                window.location.replace("index.html")
            },
            error: function(error){
                console.log(error.responseText)
            }
        })
}

$(document).ready(function(){
    if (api.readCookie("accessToken") != null){
        GetUser(api.readCookie("accessToken"))
    }else{
        window.location.replace("index.html")
    }

    $("#addPartnerButton").on('click',function(event){
        event.preventDefault()
        console.log("WOO")
        $("#addPartnerButton").before("<input type='text' name='addPartnersEmailForm' placeholder='Email'><input type='text' name='addPartnersNameForm' placeholder='Name'><br>")
    })

    $("#changeEmailButton").on('click',function(event){
        event.preventDefault()
        $("#changeEmailButton").before("<input type='email' name='changeEmailForm'><br>")
        $("#changeEmailButton").hide()
    })

    $("#updateAccountButton").on('click',function(event){
        event.preventDefault()
        AddCheckees()
        ChangeEmail()
        location.reload()
    })

    $("#modalChangePasswordButton").on('click', function(event){
        event.preventDefault()
        ChangePassword()
    })

    $("#modalDeleteAccount").on('click', function(event){
        event.preventDefault()
        DeleteAccount()
    })


})