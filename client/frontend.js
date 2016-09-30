
/* author : Gaetano Mondelli
*  Each event calls the specifi handler that
*  prepares the JSON and send it as a string 
*  to the server.
* 
* 
*/

      $(document).ready(function(){         
        counter = 0;
        resizeFrom = new Object();  
        resizeFrom.width = $(window).width().toString();
        resizeFrom.height = $(window).height().toString();
       
       
        $('input').on('input',function(e){  
		setInterval(function(){ counter++; /*$("#sample").text(counter)*/}, 1000); //1 second          
            $( this ).off(e); //yolo
            });

        $(window).resize(function(e){
              screenResizeHandler();
              $( this ).off(e); //yolo
        });
		

        $('input').on('copy',function(e){  
            copyAndPasteHandler("false",e)
         });

        $('input').on('paste',function(e){  
               copyAndPasteHandler("true",e)      
         });

        $('form').submit(function(e) {
          e.preventDefault();
          timeTaken();   
         });
           
    });
     
     function copyAndPasteHandler(pasted,e)
     {
       
           eventObject = new Object();
           eventObject.eventType = "copyAndPaste";
           eventObject.websiteUrl = window.location.href;  
           eventObject.sessionId = Cookies.get("sessionId");
           eventObject.pasted = pasted;
           eventObject.formId= e.target.id         
           $.post( "\\",  JSON.stringify(eventObject) );

     }

     function screenResizeHandler()
     {
           eventObject = new Object();
           eventObject.eventType = "screenResize";
           eventObject.websiteUrl = window.location.href;  
           eventObject.sessionId = Cookies.get("sessionId");           
           eventObject.resizeFrom = resizeFrom;   
           resizeTo = new Object();
           resizeTo.width = $(window).width().toString();
           resizeTo.height = $(window).height().toString();
           eventObject.resizeTo = resizeTo;           
           $.post( "\\", JSON.stringify(eventObject) );
           // resizeFrom = resizeTo; ! Assume only one resize occurs
     }

     function timeTaken()
     {
        eventObject = new Object();
        eventObject.eventType = "timeTaken";
        eventObject.websiteUrl = window.location.href;  
        eventObject.sessionId = Cookies.get("sessionId");    
        eventObject.time = counter;//.toString();
        $.post( "\\", JSON.stringify(eventObject) );
     }
