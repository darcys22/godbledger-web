$(document).ready(function() {
    $('.js-example-basic-single').select2({
      theme: "bootstrap",
      placeholder: "Select Account",
      ajax: {
        url: '/api/reconcile/listunreconciledtransactions',
        dataType: 'json',
      }
    });
});
