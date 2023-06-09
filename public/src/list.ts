import Chart, { ChartItem } from "chart.js/auto";
import "chartjs-adapter-date-fns";
import { Result, ResultGas } from "./global";
import { Map } from "@maptiler/sdk";

export async function getList(
  lng: number,
  lat: number,
  gas: string
): Promise<Result[]> {
  const f = await fetch("./api/stations", {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: "lng=" + lng + "&lat=" + lat + "&gas=" + gas + "&lim=15",
  });

  if (f.headers.get("content-type") != "application/json; charset=utf-8") {
    alert("error while fetching the service stations");
    return [];
  }

  const results = await f.json();
  if (results["error"]) {
    alert("error while fetching the service stations");
    console.log(results["error"]);
    return [];
  }

  return results;
}

export function displayList(
  where: HTMLElement,
  res: Result[],
  favorites: number[],
  gas: string,
  sort: string,
  map: Map,
  addBool = true
) {
  if (sort === "price") {
    res.sort((a: Result, b: Result) => {
      if (a.gas[gas] && b.gas[gas])
        if (a.gas[gas].amount > b.gas[gas].amount) return 1;
        else return -1;
      else return -1;
    });
  } else {
    res.sort((a: Result, b: Result) => {
      if (a.distance > b.distance) return 1;
      else return -1;
    });
  }

  const data = {
    type: "FeatureCollection",
    features: <unknown[]>[],
  };

  where.innerHTML = "";

  for (let i = 0; i < res.length; i++) {
    const r = res[i];
    data.features.push({
      type: "Feature",
      geometry: {
        type: "Point",
        coordinates: [r.coords[1], r.coords[0]],
      },
      properties: {
        index: i + 1,
      },
    });

    const node = document.createElement("div");
    node.id = "item-" + (i + 1);
    node.className = "item";
    node.innerHTML = `<div class="item-info">
      <div>${i + 1}</div>
      <div>
        <div>${r.distance.toFixed(
          1
        )} km</div><div><button class="item-favorite">${
      favorites.includes(r.id) ? "★" : "☆"
    }</button></div></div></div><div class="item-address">${r.address_rd}<br/>${
      r.address_cp
    }</div><div class="item-price">${
      r.gas[gas] ? (r.gas[gas].amount / 1000).toFixed(3) : "not available"
    } €/L</div>`;

    const button = <HTMLButtonElement>node.getElementsByTagName("button")[0];
    button.onclick = () => {
      if (favorites.includes(r.id)) {
        const idx = favorites.indexOf(r.id);
        favorites.splice(idx, 1);
        button.innerText = "☆";
      } else {
        favorites.push(r.id);
        button.innerText = "★";
      }
      localStorage.setItem("favorites", favorites.join(","));
    };

    const buttonGo = document.createElement("button");
    buttonGo.innerText = "Go!";
    buttonGo.onclick = () => {
      window.open(
        "https://maps.google.com/?q=" + r.coords[0] + "," + r.coords[1],
        "_blank"
      );
    };
    node.appendChild(buttonGo);

    const buttonSeeMore = document.createElement("button");
    buttonSeeMore.innerText = "See more";
    buttonSeeMore.onclick = () => {
      getHistory(r, gas);
    };
    node.appendChild(buttonSeeMore);

    where.appendChild(node);
  }

  const add = map.getSource("search-results");
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  if (add !== undefined && addBool) (add as any).setData(data);
}

async function getHistory(r: Result, gas: string) {
  const f = await fetch("./api/history", {
    method: "POST",
    headers: {
      "Content-Type": "application/x-www-form-urlencoded",
    },
    body: "id=" + r.id + "&gas=" + gas,
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

  const horaires: string[] = [];
  function pad(num: number, size: number) {
    let num2 = num.toString();
    while (num2.length < size) num2 = "0" + num2;
    return num2;
  }
  r.horaires.forEach((elem) => {
    if (elem === null || elem.length == 0) {
      horaires.push("closed");
    } else {
      let builder = "";
      for (let i = 0; i < elem.length; i++) {
        if (i != 0) {
          builder += ", ";
        }
        builder +=
          pad(elem[i][0].Hour, 2) +
          ":" +
          pad(elem[i][0].Minutes, 2) +
          " to " +
          pad(elem[i][1].Hour, 2) +
          ":" +
          pad(elem[i][1].Minutes, 2);
      }
      horaires.push(builder);
    }
  });

  const moreInfo = <HTMLElement>document.getElementById("more-info");
  const mInfoCont = <HTMLElement>document.getElementById("more-info-container");

  moreInfo.style.display = "flex";
  mInfoCont.innerHTML =
    r.address_rd +
    "<br />" +
    r.address_cp +
    "<br /><br />" +
    "Automate 24/24: " +
    r.automate_2424 +
    "<br />Monday: " +
    horaires[0] +
    "<br />Tuesday: " +
    horaires[1] +
    "<br />Wednesday: " +
    horaires[2] +
    "<br />Thursday: " +
    horaires[3] +
    "<br />Friday: " +
    horaires[4] +
    "<br />Saturday: " +
    horaires[5] +
    "<br />Sunday: " +
    horaires[6] +
    "<br /><br />" +
    "Services: " +
    r.services.join(", ");

  const canvasHolder = document.createElement("div");
  canvasHolder.style.height = "200px";
  canvasHolder.style.display = "flex";
  canvasHolder.style.justifyContent = "center";

  const labels: string[] = [];
  const data: { x: string; y: number }[] = [];
  (results as ResultGas[]).forEach((element) => {
    data.push({
      y: element.amount / 1000,
      x: element.date,
    });
    labels.push(element.date);
  });

  const canvas = <HTMLCanvasElement>document.createElement("canvas");
  new Chart(canvas as ChartItem, {
    type: "line",
    data: {
      labels: labels,
      datasets: [
        {
          data: data,
          label: gas,
        },
      ],
    },
    options: {
      scales: {
        x: {
          type: "time",
        },
      },
    },
  });
  canvasHolder.appendChild(canvas);
  mInfoCont.appendChild(canvasHolder);
}
