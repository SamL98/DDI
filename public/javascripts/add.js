function displayAddBox() {
    let form = d3.select('body').append('div')
        .attr('id', 'add-input-group')
        .style('position', 'absolute')
        .style('top', 0).style('left', 0)
        .style('width', '90%').style('height', '40px')
        .style('background-color', 'lightgray')
        .style('background-color', 'lightgray')
        .style('margin', '10px 5% 20px 5%');

    let drugGroup = form.append('div')
        .attr('position', 'grid')
        .style('background-color', 'lightgray');

    drugGroup.append('label')
        .attr('for', 'drug')
        .text('Enter the drug you want to add:');

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
                            console.log(drug, drugInput, drug.toLowerCase().startsWith(drugInput))
                            console.log(matches.indexOf(drug))
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
                            .style('padding-left', '10px')
                            .style('cursor', 'pointer')
                            .text(d => { return d })
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
    var ws = new WebSocket('ws://' + location.host + '/drugs');
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