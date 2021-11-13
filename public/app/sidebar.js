let sidebar = document.querySelector(".sidebar");
let closeBtn = document.querySelector("#btn");
let searchBtn = document.querySelector(".bx-search");


closeBtn.addEventListener("click", ()=>{
  sidebar.classList.toggle("open");
  menuBtnChange();//calling the function(optional)
});

//searchBtn.addEventListener("click", ()=>{ // Sidebar open when you click on the search iocn
  //sidebar.classList.toggle("open");
  //menuBtnChange(); //calling the function(optional)
//});

// following are the code to change sidebar button(optional)
function menuBtnChange() {
 if(sidebar.classList.contains("open")){
   closeBtn.classList.replace("bx-menu", "bx-menu-alt-right");//replacing the iocns class
 }else {
   closeBtn.classList.replace("bx-menu-alt-right","bx-menu");//replacing the iocns class
 }
}

$('.adminonly').each(function(i, obj) {
  obj.hidden = true;
})

fetch('/api/user/settings', {
  method: 'GET',
  headers: {
    'Content-Type': 'application/json;charset=utf-8'
  }
})
.then(response => response.json())
.then(data => {
  let name = document.querySelector(".name");
  let role = document.querySelector(".job");
  window.user = data
  name.innerText = data.name[0].toUpperCase() + data.name.slice(1);
  role.innerText = data.role[0].toUpperCase() + data.role.slice(1);
  if (typeof Date.setLocale !== 'undefined') {
    Date.setLocale(data.datelocale);
  }

  if (data.role == "admin") {
    $('.adminonly').each(function(i, obj) {
      obj.hidden = false;
    })
  }
})

