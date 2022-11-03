function populateWifiNetworks(event) {
    // if (!event.target.open) return
    event.preventDefault()

    var list = event.target.querySelector("ul")
    console.log(list)

    event.focus()
}

document.getElementById("wifi-form").addEventListener("submit", event => {
    fetch(event.target.action, {
        method: event.target.method,
        body: new FormData(event.target),
    })
    event.preventDefault()
})
// document.querySelector("#wifi-dropdown").addEventListener("focus", populateWifiNetworks)

document.getElementById("config-form").addEventListener("submit", event => {
    fetch(event.target.action, {
        method: event.target.method,
        body: new FormData(event.target),
    })
    event.preventDefault()
})

// Custom CodeMirror YAML linter, disables the "Save" button on error
CodeMirror.registerHelper("lint", "yaml", text => {
    var found = []
    if (!window.jsyaml) {
      console.error("Error: window.jsyaml not defined, CodeMirror YAML linting cannot run.")
      return found
    }
    try { 
        jsyaml.loadAll(text)
        document.getElementById("config-save").classList.toggle("disabled", false)
    } catch(e) {
        var loc = e.mark,
            from = loc ? CodeMirror.Pos(loc.line, loc.column) : CodeMirror.Pos(0, 0),
            to = from
        found.push({ from: from, to: to, message: e.message })
        document.getElementById("config-save").classList.toggle("disabled", true)
    }
    return found
})

var cm = CodeMirror.fromTextArea(
    document.getElementById("config-textarea"),
    {
        mode: "yaml",
        lineNumbers: true,
        autoCloseBrackets: true,
        styleActiveLine: true,
        lint: true,
        gutters: ['CodeMirror-lint-markers'],
        extraKeys: {
            Tab: function (cm) {
              var spaces = Array(cm.getOption('indentUnit') + 1).join(' ')
              cm.replaceSelection(spaces)
            },
        },
    }
)
cm.on("change", () => cm.save())