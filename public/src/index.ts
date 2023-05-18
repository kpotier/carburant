import "./minireset.css";
import "./index.css";
import { map, posDot, properties } from "./global";
import { displayList, getList } from "./list";
import "./gas";
import "./fav";
import "./sort";

// Position marker
const list = <HTMLElement>document.getElementById("list");
const listClose = <HTMLElement>document.getElementById("list-close");
const listPopup = <HTMLElement>document.getElementById("list-popup");

// Items
const items = <HTMLElement>document.getElementById("items");

listClose.onclick = () => {
  list.style.visibility = "";
};

listPopup.onclick = () => {
  list.style.visibility = "visible";
};

map.on("load", function () {
  // Geolocation
  posDot.setLngLat(map.getCenter());
  posDot.addTo(map);

  async function getListFromDot() {
    const lngLat = posDot.getLngLat();
    const res = await getList(lngLat.lng, lngLat.lat, properties.gas);
    properties.res = res;
    displayList(
      items,
      res,
      properties.favorites,
      properties.gas,
      properties.sort,
      map
    );
  }

  async function setCenter(pos: GeolocationPosition) {
    map.easeTo({
      center: { lng: pos.coords.longitude, lat: pos.coords.latitude },
    });
    posDot.setLngLat([pos.coords.longitude, pos.coords.latitude]);
    getListFromDot();
  }

  posDot.on("dragend", getListFromDot);
  navigator.geolocation.getCurrentPosition(setCenter, getListFromDot, {
    enableHighAccuracy: true,
    timeout: 5000,
    maximumAge: 0,
  });

  // Layers
  map.addSource("search-results", {
    type: "geojson",
    data: {
      type: "FeatureCollection",
      features: [],
    },
  });

  map.addLayer({
    id: "point-result",
    type: "circle",
    source: "search-results",
    paint: {
      "circle-radius": 15,
      "circle-color": "#51bbd6",
      "circle-opacity": 0.7,
    },
    filter: ["==", "$type", "Point"],
  });

  map.addLayer({
    id: "index",
    type: "symbol",
    source: "search-results",
    layout: {
      "text-field": "{index}",
      "text-font": ["DIN Offc Pro Medium", "Arial Unicode MS Bold"],
      "text-size": 12,
    },
  });

  map.on("mouseenter", "point-result", function () {
    map.getCanvas().style.cursor = "pointer";
  });
  map.on("mouseleave", "point-result", function () {
    map.getCanvas().style.cursor = "";
  });

  map.on("click", "point-result", function (e) {
    if (e.features === undefined || e.features[0].properties === null) return;
    const id = e.features[0].properties["index"];
    const elem = document.getElementById("item-" + id);
    list.scrollTo({ top: elem?.offsetTop, behavior: "smooth" });
    list.style.visibility = "visible";
  });
});

// Moreinfo popup
const moreInfo = <HTMLElement>document.getElementById("more-info");
const mInfoClose = <HTMLElement>document.getElementById("more-info-close");
mInfoClose.onclick = () => {
  moreInfo.style.display = "none";
};
