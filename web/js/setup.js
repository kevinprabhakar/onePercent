function blurElement(element, size){
    var filterVal = 'blur('+size+'px)';
    $(element)
      .css('filter',filterVal)
      .css('webkitFilter',filterVal)
      .css('mozFilter',filterVal)
      .css('oFilter',filterVal)
      .css('msFilter',filterVal);
}

function AddPartnersText(){
    var s = `<div class="row">
                 <div class="col-md-4 pull-left">
                     <input type="email" class="form-control" placeholder="kevin.surya@gmail.com" name="checkerAddr">
                 </div>
                 <div class="col-md-4 pull-left">
                     <input type="text" class="form-control" placeholder="Kevin" name="checkerName">
                 </div>
             </div>`
    return s

}

function nextTab(elem) {
    $(elem).next().find('a[data-toggle="tab"]').click();
}
function prevTab(elem) {
    $(elem).prev().find('a[data-toggle="tab"]').click();
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
            //window.location.replace("http://www.google.com")
        },
        error: function(error){
            console.log(error.responseText)
            window.location.replace("index.html")
        }
    })
}

function addCheckees(accessToken, emailList){
    $.ajax({
        url:"/api/addcheckeeof",
        type: "POST",
        data: {"accessToken":accessToken, "p":JSON.stringify(emailList)},
        success: function(resp){
            console.log(resp)
        },
        error: function(error){
            console.log(error)
        }
    })
}

function addGoal(accessToken, createGoalData, emailListReal){
    $.ajax({
        url:"/api/addgoal",
        type: "POST",
        data: {"accessToken":accessToken, "p":JSON.stringify(createGoalData)},
        success: function(resp){
            console.log(resp)
            addCheckees(accessToken,emailListReal)
        },
        error: function(error){
            console.log(error)
        }
    })
}

function verifyToken(accessToken, goalName, goalDescription, emailList){
    $.ajax({
        url:"/api/verifyaccesstoken",
        type: "POST",
        data: {"accessToken":accessToken},
        success: function(resp){
            var response = JSON.parse(resp)
            var uid = response.userId

            var createGoalData = {
                        "owner" : uid,
                        "name"  : goalName,
                        "description": goalDescription,
                        "created" : moment().unix(),
                        "updateBy" : moment().unix()
                    }

            var emailListReal = {
                "checkerList" : emailList
            }

            addGoal(accessToken, createGoalData, emailListReal)


        },
        error: function(error){
            console.log(error)
        }
    })
}

$(document).ready(function(){
    if (api.readCookie("accessToken") != null){
        VerifyAccessToken(api.readCookie("accessToken"))
    }

    $("#goalName").keyup(function(){
        var textLength = $("#goalName").val().length
        $("#goalNameCounter").text((100-textLength)+ " characters remaining")
    })

    $("#goalDescription").keyup(function(){
            var textLength = $("#goalDescription").val().length
            $("#goalDescriptionCounter").text((500-textLength)+ " characters remaining")
        })

    $("#addPartnersButton").on('click',function(event){
        event.preventDefault()
        var wrapper = $("#addPartnersDiv");
        $(wrapper).append(AddPartnersText()); //add input box

    })
    $("#goalsForm").on('submit', function(event){
        event.preventDefault()
        var checkerList = []
        var goalName = $("#goalName").val()
        var goalDescription = $("#goalDescription").val()
        var checkerAddrList = document.getElementsByName("checkerAddr")
        var checkerNameList = document.getElementsByName("checkerName")
        for (i = 0; i < checkerAddrList.length; i++){
            var checker = {
                "email":checkerAddrList[i].value,
                "name":checkerNameList[i].value
            }
            checkerList.push(checker)
        }

        console.log(checkerList)
        var accessToken = api.readCookie("accessToken")

	//p : {"owner": <currUser UID as string>, "name":<Name of Goal as string>, "description":<Description of goal as string>,
	//     "created": <current time in unix seconds int64>, "updateBy":<time to update by in seconds>}

	    if ((goalName.length==0)||(goalDescription.length==0)||(checkerAddrList.length==0)||(checkerNameList.length==0)){
	        window.location.replace("index.html")
	    }
        verifyToken(accessToken, goalName, goalDescription, checkerList)
        window.location.replace("app.html")
    })






//Initialize tooltips
    $('.nav-tabs > li a[title]').tooltip();

    //Wizard
    $('a[data-toggle="tab"]').on('show.bs.tab', function (e) {

        var $target = $(e.target);

        if ($target.parent().hasClass('disabled')) {
            return false;
        }
    });

    $(".next-step").click(function (e) {

        var $active = $('.wizard .nav-tabs li.active');
        $active.next().removeClass('disabled');
        nextTab($active);

    });
    $(".prev-step").click(function (e) {

        var $active = $('.wizard .nav-tabs li.active');
        prevTab($active);

    });
})


