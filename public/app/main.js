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

$(document).ready(function() {
    $('.js-example-basic-single').select2({
      theme: "bootstrap",
      placeholder: "Select Account",
      ajax: {
        url: '/api/accounts',
        dataType: 'json',
      }
    });
});

$('input[name="date"]').on("blur", function(){
  dt = Date.create($(this).val()).format('{yyyy}-{MM}-{dd}')
  $(this).val(dt)
});

$('#addJournal').on('submit', function (e) {
  if (e.isDefaultPrevented()) {
    // handle the invalid form...
  } else {
    e.preventDefault();
    

    window.journal.save($('#addJournal').serializeObject());
    $('#addJournal')[0].reset();
    tableCreate();
    clearJournalDateDescription();
    clearJournalLineItems();
    journal = new Journal();
    $('#journalModal').modal('toggle');
  }
})

const refreshButton = document.getElementById('refresh');

refreshButton.addEventListener('click', async _ => {
  try {
    fetch('/api/journals/')
    .then(response => response.json())
    .then(data => {
      window.transactions = data.Journals;
      tableCreate()
    })
    .catch(error => console.error(error))
} catch { error => console.error(error)
}});

function getJournal(index) {  
  try {
    fetch('/api/journals/'+index)
    .then(response => response.json())
    .then(data => {
      window.journal = data;
      for (var key in journal) {
        try {
          document.getElementById("addJournal").elements[key].value = journal[key]
        } catch(err){
          console.log(err)
        }
      }
      $("#journalModal").modal() 
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error)
  }
}

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
      $("#journalModal").modal() 
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error)
  }
}

function deleteJournal(index) {
  $.ajax({
      type: 'DELETE',
      url: `/api/journals/${index}`,
      success: function(data) {
        window.transactions.splice(index, 1);
        tableCreate();
      },
  });
}

//TODO sean 5 June 2022 - if removing this doesnt cause errors delete at a later date 
//function stripwhitecommas(str) {
  //if (!str || 0 === str.length) {
    //return str
  //} else {
    //return str.toString().replace(/[\s,]+/g,'').trim(); 
  //}
//}

//function stripCents(str) {
  //if (!str || 0 === str.length) {
    //return str
  //} else {
    //if (str.indexOf('.') !== -1) {
      //str = str.substring(0, str.indexOf('.'));
    //}
    ////return str.replace(/[^0-9,]|,[0-9]*$/,''); 
    //return str.replace("/[^\d]/",""); 
  //}
//}

//function toTitleCase(str)
//{
    //return str.replace(/\w\S*/g, function(txt){return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();});
//}

function makeJSON() {
  window.JSONFile = {Transactions: window.transactions};
  var text = JSON.stringify(window.JSONFile, null, '\t');
  download("transactions.json", text);
}

