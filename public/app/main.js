class Journal {
    constructor() {
      this.date = new Date();
      this.narration = "Display Me";
      this.lineitems = [];
      
    }

    DisplayData() {
        alert(this.narration);
    }

    addNewLineItem() {
        alert(this.narration);
    }

    save(journalForm) {
      this.narration = journalForm.narration;
      console.log(journalForm)
    }

  
}

var journal = new Journal();

// ...and hook up the add new line item button
var newLineItemButton = document.getElementById("addNewLineItemButton");
if (newLineItemButton.addEventListener) {
    newLineItemButton.addEventListener('click', function() {
        journal.addNewLineItem();
    }, false);
}
else if (newLineItemButton.attachEvent) {
    newLineItemButton.attachEvent('onclick', function() {
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
        url: '/api/accounts/list',
        dataType: 'json',
      }
    });
});

//$('#addJournal').validator().on('submit', function (e) {
$('#addJournal').on('submit', function (e) {
  if (e.isDefaultPrevented()) {
    // handle the invalid form...
  } else {
    e.preventDefault();

    window.journal.save($('#addJournal').serializeObject());
    $('#addJournal')[0].reset();
		//tableCreate();
    $('#journalModal').modal('toggle');
  }
})

const refreshButton = document.getElementById('refresh');

refreshButton.addEventListener('click', async _ => {
  try {
    fetch('/api/journals/')
    .then(response => response.json())
    .then(data => {
      window.transactions = JSON.parse(data);
      tableCreate()
    })
    .catch(error => console.error(error))
} catch { error => console.error(error)
}});

function editJournal(index) {
  var journal = window.transactions[index];
  deleteJournal(index) 
  for (var key in journal) {
    try {
      document.getElementById("addJournal").elements[key].value = journal[key]
    } catch(err){
    }
  }
  $("#journalModal").modal() 
  //$('#addJournal').validator()

}
function deleteJournal(index) {
  window.transactions.splice(index, 1);
  //tableCreate();
}

function stripwhitecommas(str) {
  if (!str || 0 === str.length) {
    return str
  } else {
    return str.toString().replace(/[\s,]+/g,'').trim(); 
  }
}

function stripCents(str) {
  if (!str || 0 === str.length) {
    return str
  } else {
    if (str.indexOf('.') !== -1) {
      str = str.substring(0, str.indexOf('.'));
    }
    //return str.replace(/[^0-9,]|,[0-9]*$/,''); 
    return str.replace("/[^\d]/",""); 
  }
}

function toTitleCase(str)
{
    return str.replace(/\w\S*/g, function(txt){return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();});
}

function download(filename, text) {
  var element = document.createElement('a');
  element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(text));
  element.setAttribute('download', filename);

  element.style.display = 'none';
  document.body.appendChild(element);

  element.click();

  document.body.removeChild(element);
}

function makeJSON() {
  window.JSONFile = {Transactions: window.transactions};
  var text = JSON.stringify(window.JSONFile, null, '\t');
  download("transactions.json", text);
}

function moneyNumber(x) {
    return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
}

function tableCreate() {
    var tbdy = document.getElementById('transactionstable');
    tbdy.innerHTML = '';
    for (var i = 0; i < window.transactions.length; i++) {
        var tr = document.createElement('tr');
        var td = document.createElement('td');
        td.appendChild(document.createTextNode(window.transactions[i].date))
        tr.appendChild(td)
        var td = document.createElement('td');
        td.appendChild(document.createTextNode(window.transactions[i].id))
        tr.appendChild(td)
        var td = document.createElement('td');
        td.appendChild(document.createTextNode(window.transactions[i].desc))
        tr.appendChild(td)
        var td = document.createElement('td');
        td.appendChild(document.createTextNode("$" + moneyNumber(window.transactions[i].amount)));
        tr.appendChild(td)
        var td = document.createElement('td');
        var btn = document.createElement('button');
        btn.className = 'btn btn-warning btn-rounded btn-sm';
        btn.setAttribute('data-param', i);
        //btn.onclick = function () {editJournal(this.getAttribute('data-param'));}; 
        btn.innerHTML = "Edit";
        td.appendChild(btn)
        tr.appendChild(td)
        var td = document.createElement('td');
        var btn = document.createElement('button');
        btn.className = 'btn btn-danger btn-rounded btn-sm';
        btn.setAttribute('data-param', i);
        btn.onclick = function () {deleteJournal(this.getAttribute('data-param'));}; 
        btn.innerHTML = "Delete";
        td.appendChild(btn)
        tr.appendChild(td)
        tbdy.appendChild(tr);
    }
}

function formatcomma(element) {
  return element.toString().replace(/ /g,'').replace(/\B(?=(\d{3})+(?!\d))/g, " ");
}

function formatdate(element) {
  element.value = moment(element.value, ["DDMMYYYY","DDMMMMYYYY", "DoMMMMYYYY", "DoMMYYYY"], false).format('Do MMMM YYYY');
}

$.fn.serializeObject = function()
{
    var o = {};
    var a = this.serializeArray();
    $.each(a, function() {
        if (o[this.name] !== undefined) {
            if (!o[this.name].push) {
                o[this.name] = "";
            }
            o[this.name] = this.value || '';
        } else {
          if (this.name.includes("[")){
            var separators = ['\\\[', '\\\]'];
            var tokens = this.name.split(new RegExp(separators.join('|'), 'g'));
            var filtered = tokens.filter(function (el) {
              return el != "";
            });
            console.log(filtered[0]);
            //o[filtered[0]][filtered[1]][filtered[2]] = this.value || '';
          } else {
            o[this.name] = this.value || '';
          }
        }
    });
    return o;
};

function main() {
  window.transactions = [];
  window.now = moment();
}
main();
