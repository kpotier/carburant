import { properties } from "./global";
import { getList } from "./sort";

// Favorites
const favLS = localStorage.getItem("favorites");
if (favLS !== null) {
  const l = favLS.split(",");
  l.forEach((e) => {
    properties.favorites.push(Number(e));
  });
}

// Gas selection
const selPopup = <HTMLElement>document.getElementById("selection-popup");
const sel = <HTMLElement>document.getElementById("selection");
const selForm = <HTMLFormElement>document.getElementById("selection-form");

const gasLS = localStorage.getItem("gas");
if (gasLS !== null) {
  properties.gas = gasLS;
}
selPopup.innerText = properties.gas + " ▶";
if (gasLS !== null) sel.style.display = "none";
else sel.style.display = "flex";

selForm.onsubmit = (e: Event) => {
  e.preventDefault();
  const form = new FormData(selForm);
  const g = form.get("gas")?.toString();
  if (g !== undefined) {
    localStorage.setItem("gas", g);
    properties.gas = g;
    sel.style.display = "none";
    selPopup.innerText = properties.gas + " ▶";
  }
  getList();
  return false;
};

selPopup.onclick = (e: Event) => {
  e.preventDefault();
  sel.style.display = "flex";
  return false;
};
