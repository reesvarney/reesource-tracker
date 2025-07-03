<script lang="ts">
  import { onMount } from "svelte";
  import { Button } from "$lib/components/ui/button";
  import * as Select from "$lib/components/ui/select";
  import * as Card from "$lib/components/ui/card";
  import {
    UUIDBlobToString,
    UUIDStringToBlob,
  } from "$lib/components/uuid_helper";

  let sampleId = "";
  let sample: any = $state(null);
  let mods: any[] = $state([]);

  // For dropdowns
  let products = $state([]);
  let locations = $state([]);

  async function fetchSample() {
    // Get sampleId from URL (query param: ?sample=...)
    const params = new URLSearchParams(window.location.search);
    sampleId = params.get("sample") || "";
    if (!sampleId) {
      return;
    }
    // Fetch sample data
    const res = await fetch(`/api/sample/${sampleId}`);
    let data = null;
    if (res.ok) {
      data = await res.json();
      sample = data["sample"];
      mods = data["mods"] ?? [];
    } else {
      console.log("!", sampleId, UUIDStringToBlob(sampleId));
      sample = {
        ID: UUIDStringToBlob(sampleId),
        ProductID: "",
        LocationID: "",
        State: "",
      };
      mods = [];
    }
    // Fetch dropdown data
    products = await (await fetch("/api/products")).json();
    locations = await (await fetch("/api/locations")).json();
  }

  async function saveSample() {
    const form = new FormData();
    form.append("location_id", sample.LocationID);
    form.append("product_id", sample.ProductID);
    form.append("state", sample.State);
    const res = await fetch(`/api/sample/${sampleId}`, {
      method: "POST",
      body: form,
    });
    if (!res.ok) throw new Error(await res.text());
    // Optionally, show a success message or redirect
    await fetchSample();
  }
  $inspect(sample);
  $inspect(products);
  $inspect(locations);

  onMount(fetchSample);
</script>

<Card.Header>
  <Card.Title>Edit Sample</Card.Title>
</Card.Header>
<Card.Content>
  {#if sample}
    <form on:submit|preventDefault={saveSample} class="space-y-4">
      <div>
        <label class="block font-bold" for="sample-id">Sample ID</label>
        <div id="sample-id">{UUIDBlobToString(sample.ID)}</div>
      </div>
      <div>
        <label class="block font-bold" for="variant-select"
          >Product Variant</label
        >
        <Select.Root
          value={sample.VariantID}
          onValueChange={(e) => (sample.VariantID = e)}
          type="single"
        >
          <Select.Trigger id="variant-select" class="w-full">
            {sample.VariantID !== ""
              ? products.find((v) => v["VariantID"] == sample.VariantID)?.[
                  "CombinedName"
                ]
              : "Select a product variant"}
          </Select.Trigger>
          <Select.Content class="w-full">
            {#each products as variant}
              <Select.Item value={variant["VariantID"]}>
                {variant["CombinedName"]}
              </Select.Item>
            {/each}
          </Select.Content>
        </Select.Root>
      </div>
      <div>
        <label class="block font-bold" for="location-select">Location</label>
        <Select.Root type="single" bind:value={sample.LocationID}>
          <Select.Trigger id="location-select" class="w-full">
            {sample.LocationID !== ""
              ? locations.find((l) => {
                  console.log(l);
                  return l["ID"] == sample.LocationID;
                })?.["Name"]
              : "Select a location"}
          </Select.Trigger>
          <Select.Content class="w-full">
            {#each locations as location}
              <Select.Item value={location["ID"]}
                >{location["Name"]}</Select.Item
              >
            {/each}
          </Select.Content>
        </Select.Root>
      </div>
      <div>
        <label class="block font-bold" for="state-input">State</label>
        <Select.Root bind:value={sample.State} type="single">
          <Select.Trigger id="state-input" class="w-full">
            {sample.State
              ? sample.State === "in_use"
                ? "In Use"
                : sample.State === "available"
                  ? "Available"
                  : sample.State === "destroyed"
                    ? "Destroyed"
                    : sample.State
              : "Select a state"}
          </Select.Trigger>
          <Select.Content class="w-full">
            <Select.Item value="in_use">In Use</Select.Item>
            <Select.Item value="available">Available</Select.Item>
            <Select.Item value="destroyed">Destroyed</Select.Item>
          </Select.Content>
        </Select.Root>
      </div>
      <div>
        <Button type="submit">Save</Button>
      </div>
    </form>
    <div class="mt-6">
      <h3 class="font-bold mb-2">Mods</h3>
      {#if mods.length === 0}
        <div>No mods for this sample.</div>
      {:else}
        <ul class="list-disc ml-6">
          {#each mods as mod}
            <li>{mod.ModName} ({mod.ModType})</li>
          {/each}
        </ul>
      {/if}
    </div>
  {/if}
</Card.Content>
