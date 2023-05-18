import { map, posDot, properties } from "./global";
import { displayList, getList } from "./list";

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

selPopup.onclick = (e: Event) => {
  e.preventDefault();
  sel.style.display = "flex";
  return false;
};

selForm.onsubmit = async (e: Event) => {
  e.preventDefault();
  const form = new FormData(selForm);
  const g = form.get("gas")?.toString();
  if (g !== undefined) {
    localStorage.setItem("gas", g);
    properties.gas = g;
    sel.style.display = "none";
    selPopup.innerText = properties.gas + " ▶";
  }
  const lngLat = posDot.getLngLat();
  const res = await getList(lngLat.lng, lngLat.lat, properties.gas);
  properties.res = res;
  displayList(
    <HTMLElement>document.getElementById("items"),
    res,
    properties.favorites,
    properties.gas,
    properties.sort,
    map,
    true
  );
  return false;
};
