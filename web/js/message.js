function getParameterByName(name, url) {
    if (!url) url = window.location.href;
    name = name.replace(/[\[\]]/g, "\\$&");
    var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)"),
        results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, " "));
}

function SendMessage(messageAccessToken, fromEmail, toEmail, messageSubject, messageText){
    var params = {
        "messageAccessToken":messageAccessToken,
        "fromEmail":fromEmail,
        "toEmail":toEmail,
        "messageSubject":messageSubject,
        "messageText": messageText
    }

    $.ajax({
        url:"/api/sendUserEmail",
        type:"POST",
        data:params,
        success: function(resp){
            window.location.replace("index.html")
        },
        error: function(error){
            console.log(error)
        }
    })
}

function AddSubmitFunctionality(messageAccessToken){
    var access = {
        "messageAccessToken" : messageAccessToken
    }

    $.ajax({
        url: "/api/verifymessageaccesstoken",
        type: "POST",
        data: access,
        success: function(resp){
            console.log(resp)
            var response = JSON.parse(resp)
            var fromEmail = response.fromEmail
            var toEmail = response.toEmail

            $("#sendMessageButton").click(function(){
                var messageSubject = $("#messageSubject").val()
                var messageText = $("#messageText").val()
                SendMessage(messageAccessToken, fromEmail,toEmail,messageSubject,messageText)
            })
        },
        err: function(error){
            console.log(error)
        }
    })
}

$(document).ready(function(){
    var link = window.location.href
    var messageAccessToken = getParameterByName("messageAccessToken",link)
    if ((messageAccessToken == '')||(messageAccessToken == null)){
        console.log("No Message Token")
    }else{
        AddSubmitFunctionality(messageAccessToken)
    }
})