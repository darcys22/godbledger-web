$('#addContractor').validator().on('submit', function (e) {
  if (e.isDefaultPrevented()) {
    // handle the invalid form...
  } else {
    e.preventDefault();
		window.transactions.push($('#addContractor').serializeObject());
		$('#addContractor')[0].reset();
		$("#fybox").val(window.endFY.format("YYYY"));
		tableCreate();
    $('#contractorModal').modal('toggle');
    openvalidate();
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

//Todo(sean):
function openvalidate() {
  var validatebutton = document.getElementById('convert');
  if (window.transactions.length > 0 && !jQuery.isEmptyObject(window.payer)) {
    validatebutton.title = "Validate before creating TPAR File"
    validatebutton.disabled = false
    validatebutton.onclick = function() {validateTPAR()};
  } else {
    validatebutton.title = "Please add a Contractor and Payer details to create file"
    validatebutton.disabled = true
    validatebutton.onclick = function(){};
  }
}
function editContractor(index) {
  var contractor = window.transactions[index];
  deleteContractor(index) 
  for (var key in contractor) {
    try {
      document.getElementById("addContractor").elements[key].value = contractor[key]
    } catch(err){
    }
  }
  $("#contractorModal").modal() 
  //$('#addContractor').validator()

}
function deleteContractor(index) {
  window.transactions.splice(index, 1);
  tableCreate();
  openvalidate();
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
        btn.className = 'btn btn-warning';
        btn.setAttribute('data-param', i);
        //btn.onclick = function () {editContractor(this.getAttribute('data-param'));}; 
        btn.innerHTML = "Edit";
        td.appendChild(btn)
        tr.appendChild(td)
        var td = document.createElement('td');
        var btn = document.createElement('button');
        btn.className = 'btn btn-danger';
        btn.setAttribute('data-param', i);
        btn.onclick = function () {deleteContractor(this.getAttribute('data-param'));}; 
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


function main() {
  window.transactions = [];
  window.now = moment();
}
main();
