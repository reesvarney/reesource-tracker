<script lang="ts">
  import * as Table from "$lib/components/ui/table";
  import { Button } from "$lib/components/ui/button";
  import * as Card from "$lib/components/ui/card";
  import { onMount } from "svelte";
  import { UUIDBlobToString } from "$lib/components/uuid_helper";

  let products = $state([]);
  let locations = $state([]);
  let samples = $state([]);
  let filter = $state("");
  let selectedProduct = $state("");
  let selectedLocation = $state("");
  const filteredSamples = $derived(() => {
    let result = samples;
    if (selectedProduct) {
      result = result.filter(
        (sample: any) => sample["VariantID"] === selectedProduct
      );
    }
    if (selectedLocation) {
      result = result.filter(
        (sample: any) => sample["LocationID"] === selectedLocation
      );
    }
    if (filter) {
      const lower = filter.toLowerCase();
      result = result.filter((sample: any) => {
        const productName =
          GetProductName(sample["VariantID"])?.toLowerCase() || "";
        const locationName =
          GetLocationName(sample["LocationID"])?.toLowerCase() || "";
        return (
          String(sample["ID"]).toLowerCase().includes(lower) ||
          productName.includes(lower) ||
          String(sample["CurrentModsSummary"]).toLowerCase().includes(lower) ||
          locationName.includes(lower)
        );
      });
    }
    return result;
  });

  async function GetSamples() {
    products = (await (await fetch("/api/products")).json()) ?? [];
    samples = (await (await fetch("/api/samples")).json()) ?? [];
    locations = (await (await fetch("/api/locations")).json()) ?? [];
  }

  function GetProductName(variant_id: string) {
    const found = products.find(
      (variant) => variant["VariantID"] === variant_id
    );
    if (found) {
      return `${found["CombinedName"]}`;
    } else {
      return "Unknown Product";
    }
  }

  function GetLocationName(location_id: string) {
    const found = locations.find((location) => location["ID"] === location_id);
    if (found) {
      return found["Name"];
    } else {
      return "Unknown Location";
    }
  }

  function EditSample(sampleId: string) {
    // Navigate to the sample edit page
    window.location.href = `/app?sample=${sampleId}`;
  }

  onMount(() => {
    GetSamples();
  });

  $inspect(locations);
  $inspect(samples);
  $inspect(products);
</script>

<div class="mb-4 flex flex-wrap items-center gap-2">
  <!-- Product Variant Dropdown -->
  <select
    class="select select-bordered max-w-xs"
    bind:value={selectedProduct}
    aria-label="Filter by product"
  >
    <option value="">All Products</option>
    {#each products as variant}
      <option value={variant["VariantID"]}>
        {variant["CombinedName"]} - Issue {variant["VariantPcbIssue"]?.[
          "String"
        ]}
      </option>
    {/each}
  </select>
  <!-- Location Dropdown -->
  <select
    class="select select-bordered max-w-xs"
    bind:value={selectedLocation}
    aria-label="Filter by location"
  >
    <option value="">All Locations</option>
    {#each locations as location}
      <option value={location["ID"]}>{location["Name"]}</option>
    {/each}
  </select>
</div>
<Table.Root>
  <Table.Header>
    <Table.Row>
      <Table.Head>ID</Table.Head>
      <Table.Head>Product</Table.Head>
      <Table.Head>Mods</Table.Head>
      <Table.Head>Location</Table.Head>
      <Table.Head>Status</Table.Head>
    </Table.Row>
  </Table.Header>
  <Table.Body>
    {#each filteredSamples() as sample}
      <Table.Row>
        <Table.Cell>{UUIDBlobToString(sample["ID"])}</Table.Cell>
        <Table.Cell>{GetProductName(sample["VariantID"])}</Table.Cell>
        <Table.Cell>{sample["CurrentModsSummary"]}</Table.Cell>
        <Table.Cell>{GetLocationName(sample["LocationID"])}</Table.Cell>
        <Table.Cell>{sample["State"]}</Table.Cell>
        <Table.Cell
          ><Button onclick={() => EditSample(UUIDBlobToString(sample["ID"]))}
            >Edit</Button
          ></Table.Cell
        >
      </Table.Row>
    {/each}
  </Table.Body>
</Table.Root>
