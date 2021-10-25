let reportRequestMap = new Map()
const TBRequest = {
  reports: [
  {
    options: {
      title: "TrialBalance",
      startdate: "1970-01-01",
      enddate: "3020-06-30"
    },
    columns: [
      "Accountname",
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
      startdate: "1970-01-01",
      enddate: "3020-06-30"
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

function getReport(reportName, config) {  
  if (arguments.length == 1) // Means second parameter is not passed
  {
    config = reportRequestGenerator(reportName);
  }
  try {
    fetch('/api/reports/',{
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=utf-8'
      },
      body: JSON.stringify(config)
    })
    .then(response => response.json())
    .then(data => {
      //Clear the page and create a table
      clearMain();
      createConfigWellAndReportsTable(data.options);
      if (data.result === null) {
        var container = document.getElementById('maincontainer');
        var title = document.createElement('h2');
        title.classList.add("m-3")
        title.classList.add("text-center")
        title.textContent = "No Results"
        container.appendChild(title)
      } else {
        $('#reportstable').DataTable({
          dom: 'Bfrtip',
          columns: data.columns.map((item) => ({ title: item, className: "dt-right"})),
          data: data.result.map((item) => (item.row))
        });
      }
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error)
  }

}

function clearMain() {
  $(document.getElementById("maincontainer")).empty();
}

function createConfigWellAndReportsTable(config) {
  var container = document.getElementById('maincontainer');
  var title = document.createElement('h2');
  title.classList.add("m-3")
  title.classList.add("text-center")
  title.textContent = config.title
  container.appendChild(title)
  var form = document.createElement('form');
  form.id = "updateReport"
  var configs = document.createElement('div');
  configs.classList.add("m-3")
  configs.classList.add("card"); 
  configs.classList.add("card-body"); 
  configs.classList.add("bg-light"); 
  configs.classList.add('m-5')
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
  datedropbutton.classList.add('m-1')
  datedropbutton.setAttribute("aria-has-popup", "true")
  datedropbutton.setAttribute("aria-expanded", "false")
  datedropbutton.id = "dropdownMenuButton"
  datedropbutton.type = "button"
  datedropbutton.setAttribute("data-toggle", "dropdown")
  datedropbutton.textContent = "All Time"
  var selector = document.createElement('div')
  selector.classList.add('dropdown-menu')
  selector.setAttribute("aria-labelledby", "dropdownMenuButton")
  var option1 = document.createElement('a')
  option1.classList.add('dropdown-item')
  option1.text = "Last Financial Year"
  option1.onclick = function(){dateQuickSelect("last-financial-year")};
  selector.appendChild(option1);
  var option2 = document.createElement('a')
  option2.classList.add('dropdown-item')
  option2.text = "Last Calendar Year"
  option2.onclick = function(){dateQuickSelect("last-calendar-year")};
  selector.appendChild(option2);
  var option3 = document.createElement('a')
  option3.classList.add('dropdown-item')
  option3.text = "This Financial Year"
  option3.onclick = function(){dateQuickSelect("this-financial-year")};
  selector.appendChild(option3);
  var option4 = document.createElement('a')
  option4.classList.add('dropdown-item')
  option4.text = "This Calendar Year"
  option4.onclick = function(){dateQuickSelect("this-calendar-year")};
  selector.appendChild(option4);
  var option5 = document.createElement('a')
  option5.classList.add('dropdown-item')
  option5.text = "Last Month"
  option5.onclick = function(){dateQuickSelect("last-month")};
  selector.appendChild(option5)
  var option6 = document.createElement('a')
  option6.classList.add('dropdown-item')
  option6.text = "This Month"
  option6.onclick = function(){dateQuickSelect("this-month")};
  selector.appendChild(option6)
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
  startinput.style.textAlign = "center"
  startinput.classList.add('input-sm')
  startinput.classList.add('form-control')
  startinput.classList.add('m-1')
  datepicker.appendChild(startinput)
  var spanaddon = document.createElement('div')
  spanaddon.classList.add('input-group-text')
  spanaddon.classList.add('input-group-prepend')
  spanaddon.classList.add('input-group-append')
  spanaddon.classList.add('m-1')
  spanaddon.textContent = "to"
  datepicker.appendChild(spanaddon)
  var endinput = document.createElement('input')
  endinput.type = "text"
  endinput.name = "end"
  endinput.style.textAlign = "center"
  endinput.classList.add('input-sm')
  endinput.classList.add('form-control')
  endinput.classList.add('m-1')
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
  update.type = 'submit'
  update.classList.add('btn')
  update.classList.add('btn-primary')
  update.textContent = "update"
  buttoncol.appendChild(update)
  updatespan.appendChild(spacecol)
  updatespan.appendChild(buttoncol)
  configs.appendChild(updatespan)
  form.appendChild(configs)
  container.appendChild(form)
  var table = document.createElement('table');
  table.id = "reportstable"
  table.classList.add("m-3")
  container.appendChild(table);

  $('#updateReport').on('submit', function (e) {
    if (e.isDefaultPrevented()) {
      // handle the invalid form...
    } else {
      e.preventDefault();
      console.log(config);
      newconfig = reportRequestGenerator(config.title);
      newconfig.reports[0].options.startdate = $('input[name=start]').val();
      newconfig.reports[0].options.enddate = $('input[name=end]').val();
      requestParams = getReport(newconfig.title, newconfig);
    }
  });

  $('input[name="start"], input[name="end"]').on("blur", function(){
    dt = Date.create($(this).val()).format('{yyyy}-{MM}-{dd}')
    $(this).val(dt)
    $('#dropdownMenuButton').html("Custom");
  });

}

function dateQuickSelect(rangeType) {
  console.log(rangeType);
  switch (rangeType) {
  case "last-financial-year":
    $('#dropdownMenuButton').html("Last Financial Year");
    startDate = Date.create("the beginning of july").addMonths(-12);
    $('input[name=start]').val(startDate.format('{yyyy}-{MM}-{dd}'));
    $('input[name=end]').val(startDate.addMonths(12).addDays(-1).format('{yyyy}-{MM}-{dd}'));
    break;
  case "last-calendar-year":
    $('#dropdownMenuButton').html("Last Calendar Year");
    $('input[name=start]').val(Date.create("the beginning of last january").format('{yyyy}-{MM}-{dd}'));
    $('input[name=end]').val(Date.create("the end of last december").format('{yyyy}-{MM}-{dd}'));
    break;
  case "this-financial-year":
    $('#dropdownMenuButton').html("This Financial Year");
    startDate = Date.create("the beginning of july");
    $('input[name=start]').val(startDate.format('{yyyy}-{MM}-{dd}'));
    $('input[name=end]').val(startDate.addMonths(12).addDays(-1).format('{yyyy}-{MM}-{dd}'));
    break;
  case "this-calendar-year":
    $('#dropdownMenuButton').html("This Calendar Year");
    $('input[name=start]').val(Date.create("the beginning of january").format('{yyyy}-{MM}-{dd}'));
    $('input[name=end]').val(Date.create("the end of december").format('{yyyy}-{MM}-{dd}'));
    break;
  case "last-month":
    $('#dropdownMenuButton').html("Last Month");
    $('input[name=start]').val(Date.create("the beginning of last month").format('{yyyy}-{MM}-{dd}'));
    $('input[name=end]').val(Date.create("the end of last month").format('{yyyy}-{MM}-{dd}'));
    break;
  case "this-month":
    $('#dropdownMenuButton').html("This Month");
    $('input[name=start]').val(Date.create("the beginning of this month").format('{yyyy}-{MM}-{dd}'));
    $('input[name=end]').val(Date.create("the end of this month").format('{yyyy}-{MM}-{dd}'));
    break;
  default:
    console.log("Unknown date range selected");
}
}

function main() {
}
main();
