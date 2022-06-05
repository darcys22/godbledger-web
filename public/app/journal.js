class LineItem {
    constructor() {
      this._date = moment().format();
      this._description = "";
      this._currency= window.user.defaultcurrency;
      this._account = "";
      this._amount = 0;
    }

    set amount(amount) {
      this._amount = amount;
    }
    get amount() {
      return this._amount;
    }

    set account(account) {
      this._account = account;
    }
    get account() {
      return this._account;
    }

    set description(description) {
      this._description = description;
    }
    get description() {
      return this._description;
    }

    set date(date) {
        this._date = moment(date).format();
    }
    get date() {
        return this._date;
    }

    isEmpty() {
			return ((!this._account || 0 === this._account.length) || !this._account.trim()) && (!this._amount || this._amount === 0);
    }
}
class Journal {
    constructor() {
      this._date = moment().format();
      this.id = "";
      this._narration = "Display Me";
      this._lineitems = [];
      this._lineItemCount = 0;

      var i;
      for (i = 0; i < 3; i++) {
        this.addNewLineItem();
      } 
    }

    setID(journalID) {
      this.id += journalID;
    }

    addNewLineItem() {
      this._lineItemCount += 1;
      addLineItem(this._lineItemCount);
    }

    save(journalForm) {

      var lineitemKeys = Object.keys(journalForm).filter(function(name) {
        return name.includes("line-item");
      });
      var journalDate = moment(journalForm.date, "YYYY-MM-DD").format();
      var i = 0;
      for (i = 0; i < lineitemKeys.length; i++) {
        if (lineitemKeys[i].includes("[")){
          var separators = ['\\\[', '\\\]'];
          var tokens = lineitemKeys[i].split(new RegExp(separators.join('|'), 'g'));
          var filtered = tokens.filter(function (el) {
            return el != "";
          });
          if (this._lineitems[parseInt(filtered[1], 10)] != undefined) {
            var lineitem = this._lineitems[parseInt(filtered[1], 10)];
          } else {
            var lineitem = new LineItem();
            lineitem._date = journalDate;
          }
          switch(filtered[2]) {
            case "narration":
              lineitem._description = journalForm[lineitemKeys[i]];
              break;
            case "account":
              lineitem._account = $(`select[name ="${lineitemKeys[i]}"]`).text();
              break;
            case "debit":
              if (lineitem._amount == 0 && journalForm[lineitemKeys[i]]) {
                lineitem._amount = parseFloat(journalForm[lineitemKeys[i]],10).toFixed(2) * 1;
              }
              break;
            case "credit":
              if (lineitem._amount == 0 && journalForm[lineitemKeys[i]]) {
                lineitem._amount = parseFloat(journalForm[lineitemKeys[i]],10).toFixed(2) * -1;
              }
              break;
            default:
              console.log("could not identify" + lineitemKeys[i])
          }

          this._lineitems[parseFloat(filtered[1], 10)] = lineitem;
        }
      }
      this._narration = journalForm.narration;
      this._lineitems.splice(0, 1);
      this._lineitems = this._lineitems.filter(function (el) {
        return !el.isEmpty();
      });
      this._date = journalDate;
      this._lineItemCount = this._lineitems.length;
      window.transactions.unshift({"id":"","_date":this._date,"_description":this._narration,"_amount":Math.abs(this._lineitems[0]._amount).toFixed(2)})
      if(this.id == "")
      {
        $.ajax({
            type: 'POST',
            url: '/api/journals',
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
      } else {
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
}

function addLineItem(index) {
  var tbdy = document.getElementById('journal');
  var tr = document.createElement('tr');
  var td = document.createElement('td');

  //ID of the Journal
  td.appendChild(document.createTextNode(index));
  tr.appendChild(td);

  //Input for Narration of line item
  var td = document.createElement('td');
  var input  = document.createElement('input');
  input.className = 'form-control';
  input.setAttribute('data-lpignore', "true");
  input.name = `line-item[${index}][narration]`;
  input.type = "text";
  input.setAttribute('tabindex', index*4+3);
  td.appendChild(input);
  tr.appendChild(td);

  //Select element for Account of line item
  var td = document.createElement('td');
  var select = document.createElement('select');
  select.className = 'js-example-basic-single form-control';
  select.name = `line-item[${index}][account]`;
  select.setAttribute('tabindex', index*4+4);
  td.appendChild(select);
  tr.appendChild(td);

  //Input for Debit Amount of line item
  var td = document.createElement('td');
  var input  = document.createElement('input');
  input.className = 'form-control money';
  input.setAttribute('data-lpignore', "true");
  input.name = `line-item[${index}][debit]`;
  input.type = "number";
  input.setAttribute("onchange", `document.getElementsByName("line-item[${index}][credit]")[0].value="";updateTotal();`);;
  input.setAttribute("min", 0);;
  input.setAttribute("step", 0.01);;
  input.setAttribute('tabindex', index*4+5);
  td.appendChild(input);
  tr.appendChild(td);

  //Input for Credit Amount of line item
  var td = document.createElement('td');
  var input  = document.createElement('input');
  input.className = 'form-control money';
  input.setAttribute('data-lpignore', "true");
  input.name = `line-item[${index}][credit]`;
  input.type = "number";
  input.setAttribute("onchange", `document.getElementsByName("line-item[${index}][debit]")[0].value="";updateTotal();`);;
  input.setAttribute("min", 0);;
  input.setAttribute('tabindex', index*4+6);
  input.setAttribute("step", 0.01);;
  td.appendChild(input);
  tr.appendChild(td);

  //Append the Row to the Table
  tbdy.appendChild(tr);

  $(`select[name ="line-item[${index}][account]"]`).select2({
    theme: "bootstrap",
    placeholder: "Select Account",
    ajax: {
      url: '/api/accounts',
      dataType: 'json',
    }
  }).on("change",function(){
    updateTotal()
  });
}

function updateTotal()
{
  $('#saveJournalButton').prop('disabled', false);
  var DRTotal = 0;
  var CRTotal = 0;

  for (var i = 1; i <= journal._lineItemCount; i++) {
    var DRAmount = parseFloat(document.getElementsByName(`line-item[${i}][debit]`)[0].value);
    if (!isNaN(DRAmount) && DRAmount >=0) { 
      DRTotal += DRAmount; 
      $(`input[name ="line-item[${i}][debit]"]`).val(DRAmount.toFixed(2))
              
      if (!$(`select[name ="line-item[${i}][account]"]`).text()) {
        $('#saveJournalButton').prop('disabled', true);
      }
    } else if (document.getElementsByName(`line-item[${i}][debit]`)[0].value) {
      $('#saveJournalButton').prop('disabled', true);
    } else {
      var CRAmount = parseFloat(document.getElementsByName(`line-item[${i}][credit]`)[0].value);
      if (!isNaN(CRAmount) && CRAmount >=0) { 
        CRTotal += CRAmount; 
      $(`input[name ="line-item[${i}][credit]"]`).val(CRAmount.toFixed(2))
        if (!$(`select[name ="line-item[${i}][account]"]`).text()) {
          $('#saveJournalButton').prop('disabled', true);
        }
      } else if (document.getElementsByName(`line-item[${i}][credit]`)[0].value) {
        $('#saveJournalButton').prop('disabled', true);
      }
    }
  }

  if ((Math.abs(DRTotal - CRTotal) >= 0.01) && (DRTotal > 0)) {
    $('#saveJournalButton').prop('disabled', true);
  }

  document.getElementById('invoiceTotalDebit').value = DRTotal.toFixed(2);
  document.getElementById('invoiceTotalCredit').value = CRTotal.toFixed(2);
}

function clearJournalDateDescription() {
  $('input[name=date').val('');
  $('input[name=narration').val('');
}

function clearJournalLineItems() {
  var rows = $('#journal > tr');
  rows.each(function(idx, li) {
    var lineItem = $(li);
    lineItem.remove();
  });
}
