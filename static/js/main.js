let attention = Prompt();
    
// Example starter JavaScript for disabling form submissions if there are invalid fields
(function () {
'use strict'

// Fetch all the forms we want to apply custom Bootstrap validation styles to
let forms = document.querySelectorAll('.needs-validation')

// Loop over them and prevent submission
Array.prototype.slice.call(forms)
    .forEach(function (form) {
        form.addEventListener('submit', function (event) {
            if (!form.checkValidity()) {
            event.preventDefault()
            event.stopPropagation()
            }

            form.classList.add('was-validated')
        }, false)
    })
})()

function notify(msg, msgType) {
    notie.alert({
        type: msgType,
        text: msg,
    })
}

function notifyModal(title, text, icon, confirmButtonText) {
    Swal.fire({
        title: title,
        html: text,
        icon: icon,
        confirmButtonText: confirmButtonText
    })
}

// {{with .Error}}
// notify("{{.}}", "error")
{{end}}
function Prompt() {
    let toast = function(c) {
        const { msg = "", 
                icon="success",
                position = "top-end",

            } = c;
        console.log("clicked button and got toast");
        const Toast = Swal.mixin({
            toast: true,
            title: msg,
            position: position,
            icon: icon,
            showConfirmButton: false,
            timer: 3000,
            timerProgressBar: true,
            didOpen: (toast) => {
                toast.addEventListener('mouseenter', Swal.stopTimer)
                toast.addEventListener('mouseleave', Swal.resumeTimer)
            }
        })

        Toast.fire({})
    }
    
    let success = function(c) {
        const {
            msg = "",
            title = "",
            footer = ""
        } = c;
        console.log("clicked button and got success");
        Swal.fire({
            icon: 'success',
            title: title,
            text: msg,
            footer: footer
        })
    }
    
    let error = function(c) {
        const {
            msg = "",
            title = "",
            footer = ""
        } = c;
        console.log("clicked button and got success");
        Swal.fire({
            icon: 'error',
            title: title,
            text: msg,
            footer: footer
        })
    }
    
    async function custom(c) {
        const {
            msg = "",
            title = ""
        } = c

        const { value: result } = await Swal.fire({
            title: title,
            html: msg,
            backdrop: false,
            focusConfirm: false,
            showCancelButton: true,
            confirmButtonColor:"#ffc107",

            willOpen: () => {
                if(c.willOpen !== undefined) {
                    c.willOpen();
                }
            },
            preConfirm: () => {
               if(c.preConfirm !== undefined) {
                   c.preConfirm();
               }
            },
            didOpen: () => {
               if(c.didOpen !== undefined){
                   c.didOpen();
               }
            },
        })

        if(result) {
            if(result.dismiss !== Swal.DismissReason.cancel){
                if(result.value !== ""){
                    if(c.callback !== undefined) {
                        c.callback(result)
                    }
                } else {
                    c.callback(false)
                }
            } else {
                c.callback(false)
            }
        }
    }

    return {
        toast: toast,
        success: success,
        error:error,
        custom: custom
    }
}