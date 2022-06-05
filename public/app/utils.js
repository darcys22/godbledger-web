var truncate = function (fullStr, strLen, separator) {
    if (fullStr.length <= strLen) return fullStr;

    separator = separator || '...';

    var sepLen = separator.length,
        charsToShow = strLen - sepLen,
        frontChars = Math.ceil(charsToShow/2),
        backChars = Math.floor(charsToShow/2);

    return fullStr.substr(0, frontChars) +
           separator +
           fullStr.substr(fullStr.length - backChars);
};

function formatdate(element) {
  return moment(element).format('Do MMMM YYYY');
}
function formatformaldate(element) {
  return moment(element).format('YYYY-MM-DD');
}
function moneyNumber(x, decimals = 0) {
  xstr = x.toString();
  truncstr = xstr.substring(0, xstr.length - decimals);
  truncstrcomma = truncstr.replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  if (decimals > 0) {
      truncstrcomma = truncstrcomma + "."+xstr.substring(xstr.length - decimals, xstr.length);
  }
  return truncstrcomma;
}

const copyToClipboard = str => {
  const el = document.createElement('textarea');
  el.value = str;
  el.setAttribute('readonly', '');
  el.style.position = 'absolute';
  el.style.left = '-9999px';
  document.body.appendChild(el);
  const selected =
    document.getSelection().rangeCount > 0
      ? document.getSelection().getRangeAt(0)
      : false;
  el.select();
  document.execCommand('copy');
  document.body.removeChild(el);
  if (selected) {
    document.getSelection().removeAllRanges();
    document.getSelection().addRange(selected);
  }
};

function download(filename, text) {
  var element = document.createElement('a');
  element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(text));
  element.setAttribute('download', filename);

  element.style.display = 'none';
  document.body.appendChild(element);

  element.click();

  document.body.removeChild(element);
}
