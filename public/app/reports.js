const reportRequest = {
  reports: [
  {
    options: {
      title: "TrialBalance",
      startdate: "2019-07-01",
      enddate: "2020-06-30"
    },
    columns: [
      "AccountName",
      "Amount",
      "Currency"
    ]
 }
]}

function getReport(reportName) {  
  console.log(JSON.stringify(reportRequest))
  try {
    fetch('/api/reports/',{
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=utf-8'
      },
      body: JSON.stringify(reportRequest)
    })
    .then(response => response.json())
    .then(data => {
      console.log(data)
      //Clear the page and create a table
      clearMain();
      createTableFromHeaders(data.columns);
      //Copy the data directly into a table for it
      createTableFromReport(data.result);
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error)
  }

}

function clearMain() {
  $(document.getElementById("maincontainer")).empty();;
}
              
function createTableFromHeaders(columns) {
  var container = document.getElementById('maincontainer');
  var table = document.createElement('table');
  table.classList.add('table');
  table.id = "reportstable"
  var thead = document.createElement('thead');
  var tr = document.createElement('tr');
  for (var i = 0; i < columns.length; i++) {
    var th = document.createElement('th');
    th.appendChild(document.createTextNode(columns[i]))
    tr.appendChild(th)
  }
  thead.appendChild(tr);
  table.appendChild(thead);
  container.appendChild(table);
}

function createTableFromReport(rows) {
  console.log(rows)
  var table = document.getElementById('reportstable');
  var tbdy = document.createElement('tbody');
  for (var i = 0; i < rows.length; i++) {
    var tr = document.createElement('tr');
    // Iterate over all the elements in each row
    for (var j = 0; j < rows[i].row.length; j++) {
      var td = document.createElement('td');
      td.appendChild(document.createTextNode(rows[i].row[j]))
      tr.appendChild(td)
    }
    tbdy.appendChild(tr);
  }
  table.append(tbdy);
}
