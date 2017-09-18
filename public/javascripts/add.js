function displayAddBox(source) {
    let form = d3.select('body').append('div')
        .attr('id', 'add-input-group')
        .style('position', 'absolute')
        .style('top', 0).style('left', 0)
        .style('width', '90%').style('height', '40px')
        .style('background-color', 'lightgray')
        .style('background-color', 'lightgray')
        .style('margin', '10px 5% 20px 5%')

    let drugGroup = form.append('div')
        .attr('position', 'grid')
        .style('background-color', 'lightgray')

    drugGroup.append('label')
        .attr('for', 'drug')
        .style('padding', '5px 5px 0 5px')
        .text('Enter the drug you want to add:')

    drugGroup.append('input')
        .attr('type', 'button')
        .attr('value', '^')
        .style('float', 'right')
        .style('background-color', 'lightgray')
        .style('color', 'black')
        .style('padding', '0 0 0 0')
        .style('border', 'none')
        .style('width', '20px').style('height', '20px')
        .style('margin', '5px 5px 0 0')
        .on('click', e => {
            $('#add-input-group').animate({ height: 0 }, 1000, () => {
                form.remove()
            })
        })

    var matches = []
    drugGroup.append('input')
        .attr('class', 'form-control')
        .attr('id', 'add-input')
        .attr('type', 'drug')
        .style('border-radius', '5px')
        .style('font-family', 'Arial')
        .style('margin', '10px 0 10px 0')
        .on('input', () => {
            let drugInput = document.getElementById('add-input').value.toLowerCase()
            if (drugInput.length >= 2) {
                drugs.forEach(drug => {
                    if (matches.includes(drug)) {
                        if (!drug.toLowerCase().startsWith(drugInput)) {
                            matches.splice(matches.indexOf(drug), 1)
                            d3.select('#' + drug).remove()
                        }
                    } else {
                        if (drug.toLowerCase().startsWith(drugInput)) {
                            matches.push(drug)
                        }
                    }
                })

                if (matches.length > 0) {
                    if (!dropdownPresent) {
                        dropdownPresent = true
                        drugGroup.insert('div', '#drug-add')
                            .attr('id', 'drug-dropdown')
                            .style('overflow-y', 'scroll')
                            .style('border-width', '5px')
                            .style('border-color', 'steelblue')
                            .style('border-radius', '5px')
                            .style('border-style', 'solid')
                    }

                    d3.select('#drug-dropdown').selectAll('p')
                        .data(matches).enter().append('p')
                            .attr('id', d => { return d })
                            .style('font-family', 'Arial')
                            .style('font-size', '1em')
                            .style('text-align', 'middle')
                            .style('background-color', 'white')
                            .style('margin', '0 0 0 0')
                            .style('padding', '10px 3.5px 10px 3.5px')
                            .style('cursor', 'pointer')
                            .text(d => { return d })
                            .on('click', d => {
                                console.log(source)
                                getDrug(d, source.name, drugJSON => {
                                    if (source.children === undefined) {
                                        source.children = [drugJSON]
                                    } else {
                                        source.children.push(drugJSON)
                                    }
                                    update(source)
                                    form.remove()
                                })
                            })
                    return
                }
            }
            
            dropdownPresent = false
            d3.select('#drug-dropdown').remove()
        });

    drugGroup.append('button')
        .attr('id', 'drug-add')
        .style('border-radius', '5px')
        .style('border', 'none')
        .style('font-family', 'Arial')
        .style('float', 'right')
        .text('Add')
        .on('click', () => {
            form.remove()
        });

    if (drugs.length == 0) {
        getDrugs()
    }
}

var drugs = []
var dropdownPresent = false

function getDrugs() {
    var ws = new WebSocket('ws://' + location.host + '/drugs')
    ws.onopen = e => {
        ws.send("Start")
    }
    ws.onmessage = e => {
        if (e.data === "End") {
            ws.close()
            return
        }
        drugs.push(JSON.parse(e.data).Name)
    }
}

function getDrug(drug, base, callback) {
    var url = 'http://' + location.host + '/drug?\added=' + encodeURIComponent(drug)
    base.split(",").forEach(name => {
        url += '&base=' + name
    })

    $.get(url, data => {
        let or = parseFloat(data)
        callback({
            "or": or,
            "name": drug,
            "p": 0.05
        })
    })
}