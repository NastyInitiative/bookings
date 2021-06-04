document.querySelector("#check-availability-button").addEventListener("click", function() {
            
    let html =/* html */ `
        <form action="" id="check-availability-form" method="POST" novalidate class="needs-validation">
            <div class="row">
                <div class="col">
                    <div class="row" id="reservation-date-modal">
                        <div class="col">
                            <input type="text" disabled required class="form-control" name="start" id="start" placeholder="Arrival date">
                        </div>
                        <div class="col">
                            <input type="text" disabled required class="form-control" name="end" id="end" placeholder="Departure date">
                        </div>
                    </div>
                </div>
            </div>
        </form>
    `;

    attention.custom({
        msg: html, 
        title: "Choose your dates",
        willOpen: () => {
            const elem = document.querySelector("#reservation-date-modal");
            const rp = new DateRangePicker(elem, {
                format: "dd/mm/yyyy",
                showOnFocus: true,
            })
        },
        didOpen: () => {
            document.getElementById("start").removeAttribute("disabled");
            document.getElementById("end").removeAttribute("disabled");
        },
        preConfirm: () => {
            return [
                document.getElementById('start').value,
                document.getElementById('end').value
            ]
        },
        callback: function(result) {
            console.log("called");

            let form = document.getElementById("check-availability-form");
            let csrfToken = document.getElementById("csrf_token").value;
            // console.log(csrfToken);
            let formData = new FormData(form);
            formData.append("csrf_token", csrfToken);

            fetch('/search-availability-json', {
                method: "post",
                body: formData,
            })
            .then(response => response.json())
            .then(data => {
                console.log("data ::", data);
            })
        }
    })
})