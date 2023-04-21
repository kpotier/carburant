import "./minireset.css";
import "./index.css";
import { map, posDot } from "./global";
import { getList } from "./sort";
import "./gas";
import "./fav";

// Position marker
const list = <HTMLElement>document.getElementById("list");
const listClose = <HTMLElement>document.getElementById("list-close");
const listPopup = <HTMLElement>document.getElementById("list-popup");

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

  function setCenter(pos: GeolocationPosition) {
    map.easeTo({
      center: { lng: pos.coords.longitude, lat: pos.coords.latitude },
    });
    posDot.setLngLat([pos.coords.longitude, pos.coords.latitude]);
    getList();
  }

  posDot.on("dragend", getList);
  navigator.geolocation.getCurrentPosition(setCenter, getList, {
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
