import { map, properties } from "./global";
import { displayList } from "./list";

const sDistArrow = <HTMLElement>document.getElementById("sort-distance-arrow");
const sPricArrow = <HTMLElement>document.getElementById("sort-price-arrow");
const sFavPricArrow = <HTMLElement>(
  document.getElementById("sort-fav-price-arrow")
);
const sFavDistArrow = <HTMLElement>(
  document.getElementById("sort-fav-distance-arrow")
);

const favItems = <HTMLElement>document.getElementById("items-fav");
const items = <HTMLElement>document.getElementById("items");

let sort = "distance";
if (localStorage.getItem("sort") === "price") {
  sort = "price";
  sPricArrow.innerText = "▼";
  sFavPricArrow.innerText = "▼";
} else {
  sDistArrow.innerText = "▼";
  sFavDistArrow.innerText = "▼";
}
properties.sort = sort;

const sortDist = <HTMLButtonElement>document.getElementById("sort-distance");
const sortFavDist = <HTMLButtonElement>(
  document.getElementById("sort-fav-distance")
);

const sortPric = <HTMLButtonElement>document.getElementById("sort-price");
const sortFavPric = <HTMLButtonElement>(
  document.getElementById("sort-fav-price")
);

sortFavDist.onclick = (e) => {
  sortDistFn(e);
};

sortDist.onclick = (e: Event) => {
  sortDistFn(e);
};

function sortDistFn(e: Event) {
  e.preventDefault();
  localStorage.setItem("sort", "distance");
  sort = "distance";
  properties.sort = sort;
  displayList(
    items,
    properties.res,
    properties.favorites,
    properties.gas,
    sort,
    map
  );
  displayList(
    favItems,
    properties.resFav,
    properties.favorites,
    properties.gas,
    sort,
    map,
    false
  );
  sPricArrow.innerText = "";
  sFavPricArrow.innerText = "";
  sFavDistArrow.innerText = "▼";
  sDistArrow.innerText = "▼";
  return false;
}

function SortPricFn(e: Event) {
  e.preventDefault();
  localStorage.setItem("sort", "price");
  sort = "price";
  properties.sort = sort;
  displayList(
    items,
    properties.res,
    properties.favorites,
    properties.gas,
    sort,
    map
  );
  displayList(
    favItems,
    properties.resFav,
    properties.favorites,
    properties.gas,
    sort,
    map,
    false
  );
  sPricArrow.innerText = "▼";
  sFavPricArrow.innerText = "▼";
  sDistArrow.innerText = "";
  sFavDistArrow.innerText = "";
  return false;
}

sortPric.onclick = (e: Event) => {
  SortPricFn(e);
};

sortFavPric.onclick = (e: Event) => {
  SortPricFn(e);
};
