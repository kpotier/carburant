import { Result, map, posDot, properties } from "./global";
import { displayList } from "./list";

const favPopup = <HTMLElement>document.getElementById("fav-popup");
const fav = <HTMLElement>document.getElementById("fav");
const favItems = <HTMLElement>document.getElementById("items-fav");
const favClose = <HTMLElement>document.getElementById("fav-close");

const favLS = localStorage.getItem("favorites");
if (favLS !== null) {
  const l = favLS.split(",");
  l.forEach((e) => {
    properties.favorites.push(Number(e));
  });
}

map.on("load", () => {
  favPopup.onclick = async () => {
    fav.style.visibility = "visible";
    const favorites = localStorage.getItem("favorites");

    let results: Result[] = [];
    if (favorites != "") {
      const lngLat = posDot.getLngLat();
      const f = await fetch("./api/favorites", {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded",
        },
        body: "lng=" + lngLat.lng + "&lat=" + lngLat.lat + "&list=" + favorites,
      });

      if (f.headers.get("content-type") != "application/json; charset=utf-8") {
        alert("error while fetching the service stations");
        return;
      }

      const res = await f.json();
      if (res["error"]) {
        alert("error while fetching the service stations");
        console.log(res["error"]);
        return;
      } else {
        results = res;
      }
    }

    properties.resFav = results;
    displayList(
      favItems,
      results,
      properties.favorites,
      properties.gas,
      properties.sort,
      map,
      false
    );
  };
});

favClose.onclick = () => {
  fav.style.visibility = "";
};
