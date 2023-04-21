import { map, posDot, properties } from "./global";
import { toSort } from "./sort";

const favPopup = <HTMLElement>document.getElementById("fav-popup");
const fav = <HTMLElement>document.getElementById("fav");
const favItems = <HTMLElement>document.getElementById("items-fav");
const favClose = <HTMLElement>document.getElementById("fav-close");

map.on("load", () => {
  favPopup.onclick = async () => {
    fav.style.visibility = "visible";

    const lngLat = posDot.getLngLat();
    const f = await fetch("./api/favorites", {
      method: "POST",
      headers: {
        "Content-Type": "application/x-www-form-urlencoded",
      },
      body:
        "lng=" +
        lngLat.lng +
        "&lat=" +
        lngLat.lat +
        "&list=" +
        localStorage.getItem("favorites"),
    });

    if (f.headers.get("content-type") != "application/json; charset=utf-8") {
      alert("error while fetching the service stations");
      return;
    }

    const results = await f.json();
    if (results["error"]) {
      alert("error while fetching the service stations");
      console.log(results["error"]);
      return;
    }

    properties.resFav = results;
    toSort(favItems, results, false);
  };
});

favClose.onclick = () => {
  fav.style.visibility = "";
};
