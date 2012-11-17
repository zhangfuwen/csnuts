<script type="text/javascript">
    var curr=0;
    var totalArts={{.NumMsgs.Value}}
    var prev=false;
    var next=(curr<totalArts);
    DispNextPrev();
    function DispNextPrev() {
            next=(curr<totalArts)
                prev=(curr>0)
                if(next==false) {
                    $("#Next").css("display","none");
                }else{
                    $("#Next").css("display","block");
                }
            if(prev==true) {
                $("#Prev").css("display","block");
            }else {
                $("#Prev").css("display","none");
            }
    }
    function Next() {
        curr+=10;
        $.ajax({url:"/query?next="+curr,
                success:function(htmlobj) {
                $("#content").html(htmlobj);
                },
                error:function(){
                next=false;
                }
        });
        DispNextPrev();
    }
    function Prev() {
        curr-=10;
        $.ajax({url:"/query?next="+curr,
                success:function(htmlobj) {
                $("#content").html(htmlobj);
                },
                error:function(){
                    yet=false;
                }
        });
        DispNextPrev();
    }
    function Del(msgid) {
        $.ajax({
            url:"/del/?id="+msgid,
            success:function(htmlobj) {
                window.location.reload();
            },
            error:function(htmlobj){
                alert("Bad Request");
            }
        });
    }
    function Good(msgid) {
        $.ajax({url:"/good/?id="+msgid,
            success:function(htmlobj) {
                var good=$("#"+msgid).text();
                good++;
                $("#"+msgid).text(good);
           },
           error:function(htmlobj){
           }
        });
    }
    //输入框自动增高
    function ResizeTextarea() {
        var t=document.getElementById("newctnt");
        var h=t.scrollHeight;
        h=h>110?465:100;
        t.style.height=h+"px";
    }
</script>
