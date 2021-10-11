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
      $(".editBtn").on('click', function(event){
          event.stopPropagation();
          event.stopImmediatePropagation();
          console.log("testsing1 edit");
          console.log(this.getAttribute('data-param'));
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
              throw new Error('Network response was not ok');
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

function main() {
  getAccounts();
}
main();
