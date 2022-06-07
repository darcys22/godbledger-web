$(document).ready(function() {
    $('.selectaccount').select2({
      theme: "bootstrap",
      placeholder: "Select Account",
      ajax: {
        url: '/api/reconcile/listexternalaccounts',
        dataType: 'json',
      }
    });
});

$('.selectaccount').on("select2:select", function(event) {
  var value = $(event.currentTarget).find("option:selected").text();
  getTransactions(value);
});

function UnreconciledTransactionsRequest(account)  {
  return {
    options: {
      account: account,
      startdate: "2019-07-01",
      enddate: "2020-06-30"
    },
    columns: [
      "Date",
      "Details",
      "Amount",
      "Currency"
    ]
 }
}

function getTransactions(account) {  
  try {
    fetch('/api/reconcile/listunreconciledtransactions',{
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=utf-8'
      },
      body: JSON.stringify(UnreconciledTransactionsRequest(account))
    })
    .then(response => response.json())
    .then(data => {
      //Clear the page and create a table
      clearMain();
      createConfigWellAndTransactionsTable(data.options);
      cols = data.columns.map((item) => ({ title: item , className: "dt-right"}))
      cols.push({title:"", className: "dt-right"})
      dta = data.result.map((item) => (item.row))
      var table = $('#transactionstable')
      table.DataTable({
        dom: 'Bfrtip',
        select: true,
        columns: cols,
        data: dta,
        columnDefs: [
            {
                // The `data` parameter refers to the data for the cell (defined by the
                // `data` option, which defaults to the column being worked with, in
                // this case `data: 0`.
                "render": function ( data, type, row ) {
                    return `
                      <div class="btn-group">
                        <button type="button" class="btn btn-info btn-sm" onclick="editJournal(1,1)">New Journal</button>
                        <button type="button" class="btn btn-info btn-sm dropdown-toggle dropdown-toggle-split" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                          <span class="sr-only">Toggle Dropdown</span>
                        </button>
                        <div class="dropdown-menu">
                          <a class="dropdown-item" href="#">Something else here</a>
                        </div>
                      </div>
                    `
                },
                "targets": cols.length - 1
            }
        ]
      });
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error)
  }
  //TODO sean the editJournal funct in onclick above needs to actually pass in some parameters

}

function clearMain() {
  $(document.getElementById("maincontainer")).empty();;
}

function clearCSVColumns() {
  var rows = $('#importColumns > tbody > tr');
  rows.each(function(idx, li) {
    var csvcolumn = $(li);
    csvcolumn.remove();
  });
}

function addCSVColumns(index) {
  var tbdy = document.getElementById('importColumns').children[0];
  var tr = document.createElement('tr');
  var td = document.createElement('td');

  //ID of the Journal
  td.appendChild(document.createTextNode(index));
  tr.appendChild(td);

  //Select element for Description of column
  var td = document.createElement('td');
  var select = document.createElement('select');
  select.className = 'js-example-basic-single form-control';
  select.name = `csv-description[${index}][description]`;
  td.appendChild(select);
  tr.appendChild(td);

  //Append the Row to the Table
  tbdy.appendChild(tr);

  $(`select[name ="csv-description[${index}][description]"]`).select2({
    theme: "bootstrap",
    placeholder: "Select Description",
    data: window.CSVColumnTypes,
    width: 'auto',
  })
}

function deleteCSVColumns(index) {
  $('#importColumns > tbody > tr').eq(Number(index)).remove();
}

function processCSVColumns() {
  var rowCount = $('#importColumns > tbody > tr').length;
  if (rowCount < window.csvColumns) {
    for (let i = rowCount; i < window.csvColumns; i++) {
      addCSVColumns(i);
    }
  } else {
    for (let i = rowCount; i >= window.csvColumns; i--) {
      deleteCSVColumns(i);
    }
  }
}

function updateCSVInput() {
  if ($('input[name=numberColumns]').val() > 0)
  {
    window.csvColumns = Number($('input[name=numberColumns]').val());
    processCSVColumns();
  }
}

function createConfigWellAndTransactionsTable(config) {
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
  table.id = "transactionstable"
  table.classList.add("m-3")
  container.appendChild(table);
}

