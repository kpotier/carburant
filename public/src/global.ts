import "@maptiler/sdk/dist/maptiler-sdk.css";
import * as maptilersdk from "@maptiler/sdk";

maptilersdk.config.apiKey = "YOUR API KEY";
export const map = new maptilersdk.Map({
  container: "map",
  style: maptilersdk.MapStyle.STREETS,
  geolocate: maptilersdk.GeolocationType.POINT,
  geolocateControl: false,
  zoom: 13,
});

export const posDot = new maptilersdk.Marker({ draggable: true });

// List
export interface Result {
  id: number;
  coords: number[];
  distance: number;
  address_rd: string;
  address_cp: string;
  automate_2424: boolean;
  horaires: { Hour: number; Minutes: number }[][][];
  services: string[];
  gas: { [x: string]: ResultGas };
}

export interface ResultGas {
  date: string;
  amount: number;
}

export const properties = {
  gas: "E10",
  sort: "distance",
  favorites: <number[]>[],
  res: <Result[]>[],
  resFav: <Result[]>[],
};
