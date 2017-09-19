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
            // Only find drug matches if the user input is greater than two characters, otherwise it is too long.
            if (drugInput.length >= 2) {
                drugs.forEach(drug => {
                    if (matches.includes(drug)) {
                        // If the drug was previously in the matches, remove it if it no longer matches the input.
                        if (!drug.toLowerCase().startsWith(drugInput)) {
                            matches.splice(matches.indexOf(drug), 1)
                            d3.select('#' + drug).remove()
                        }
                    } else {
                        // Add the drug if it was not a match but now matches the input (only should occur on second letter).
                        if (drug.toLowerCase().startsWith(drugInput)) {
                            matches.push(drug)
                        }
                    }
                })

                // Only display the dropdown when multiple drugs are matches with the given input.
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
                            .style('max-height', '350px')
                    }

                    d3.select('#drug-dropdown').selectAll('p')
                        .data(matches).enter().append('p')
                            .attr('id', d => { return d })
                            .style('font-family', 'Arial')
                            .style('font-size', '1em')
                            .style('background-color', 'white')
                            .style('margin', '0 0 0 0')
                            .style('padding', '10px 3.5px 10px 3.5px')
                            .style('cursor', 'pointer')
                            .text(d => { return d })
                            .on('click', d => {
                                var backColor = 'white'
                                var frontColor = 'black'

                                // If the drug is selected, its background color is steelblue and its foreground color is white.
                                if (selected.includes(d)) {
                                    selected.splice(selected.indexOf(d), 1)
                                } else {
                                    selected.push(d)
                                    backColor = 'steelblue'
                                    frontColor = 'white'
                                }

                                d3.select('#' + d)
                                    .style('background-color', backColor)
                                    .style('color', frontColor)
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
            getDrug(source.name, drugJSON => {
                // If the source node does not have any children, set it's children to the fetched association.
                if (source.children === undefined) {
                    source.children = [drugJSON]
                } else {
                    source.children.push(drugJSON)
                }
                update(source)
                form.remove()
            })
        });

    // Get the drug names if none are present.
    if (drugs.length == 0) {
        getDrugs()
    }
}

var selected = []
var drugs = []
var dropdownPresent = false

// Get the names of all the drugs in the database through a websocket.
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

// Get the association between the base drugs and the selected drugs to add.
function getDrug(base, callback) {
    var url = 'http://' + location.host + '/drug?stub=stub&'

    selected.forEach(name => {
        url += '&added=' + encodeURIComponent(name)
    })

    base.split(",").forEach(name => {
        url += '&base=' + encodeURIComponent(name)
    })

    $.get(url, data => {
        let or = parseFloat(data)

        // If the OR is 0.0, then this means no association was found.
        if (or == 0.0) {
            alert('There is no existing association between ' + base + ' and ' + selected.join(','))
            selected = []
            d3.select('#add-input-group').remove()
            return
        }

        callback({
            "or": or,
            "name": selected.join(','),
            "p": 0.05
        })
        selected = []
    })
}