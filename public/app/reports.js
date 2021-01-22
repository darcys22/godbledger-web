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
  //console.log(JSON.stringify(reportRequest))
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
      //Clear the page and create a table
      clearMain();
      var container = document.getElementById('maincontainer');
      var table = document.createElement('table');
      table.id = "reportstable"
      container.appendChild(table);
        $('#reportstable').DataTable({
          columns: data.columns.map((item) => ({ title: item})),
          data: data.result.map((item) => (item.row))
        });
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error)
  }

}

function clearMain() {
  $(document.getElementById("maincontainer")).empty();;
}
