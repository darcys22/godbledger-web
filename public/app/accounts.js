$('#addAccount').on('submit', function (e) {
  if (e.isDefaultPrevented()) {
    // handle the invalid form...
  } else {
    e.preventDefault();
    name = $('input[name=account]').val();
    fetch('/api/accounts',{
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=utf-8'
      },
      body: JSON.stringify({
        name: name,
        tags: [
          "main"
        ]
      })
    })
    .then(response => response.json())
    .then(data => {
      console.log(data);
      getAccounts();
    })
    .catch(error => console.error(error))
    $('input[name=account]').val('');
  }
})

function getAccounts() {  
  try {
    fetch('/api/accounts', {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json;charset=utf-8'
      }
    })
    .then(response => response.json())
    .then(data => {
      //Clear the page and create a table
      clearMain();
      createAccountsTable();
      var table = $('#accountstable')
      dta = data.results.map((item) => ([item.id, item.text]));
      table.DataTable({
        columns: [
          { title: "ID", className: "dt-right" },
          { title: "Account", className: "dt-right" },
          { title: "", className: "dt-right" }
        ],
        data: dta,
        columnDefs: [
            {
                "render": function ( data, type, row ) {
                    var span = document.createElement('span');
                    var editBtn = document.createElement('button');
                    editBtn.className = 'btn btn-info btn-rounded btn-sm m-2 editBtn';
                    editBtn.setAttribute('data-param',row[1]);
                    editBtn.innerHTML = "Edit";
                    span.appendChild(editBtn)
                    var deleteBtn = document.createElement('button');
                    deleteBtn.className = 'btn btn-danger btn-rounded btn-sm m-2 deleteBtn';
                    deleteBtn.setAttribute('data-param',row[1]);
                    deleteBtn.innerHTML = "Delete";
                    span.appendChild(deleteBtn)
                    return span.outerHTML;
                },
                "targets": 2
            }
        ]
      });
      //TODO sean remove this
      $(".editBtn").on('click', function(event){
          event.stopPropagation();
          event.stopImmediatePropagation();
          editAccount(this.getAttribute('data-param'));
      });

      $(".deleteBtn").on('click', function(event){
          event.stopPropagation();
          event.stopImmediatePropagation();
          fetch('/api/accounts/'+this.getAttribute('data-param'),{
            method: 'DELETE',
            headers: {
              'Content-Type': 'application/json;charset=utf-8'
            }
          })
          .then(response => {
            if (!response.ok) {
              response.text().then(data => {
                if (data.includes("FOREIGN KEY constraint failed")) {
                  showMessage("Could not delete Account: Transactions exist using this account", "Error")
                } else {
                  showMessage(data, "Error")
                }
              })
            }
            getAccounts();
            return response;
          })
          .catch(error => console.error(error))
      });
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error)
  }

}

function editAccount(id) {  
  console.log(id);
  try {
    fetch('/api/accounts/'+id)
    .then(response => response.json())
    .then(data => {
      console.log(data);
      //$('#editAccounts')[0].reset();
      //TODO sean this needs to populate the edit Account modal
      //account = new Account();
      //clearAccountModal();
      //clearAccountTags();
      //account._tagCount = 0;
      //account.setID(id);
      //document.getElementsByName("code")[0].value = data._code;
      //document.getElementsByName("name")[0].value = data._name;
      //for (var tag in data._tags) {
        //account.addNewTag();
        //document.getElementsByName("line-item[" +(journal._lineItemCount)+ "][narration]")[0].value = data._lineItems[lineItem]._description;
      //}
      $("#accountsModal").modal() 
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error)
  }

}

function clearMain() {
  $(document.getElementById("maincontainer")).empty();;
}

function createAccountsTable() {
  var container = document.getElementById('maincontainer');
  var table = document.createElement('table');
  table.id = "accountstable"
  table.classList.add("m-3")
  container.appendChild(table);
}

function showMessage(message, messagetype) {
  var cssclass;
  switch (messagetype) {
    case 'Success':
        cssclass = 'alert-success'
        break;
    case 'Error':
        cssclass = 'alert-danger'
        break;
    case 'Warning':
        cssclass = 'alert-warning'
        break;
    default:
        cssclass = 'alert-info'
  }
  $('#alert_container').append('<div id="alert_div" style="margin: 0 0.5%; -webkit-box-shadow: 3px 4px 6px #999;" class="alert ' + cssclass + '"><a href="#" class="close" data-dismiss="alert" aria-label="close">&times;</a><strong>' + messagetype + '!</strong> <span>' + message + '</span></div>');

  setTimeout(function () {
    $("#alert_div").fadeTo(2000, 500).slideUp(500, function () {
      $("#alert_div").remove();
    });
  }, 3000);//3000=5 seconds
}

function main() {
  getAccounts();
}
main();
