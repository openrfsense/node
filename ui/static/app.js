function formSubmit(event) {
    fetch(event.target.action, {
        method: event.target.method,
        body: new FormData(event.target),
    })
    event.preventDefault()
}

function main() {
    document.querySelector("#wifi-form").addEventListener("submit", formSubmit)
}

document.addEventListener("DOMContentLoaded", main)