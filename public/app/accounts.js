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

class Account {
    constructor() {
      this.id = "";
      this.name = "Display Me";
      this.tags = [];
      this._tagCount = 0;
    }

    addNewTag(tag) {
      this.tags.push(tag);
      this._tagCount += 1;
      addTag(this._tagCount);
    }

    save(editAccountForm) {

      var tags = Object.keys(editAccountForm).filter(function(name) {
        return name.includes("tag");
      });
      var i = 0;
      for (i = 0; i < tags.length; i++) {
        this._lineitems[parseFloat(filtered[1], 10)] = lineitem;
      }
      this.id = editAccountForm.id;
      this.name  = editAccountForm.name;
      $.ajax({
          type: 'POST',
          url: '/api/journals/'+this.id,
          data: JSON.stringify(this,
            function(k, v) {
              if (k === '_amount') {
                return v.toString();
              }
              return v;
            }
          ),
          success: function(data) {},
          contentType: "application/json",
          dataType: 'json'
      });
    }
}

function addTag(index) {
  var tbdy = document.getElementById('tags');
  var tr = document.createElement('tr');
  tr.id = `tag[${index}]`;
  var td = document.createElement('td');

  //ID of the tag
  td.appendChild(document.createTextNode(index));
  td.align="center"
  tr.appendChild(td);

  //Input for name of tag
  var td = document.createElement('td');
  var input  = document.createElement('input');
  input.className = 'form-control';
  input.setAttribute('data-lpignore', "true");
  input.name = `tag[${index}][name]`;
  input.type = "text";
  input.setAttribute('tabindex', index*4+3);
  td.appendChild(input);
  tr.appendChild(td);

  // Delete button
  var td = document.createElement('td');
  td.align="center"
  var btn = document.createElement('button');
  btn.className = 'btn btn-danger btn-rounded btn-sm';
  btn.setAttribute('data-param', index);
  btn.innerHTML = "Delete";
  td.appendChild(btn)
  tr.appendChild(td)

  //Append the Row to the Table
  tbdy.appendChild(tr);

  if (btn.addEventListener) {
      btn.addEventListener('click', function(e) {
          e.preventDefault();
          deleteTag(this.getAttribute('data-param'));
      }, false);
  }
  else if (btn.attachEvent) {
      btn.attachEvent('onclick', function(e) {
          e.preventDefault();
          deleteTag(this.getAttribute('data-param'));
      });
  }
}

function deleteTag(index) {
  var accountName = window.account.name;
  var tag = window.account.tags[index-1];
  try {
    fetch(`/api/accounttags/deletetag/${accountName}/${tag}`,{
      method: 'DELETE'
    })
    .then(data => {
      window.account.tags.splice(index - 1, 1);
      const elem = document.getElementById("tag[" +(index)+ "]");
      elem.parentNode.removeChild(elem);
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error)
  }
}

function editAccount(id) {  
  try {
    fetch('/api/accounts/'+id)
    .then(response => response.json())
    .then(data => {
      $('#editAccount')[0].reset();
      clearModalTags();
      window.account = new Account();
      window.account._tagCount = 0;
      window.account.id = data.id;
      window.account.name = data.name;
      document.getElementsByName("code")[0].value = data.id;
      document.getElementsByName("name")[0].value = data.name;
      const tags = data.tags;
      for (const tg of tags) {
        window.account.addNewTag(tg);
        document.getElementsByName("tag[" +(window.account._tagCount)+ "][name]")[0].value = tg;
      }
      $("#accountsModal").modal() 
    })
    .catch(error => console.error(error))
  } catch { error => console.error(error)
  }

}

function clearModalTags() {
  var rows = $('#tags > tr');
  rows.each(function(idx, li) {
    var lineItem = $(li);
    lineItem.remove();
  });
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

function importAccounts(name) {  
  fetch('/api/accounts/import',{
    method: 'POST',
    headers: {
      'Content-Type': 'application/json;charset=utf-8'
    },
    body: JSON.stringify({
      name: name,
    })
  })
  .then(response => response.json())
  .then(data => {
    getAccounts();
  })
  .catch(error => console.error(error))
}

// hook up the import default chart of accounts
var importDefaultAccountsButton = document.getElementById("importDefaultAccounts");
if (importDefaultAccountsButton.addEventListener) {
    importDefaultAccountsButton.addEventListener('click', function(e) {
        e.preventDefault();
        importAccounts("default");
        //window.location.reload(true)
    }, false);
}
else if (newTagButton.attachEvent) {
    importDefaultAccountsButton.attachEvent('onclick', function(e) {
        e.preventDefault();
        importAccounts("default");
        //window.location.reload(true)
    });
}
else {
    // Very old browser, complain
}



window.account = new Account();

// ...and hook up the add new tag button
var newTagButton = document.getElementById("addNewTagButton");
if (newTagButton.addEventListener) {
    newTagButton.addEventListener('click', function(e) {
        e.preventDefault();
        var newTagName = document.getElementById("newTagName");
        window.account.addNewTag(newTagName);
        document.getElementsByName("tag[" +(window.account._tagCount)+ "][name]")[0].value = newTagName.value;
    }, false);
}
else if (newTagButton.attachEvent) {
    newTagButton.attachEvent('onclick', function(e) {
        e.preventDefault();
        var newTagName = document.getElementById("newTagName");
        window.account.addNewTag(newTagName);
        document.getElementsByName("tag[" +(window.account._tagCount)+ "][name]")[0].value = newTagName.value;
    });
}
else {
    // Very old browser, complain
}

function main() {
  getAccounts();
}
main();