window.CSVColumnTypes = [
    {
        id: 0,
        text: 'empty'
    },
    {
        id: 1,
        text: 'date'
    },
    {
        id: 2,
        text: 'description'
    },
    {
        id: 3,
        text: 'amount'
    },
    {
        id: 4,
        text: 'debit'
    },
    {
        id: 5,
        text: 'credit'
    }
];

function editJournal(index,id) {  
  try {
    fetch('/api/journals/'+id)
    .then(response => response.json())
    .then(data => {
      $('#addJournal')[0].reset();
      journal = new Journal();
      clearJournalDateDescription();
      clearJournalLineItems();
      journal._lineItemCount = 0;
      journal.setID(id);
      document.getElementsByName("date")[0].value = formatformaldate(data._date);
      document.getElementsByName("narration")[0].value = data._narration;
      for (var lineItem in data._lineItems) {
        journal.addNewLineItem();
        document.getElementsByName("line-item[" +(journal._lineItemCount)+ "][narration]")[0].value = data._lineItems[lineItem]._description;
        var amount = parseInt(data._lineItems[lineItem]._amount);
        if (amount > 0) {
          document.getElementsByName("line-item[" +journal._lineItemCount+ "][debit]")[0].value = amount;
        } else {
          document.getElementsByName("line-item[" +journal._lineItemCount+ "][credit]")[0].value = -amount;
        }
        var account = data._lineItems[lineItem]._account;
        var accountSelect = $(`select[name ="line-item[${journal._lineItemCount}][account]"]`);
        var option = new Option(account, '0', true, true);
        accountSelect.append(option).trigger('change');
      }
      updateTotal();
      $("#reconcileModal").modal() 
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error)
  }
}

$('#reconcileModal').on('shown.bs.modal', function () {
  $('input[name=date').trigger('focus');
})

$('#reconcileModal').on('hidden.bs.modal', function () {
  clearJournalDateDescription();
  clearJournalLineItems();
  journal = new Journal();
  $('#addJournal')[0].reset();
  updateTotal()
  $('#saveJournalButton').prop('disabled', true);
  window.now = moment();
})

var journal = new Journal();

// ...and hook up the add new line item button
var newLineItemButton = document.getElementById("addNewLineItemButton");
if (newLineItemButton.addEventListener) {
    newLineItemButton.addEventListener('click', function(e) {
        e.preventDefault();
        journal.addNewLineItem();
    }, false);
}
else if (newLineItemButton.attachEvent) {
    newLineItemButton.attachEvent('onclick', function(e) {
        e.preventDefault();
        journal.addNewLineItem();
    });
}
else {
    // Very old browser, complain
}

function handleFiles(event) {
  window.csvupload = event.target.files[0];
  var csvlabel = document.getElementById("CSVLabel");
  csvlabel.innerHTML = window.csvupload.name;
}

var csvcheckbutton = document.getElementById("saveCSVButton");
if (csvcheckbutton.addEventListener) {
    csvcheckbutton.addEventListener('click', function(e) {
        e.preventDefault();
        submitUpload();
    }, false);
}
else if (csvcheckbutton.attachEvent) {
    csvcheckbutton.attachEvent('onclick', function(e) {
        e.preventDefault();
        submitUpload();
    });
}
else {
    // Very old browser, complain
}

function submitUpload() {
  var columns = [];
  $('#importColumns > tbody > tr').each(function(index, tr) { 
    columns.push($(tr).find('select :selected').text() || '');
  });

  var startRow = 0
  if ($('input[name=startRow]').val() > 0)
    startRow = Number($('input[name=startRow]').val());

  var lastRow = 0
  if ($('input[name=lastRow]').val() > 0)
    lastRow = Number($('input[name=lastRow]').val());

  const CSVRequest = {
    options: {
      startRow,
      lastRow
    },
    columns,
    file: window.csvupload
  }

  console.log(CSVRequest);

  try {
    fetch('/api/reconcile/UploadCSV',{
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=utf-8'
      },
      body: JSON.stringify(CSVRequest)
    })
    .then(response => response.json())
    .then(data => {
      console.log("Successfull send to backend")
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error) }
  // TODO sean step 4 submit for review
  // TODO show a success/reject and proceed
}

function main() {
  clearCSVColumns();
  $('input[name=numberColumns]').val("5");
  window.csvColumns = 5;
  processCSVColumns();
}
main();
