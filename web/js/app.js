function hideAll(){
    $("#addOnePercent").hide()
    $("#accountPage").hide()
    $("#dashboardPage").hide()
}

function SubmitPost(Id, actionText, feelingText, learningText, uid){

	//p : {"action":<action as string>, "feeling":<feeling as string>, "learning":<learning as string>, "owner":<owner uid as string>,
	//     "goal":<goal id as string>, "created":<int64 unix time of creation>}
    var params = {
        "action":actionText,
        "feeling":feelingText,
        "learning":learningText,
        "owner":uid,
        "goal":Id,
        "created":moment().unix()
    }

    $.ajax({
        url: "/api/addpost",
        type: "POST",
        data: {"accessToken":api.readCookie("accessToken"),"p":JSON.stringify(params)},
        success: function(resp){
            console.log(resp)
        }, error: function(error){
            console.log(error.responseText)
        }
    })
}

function DoGoal(uid){
    $.ajax({
        url: "/api/goal",
        type: "POST",
        data: {"accessToken":api.readCookie("accessToken")},
        success: function(resp){
            var goalId = resp[0].Id
            var actionText = $("#action").val()
            var feelingText = $("#feeling").val()
            var learningText = $("#learning").val()

            SubmitPost(goalId, actionText, feelingText, learningText, uid)
        }, error: function(error){
            console.log(error)
        }
    })
}

function GetGoalAndSubmitPost(){
    $.ajax({
        url:"/api/verifyaccesstoken",
        type: "POST",
        data: {"accessToken":api.readCookie("accessToken")},
        success: function(resp){
            var response = JSON.parse(resp)
            var uid = response.userId

            DoGoal(uid)

        },
        error: function(error){
            console.log(error.responseText)
        }
    })

}

$(document).ready(function(){
    hideAll()
    $("#dashboardPage").show()


    $("#action").keyup(function(){
        var textLength = this.value.match(/\S+/g).length;
        if(textLength > 16){
            var trimmed = $("#action").val().split(/\s+/, 16).join(" ");
            $("#action").val(trimmed + " ")
            $("#actionDescriptionCounter").text((0)+ " words remaining")
        }else{
            $("#actionDescriptionCounter").text((17-textLength)+ " words remaining")
        }
    })
    $("#feeling").keyup(function(){
        var textLength = this.value.match(/\S+/g).length;
        if(textLength > 8){
            var trimmed = $("#feeling").val().split(/\s+/, 8).join(" ");
            $("#feeling").val(trimmed + " ")
            $("#feelingDescriptionCounter").text((0)+ " words remaining")
        }else{
            $("#feelingDescriptionCounter").text((8-textLength)+ " words remaining")
        }
    })
    $("#learning").keyup(function(){
        var textLength = this.value.match(/\S+/g).length;
        if(textLength > 16){
            var trimmed = $("#learning").val().split(/\s+/, 16).join(" ");
            $("#learning").val(trimmed + " ")
            $("#learningDescriptionCounter").text((0)+ " words remaining")
        }else{
            $("#learningDescriptionCounter").text((17-textLength)+ " words remaining")
        }
    })

    $("#submitPostButton").on('click',function(event){
        event.preventDefault
        GetGoalAndSubmitPost()
        window.location.reload()
        hideAll()
        $("#dashboardPage").show();

    })
    $('a[href="#onePercentLink"]').click(function(){
        hideAll()
        $("#addOnePercent").show();
    });
    $('a[href="#dashboardLink"]').click(function(){
        hideAll()
        $("#dashboardPage").show();
    });
    $('a[href="#accountLink"]').click(function(){
        hideAll()
        $("#accountPage").show();
    });
    $('a[href="#logoutLink"]').click(function(){
        api.eraseCookie("accessToken")
        window.location.replace("index.html")
    });
    $('a[href="#homepage"]').click(function(){
        window.location.replace("index.html")
    });

})