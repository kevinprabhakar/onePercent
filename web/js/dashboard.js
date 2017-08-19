function PopulatePage(){
    $.ajax({
        url: "/api/dashboard",
        type: "POST",
        data: {"accessToken":api.readCookie("accessToken")},
        success: function(resp){
            console.log(resp)
            var response = JSON.parse(resp)

            console.log(response)

            if (response == undefined){
                window.location.replace("setup.html")
                return
            }

            var goal = response[0].goal
            var goalName = goal.Name
            var goalDescription = goal.Description

            $("#goalName").text(goalName)
            $("#goalDescription").text(goalDescription)

            if (response[0].posts != null){
                var postList = response[0].posts

                for (i = postList.length-1; i>=0; i--){
                    var s =
                    `<tr>
                        <td>${moment(postList[i].Created).format("MM/DD/YYYY hh:mm A")}</td>
                        <td>${postList[i].Action}</td>
                        <td>${postList[i].Feeling}</td>
                        <td>${postList[i].Learning}</td>
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

$(document).ready(function(){
    //THIS FUNCTION MUST BE EDITED FOR MULTIPLE GOALS
    PopulatePage()


})