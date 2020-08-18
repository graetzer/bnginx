$(document).ready(function () {
    $(".btn-delete").on('click', function (event) {
        event.stopPropagation();
        event.stopImmediatePropagation();
        //(... rest of your JS code)

        let url = $(this).attr('data-href');
        if (!url) {
            return;
        }
        let rmv = $(this).attr('data-hide');
        let target = $(this).attr('data-target');

        let xhr = new XMLHttpRequest();
        xhr.onload = function() {
            if (xhr.readyState !== XMLHttpRequest.DONE){
                return;
            }
            let status = xhr.status;
            if (status === 0 || (status >= 200 && status < 400)) {
                if (rmv) {
                    $(rmv).remove();
                    return;
                }
    
                if (target) {
                    window.location = target;
                }
              } else {
                // Oh no! There has been an error with the request!
              }
        };
        xhr.open('DELETE', url);
        xhr.send()
    });
});
