const reportRequest = {
  reports: [
  {
    Options: {
      Title: "TrialBalance",
      StartDate: "2019-07-01",
      EndDate: "2020-06-30"
    },
    Columns: [
      "Date",
      "AccountName",
      "Description",
      "Amount"
    ]
 }
]}

function getReport(reportName) {  
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
      //Copy the data directly into a table for it
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error)
  }

}
