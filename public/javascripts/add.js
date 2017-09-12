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

    drugGroup.append('input')
        .attr('class', 'form-control')
        .attr('type', 'drug')
        .style('border-radius', '5px')
        .style('font-family', 'Arial')
        .style('margin', '10px 0 10px 0');

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

    getDrugs();
}

var drugs = []

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
    ws.onclose = e => {
        d3.select('#add-input-group')
            .append('select')
            .selectAll('option').data(drugs)
            .enter().append('option')
            .attr('value', d => { return d })
            .text(d => { return d });
    }
}