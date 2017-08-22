function PopulatePage(){
    $.ajax({
        url: "/api/dashboard",
        type: "POST",
        data: {"accessToken":api.readCookie("accessToken")},
        success: function(resp){
            var response = JSON.parse(resp)

            if (response == undefined){
                window.location.replace("setup.html")
                return
            }

            var goal = response[0].goal
            var goalName = goal.name
            var goalDescription = goal.description

            $("#goalName").text(goalName)
            $("#goalDescription").text(goalDescription)

            if (response[0].posts != null){
                var postList = response[0].posts

                for (i = postList.length-1; i>=0; i--){
                    var s =
                    `<tr>
                        <td>${moment(postList[i].created).format("MM/DD/YYYY hh:mm A")}</td>
                        <td>${postList[i].action}</td>
                        <td>${postList[i].feeling}</td>
                        <td>${postList[i].learning}</td>
                    </tr>`
                    $("#historyTableBody").append(s)
                }

                if (moment(postList[postList.length-1].Created).format("MM/DD/YYYY")==moment().format("MM/DD/YYYY")){
                    $('a[href="#onePercentLink"]').hide();
                }
            }
            $('#forDataTable').DataTable( {
                 "order": [[ 0, "desc" ]]
             } );
        },
        error: function(error){
            console.log(error.responseText)
        }
    })
}

function GetStreak(){
    $.ajax({
        url:"/api/goal",
        type:"POST",
        data:{"accessToken":api.readCookie("accessToken")},
        success:function(resp){
            var response = JSON.parse(resp)
            var goalId = response[0].id
            $.ajax({
                url:"/api/getwinstreak",
                type:"POST",
                data:{"accessToken":api.readCookie("accessToken"),"goalId":goalId},
                success:function(resp){
                    var response = JSON.parse(resp)
                    $("#longW").text("Longest Streak: " + resp.longW)
                    $("#currW").text("Current Streak: " + resp.currW)
                },
                error: function(error){

                }
            })
        },
        error: function(error){
            console.log(error)
        }
    })
}

$(document).ready(function(){
    //THIS FUNCTION MUST BE EDITED FOR MULTIPLE GOALS
    PopulatePage()


})