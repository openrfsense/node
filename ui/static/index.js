/**
 * Event handler for form submit which prevents the form from reloading
 * the page but still sends the correct form request.
 */
function formSubmit(event) {
    fetch(event.target.action, {
        method: event.target.method,
        body: new FormData(event.target),
    })
    event.preventDefault()
}

function populateWifiNetworks(event) {
    // if (!event.target.open) return
    event.preventDefault()

    var list = event.target.querySelector("ul")
    console.log(list)

    event.focus()
}

document.addEventListener("DOMContentLoaded", () => {
    document.querySelector("#wifi-form").addEventListener("submit", formSubmit)
    // document.querySelector("#wifi-dropdown").addEventListener("focus", populateWifiNetworks)
})