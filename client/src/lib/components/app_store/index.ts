import { writable } from "svelte/store";
import type { Writable } from "svelte/store";
import { Sample } from "$lib/components/sample";
import { SampleProduct } from "$lib/components/product";
import { SampleLocation } from "$lib/components/location";

export type AppData = {
  samples: Sample[];
  locations: SampleLocation[];
  products: SampleProduct[];
};

const app_data_cache = localStorage.getItem("app_data_cache");
let cached_data: AppData = {
  samples: [],
  locations: [],
  products: [],
};

if (app_data_cache) {
  const data = JSON.parse(app_data_cache);
  cached_data = data;
}

export const AppStore: Writable<AppData> = writable(cached_data);

export async function UpdateLocations() {
  const locations_req = await fetch("/api/locations");
  let new_locations: SampleLocation[] = [];
  if (!locations_req.ok) {
    console.error("Failed to fetch locations:", await locations_req.json());
  } else {
    new_locations = (await locations_req.json()) ?? [];
  }
  AppStore.update((data) => ({
    ...data,
    locations: new_locations.map(
      (location: any) => new SampleLocation(location, AppStore)
    ),
  }));
}

export async function UpdateProducts() {
  const products_req = await fetch("/api/products");
  let new_products: SampleProduct[] = [];
  if (!products_req.ok) {
    console.error("Failed to fetch products:", await products_req.json());
  } else {
    new_products = (await products_req.json()) ?? [];
  }
  AppStore.update((data) => ({
    ...data,
    products: new_products.map(
      (product: any) => new SampleProduct(product, AppStore)
    ),
  }));
}

export async function UpdateSamples() {
  const samples_req = await fetch("/api/samples");
  let new_samples: Sample[] = [];
  if (!samples_req.ok) {
    console.error("Failed to fetch samples:", await samples_req.json());
  } else {
    new_samples = (await samples_req.json()) ?? [];
  }
  AppStore.update((data) => ({
    ...data,
    samples: new_samples.map((sample: any) => new Sample(sample, AppStore)),
  }));
}

export async function UpdateAppStore() {
  await Promise.all([UpdateLocations(), UpdateProducts(), UpdateSamples()]);
}

// Connect to sync eventstream API and update store on events
function connectSyncEvents() {
  const evtSource = new EventSource("/api/sync");
  evtSource.addEventListener("samples_updated", () => {
    console.log("Samples updated event received");
    UpdateSamples();
  });
  evtSource.addEventListener("products_updated", () => {
    console.log("Products updated event received");
    UpdateProducts();
  });
  evtSource.addEventListener("locations_updated", () => {
    console.log("Locations updated event received");
    UpdateLocations();
  });
  evtSource.onerror = (err) => {
    console.error("Sync eventstream error", err);
    // Optionally, try to reconnect after a delay
  };
}

// Start sync connection immediately
connectSyncEvents();