function tableCreate() {
  var tbdy = document.getElementById('transactionstable');
  tbdy.innerHTML = '';
  for (var i = 0; i < window.transactions.length; i++) {
    // Date
    var tr = document.createElement('tr');
    var td = document.createElement('td');
    td.className = 'txntable';
    td.appendChild(document.createTextNode(formatdate(window.transactions[i]._date)))
    tr.appendChild(td)

    // Journal ID
    var td = document.createElement('td');
    td.className = 'txntable';
    var span = document.createElement('span');
    span.appendChild(document.createTextNode(truncate(window.transactions[i].id,12)))
    span.title=window.transactions[i].id;
    td.appendChild(span)

    // Journal ID copy to clipboard
    var svg = document.createElementNS("http://www.w3.org/2000/svg","svg");
    svg.setAttribute('viewBox',"0 0 16 16");
    svg.setAttribute('width',"16");
    svg.setAttribute("height","16");
    var path = document.createElementNS("http://www.w3.org/2000/svg","path");
    path.setAttribute("fill-rule","evenodd");
    path.setAttribute("d","M5.75 1a.75.75 0 00-.75.75v3c0 .414.336.75.75.75h4.5a.75.75 0 00.75-.75v-3a.75.75 0 00-.75-.75h-4.5zm.75 3V2.5h3V4h-3zm-2.874-.467a.75.75 0 00-.752-1.298A1.75 1.75 0 002 3.75v9.5c0 .966.784 1.75 1.75 1.75h8.5A1.75 1.75 0 0014 13.25v-9.5a1.75 1.75 0 00-.874-1.515.75.75 0 10-.752 1.298.25.25 0 01.126.217v9.5a.25.25 0 01-.25.25h-8.5a.25.25 0 01-.25-.25v-9.5a.25.25 0 01.126-.217z");
    var btn = document.createElement('button');
    btn.title=window.transactions[i].id;

    btn.setAttribute('data-param',window.transactions[i].id);
    btn.className = 'clipboard btn btn-sm btn-light';
    btn.onclick = function () {copyToClipboard(this.getAttribute('data-param'));}; 
    svg.appendChild(path);
    btn.appendChild(svg);
    if (window.transactions[i].id.length > 0){
      td.appendChild(btn);
    }
    tr.appendChild(td);

    // Narration
    var td = document.createElement('td');
    td.className = 'txntable';
    td.appendChild(document.createTextNode(window.transactions[i]._description))
    tr.appendChild(td)

    // Amount
    var td = document.createElement('td');
    td.className = 'txntable';
    var amount = document.createTextNode("$" + moneyNumber(window.transactions[i]._amount))
    td.appendChild(amount);
    td.className = 'txntable dollaramount';
    tr.appendChild(td)

    // Edit button
    var td = document.createElement('td');
    var btn = document.createElement('button');
    btn.className = 'btn btn-warning btn-rounded btn-sm';
    btn.setAttribute('data-param-index', i);
    btn.setAttribute('data-param-id',window.transactions[i].id);
    btn.onclick = function () {editJournal(this.getAttribute('data-param-index'),this.getAttribute('data-param-id'));}; 
    btn.innerHTML = "Edit";
    if (window.transactions[i].id.length > 0){
      td.appendChild(btn)
    }
    tr.appendChild(td)

    // Delete button
    var td = document.createElement('td');
    var btn = document.createElement('button');
    btn.className = 'btn btn-danger btn-rounded btn-sm';
    btn.setAttribute('data-param', i);
    btn.setAttribute('data-id', window.transactions[i].id);
    btn.onclick = function () {deleteJournal(this.getAttribute('data-id'));}; 
    btn.innerHTML = "Delete";
    if (window.transactions[i].id.length > 0){
      td.appendChild(btn)
    }
    tr.appendChild(td)
    tbdy.appendChild(tr);
  }
}

$('#journalModal').on('shown.bs.modal', function () {
  $('input[name=date').trigger('focus');
})

$('#journalModal').on('hidden.bs.modal', function () {
  clearJournalDateDescription();
  clearJournalLineItems();
  journal = new Journal();
  $('#addJournal')[0].reset();
  updateTotal()
  $('#saveJournalButton').prop('disabled', true);
  window.now = moment();
})

//TODO sean 5 June 2022 - if removing this doesnt cause errors delete at a later date 
//function formatcomma(element) {
  //return element.toString().replace(/ /g,'').replace(/\B(?=(\d{3})+(?!\d))/g, " ");
//}

$.fn.serializeObject = function()
{
    var o = {};
    var a = this.serializeArray();
    $.each(a, function() {
        if (o[this.name] !== undefined) {
            if (!o[this.name].push) {
                o[this.name] = "";
            }
            if ($(this).is("select")) {
              o[this.name] = $(this).find(':selected').text() || '';
            } else {
              o[this.name] = this.value || '';
            }
        } else {
            if ($(this).is("select")) {
              o[this.name] = $(this).find(':selected').text() || '';
            } else {
              o[this.name] = this.value || '';
            }
        }
    });
    return o;
};

function main() {
  $('#addJournal')[0].reset();
  updateTotal()
  $('#saveJournalButton').prop('disabled', true);
  window.transactions = [];
  window.now = moment();
  try {
    fetch('/api/journals/')
    .then(response => response.json())
    .then(data => {
      window.transactions = data.Journals;
      tableCreate()
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error) }
}
main();
