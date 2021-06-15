let reportRequestMap = new Map()
const TBRequest = {
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
reportRequestMap.set("TrialBalance",TBRequest)

const GLRequest = {
  reports: [
  {
    options: {
      title: "GeneralLedger",
      startdate: "2019-07-01",
      enddate: "2020-06-30"
    },
    columns: [
      "ID",
      "Date",
      "Description",
      "Currency",
      "Amount",
      "Account",
    ]
 }
]}
reportRequestMap.set("GeneralLedger",GLRequest)

function reportRequestGenerator(reportName) {  
  return reportRequestMap.get(reportName)
}

function getReport(reportName) {  
  try {
    fetch('/api/reports/',{
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=utf-8'
      },
      body: JSON.stringify(reportRequestGenerator(reportName))
    })
    .then(response => response.json())
    .then(data => {
      //Clear the page and create a table
      clearMain();
      createConfigWellAndReportsTable(data.options);
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

function createConfigWellAndReportsTable(config) {
  var container = document.getElementById('maincontainer');
  var title = document.createElement('h2');
  title.classList.add("m-3")
  title.classList.add("text-center")
  title.textContent = config.title
  container.appendChild(title)
  var configs = document.createElement('div');
  configs.classList.add("m-3")
  configs.classList.add("card"); 
  configs.classList.add("card-body"); 
  configs.classList.add("bg-light"); 
  var topspan = document.createElement('div')
  topspan.classList.add('row')
  topspan.classList.add('mb-2')
  var datetypecol = document.createElement('div')
  datetypecol.classList.add('col-sm-3')
  
          //<label for="minViewMode">Min view mode
            //<select class="span2 col-md-2 form-control" id="minViewMode" name="minViewMode">
              //<option value="0">0 / days</option>
              //<option value="1">1 / months</option>
              //<option value="2">2 / years</option>
              //<option value="3">3 / decades</option>
              //<option value="4">4 / centuries</option>
            //</select>
          //</label>
  var dropdown = document.createElement('div')
  dropdown.classList.add('dropdown')
  var datedropbutton = document.createElement('button')
  datedropbutton.classList.add('btn')
  datedropbutton.classList.add('btn-secondary')
  datedropbutton.classList.add('dropdown-toggle')
  datedropbutton.setAttribute("aria-has-popup", "true")
  datedropbutton.setAttribute("aria-expanded", "false")
  datedropbutton.id = "dropdownMenuButton"
  datedropbutton.type = "button"
  datedropbutton.setAttribute("data-toggle", "dropdown")
  datedropbutton.textContent = "Last Financial Year"
  var selector = document.createElement('div')
  selector.classList.add('dropdown-menu')
  selector.setAttribute("aria-labelledby", "dropdownMenuButton")
  var option1 = document.createElement('a')
  option1.classList.add('dropdown-item')
  option1.text = "last year"
  selector.appendChild(option1);
  var option2 = document.createElement('a')
  option2.classList.add('dropdown-item')
  option2.text = "last month"
  selector.appendChild(option2)
  dropdown.appendChild(datedropbutton)
  dropdown.appendChild(selector)
  datetypecol.appendChild(dropdown)
  topspan.appendChild(datetypecol)
          //<div class="input-daterange input-group" id="datepicker">
            //<input type="text" class="input-sm form-control" name="start" />
            //<span class="input-group-addon">to</span>
            //<input type="text" class="input-sm form-control" name="end" />
          //</div>
  var datepickercol = document.createElement('div')
  datepickercol.classList.add('col-sm-9')
  datepickercol.id = "datepickercontainer"
  var datepicker = document.createElement('div')
  datepicker.classList.add('input-daterange')
  datepicker.classList.add('input-group')
  datepicker.id = 'datepicker'
  var startinput = document.createElement('input')
  startinput.type = "text"
  startinput.name = "start"
  startinput.classList.add('input-sm')
  startinput.classList.add('form-control')
  datepicker.appendChild(startinput)
  var spanaddon = document.createElement('div')
  spanaddon.classList.add('input-group-text')
  spanaddon.classList.add('input-group-prepend')
  spanaddon.classList.add('input-group-append')
  spanaddon.textContent = "to"
  datepicker.appendChild(spanaddon)
  var endinput = document.createElement('input')
  endinput.type = "text"
  endinput.name = "end"
  endinput.classList.add('input-sm')
  endinput.classList.add('form-control')
  datepicker.appendChild(endinput)
  datepickercol.appendChild(datepicker)
  topspan.appendChild(datepickercol)
  configs.appendChild(topspan)

          //<label for="calendarWeeks" class="inline checkbox">
            //<input id="calendarWeeks" name="calendarWeeks" type="checkbox">
            //Calendar weeks
          //</label>
          //<label for="autoclose" class="inline checkbox">
            //<input id="autoclose" name="autoclose" type="checkbox">
            //Autoclose
          //</label>
          //<label for="todayHighlight" class="inline checkbox">
            //<input id="todayHighlight" name="todayHighlight" type="checkbox">
            //Today highlight
          //</label>
  //var checkboxrow = document.createElement('div')
  //checkboxrow.classList.add('row')
  //checkboxrow.classList.add('mb-2')

          //<br><br>

          //<button class="btn btn-warning" type="button" id="ch_bs">Switch to Bootstrap 2</button>
          //<button class="btn btn-danger" type="reset">Reset to defaults</button>
        //</div>
  var updatespan = document.createElement('div')
  updatespan.classList.add('row')
  var spacecol = document.createElement('div')
  spacecol.classList.add('col-sm-11')
  var buttoncol = document.createElement('div')
  spacecol.classList.add('col-sm-1')
  var update = document.createElement('button')
  update.classList.add('btn')
  update.classList.add('btn-primary')
  update.textContent = "update"
  buttoncol.appendChild(update)
  updatespan.appendChild(spacecol)
  updatespan.appendChild(buttoncol)
  configs.appendChild(updatespan)
  container.appendChild(configs)
  var table = document.createElement('table');
  table.id = "reportstable"
  table.classList.add("m-3")
  container.appendChild(table);
  $('#datepickercontainer .input-daterange').datepicker({
    format: "d MM yyyy",
    todayBtn: "linked",
    clearBtn: true
  });
}
