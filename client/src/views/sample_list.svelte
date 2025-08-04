<script lang="ts">
  import * as Table from "$lib/components/ui/table";
  import { Button } from "$lib/components/ui/button";
  import ProductSelect from "$lib/components/selects/product-select.svelte";
  import LocationSelect from "$lib/components/selects/location-select.svelte";
  import { AppStore } from "$lib/components/app_store";
  import * as Pagination from "$lib/components/ui/pagination";
  let selectedProduct = $state("");
  let selectedLocation = $state("");
  let page = $state(1);
  import { onMount, onDestroy } from "svelte";
  let pageSize = $state(20); // Will be dynamically set based on available space
  let tableContainer: HTMLDivElement;
  const ROW_HEIGHT = 48; // px, adjust if your row height is different
  const HEADER_HEIGHT = 56; // px, adjust if your table header is different
  const PAGINATION_HEIGHT = 64; // px, adjust if your pagination controls are different
  const FILTERS_HEIGHT = 64; // px, adjust if your filters area is different
  const MARGIN = 32; // px, extra margin for safety

  function calculatePageSize() {
    if (!tableContainer) return;
    const totalHeight = tableContainer.clientHeight;
    const available =
      totalHeight - HEADER_HEIGHT - PAGINATION_HEIGHT - FILTERS_HEIGHT - MARGIN;
    const rows = Math.max(1, Math.floor(available / ROW_HEIGHT));
    pageSize = rows;
  }

  let resizeObserver: ResizeObserver | null = null;

  onMount(() => {
    calculatePageSize();
    if (tableContainer) {
      resizeObserver = new ResizeObserver(() => {
        calculatePageSize();
      });
      resizeObserver.observe(tableContainer);
    }
    window.addEventListener("resize", calculatePageSize);
  });

  onDestroy(() => {
    if (resizeObserver && tableContainer) {
      resizeObserver.unobserve(tableContainer);
      resizeObserver.disconnect();
    }
    window.removeEventListener("resize", calculatePageSize);
  });

  // Helper to get all descendant ids for a given parent id (recursive)
  function getDescendantProductIds(parentId: string): string[] {
    const directChildren = $AppStore.products.filter((p) => p.parentProductID === parentId);
    let allIds: string[] = directChildren.map((p) => p.id);
    for (const child of directChildren) {
      allIds = allIds.concat(getDescendantProductIds(child.id));
    }
    return allIds;
  }
  function getDescendantLocationIds(parentId: string): string[] {
    const directChildren = $AppStore.locations.filter((l) => l.parentLocationID === parentId);
    let allIds: string[] = directChildren.map((l) => l.id);
    for (const child of directChildren) {
      allIds = allIds.concat(getDescendantLocationIds(child.id));
    }
    return allIds;
  }

  const filteredSamples = $derived(() => {
    let result = $AppStore.samples;
    if (selectedProduct) {
      if (selectedProduct.startsWith("any-")) {
        const parentId = selectedProduct.slice(4);
        const descendantIds = getDescendantProductIds(parentId);
        // Include both parent and all descendants
        const validIds = [parentId, ...descendantIds];
        result = result.filter(
          (sample) => sample.Product?.id && validIds.includes(sample.Product.id)
        );
      } else {
        result = result.filter(
          (sample) => sample.Product?.id === selectedProduct
        );
      }
    }
    if (selectedLocation) {
      if (selectedLocation.startsWith("any-")) {
        const parentId = selectedLocation.slice(4);
        const descendantIds = getDescendantLocationIds(parentId);
        // Include both parent and all descendants
        const validIds = [parentId, ...descendantIds];
        result = result.filter(
          (sample) =>
            sample.Location?.id && validIds.includes(sample.Location.id)
        );
      } else {
        result = result.filter(
          (sample) => sample.Location?.id === selectedLocation
        );
      }
    }
    return result;
  });

  const totalPages = $derived(() =>
    Math.max(1, Math.ceil(filteredSamples().length / pageSize))
  );
  const paginatedSamples = $derived(() =>
    filteredSamples().slice((page - 1) * pageSize, page * pageSize)
  );

  function EditSample(sampleId: string) {
    window.location.href = `/app?sample_id=${sampleId}`;
  }
</script>

<div class="flex flex-col h-full justify-stretch">
  <div
    class="mb-4 flex flex-wrap items-center gap-6 flex-row w-full"
    style="min-height:64px"
  >
    <div>
      <ProductSelect bind:bindValue={selectedProduct} filterMode={true} />
    </div>
    <div>
      <LocationSelect bind:bindValue={selectedLocation} filterMode={true} />
    </div>
  </div>
  <div class="grow" bind:this={tableContainer}>
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
        {#each paginatedSamples() as sample}
          <Table.Row>
            <Table.Cell>{sample.DisplayId}</Table.Cell>
            <Table.Cell>{sample.Product?.CombinedName}</Table.Cell>
            <Table.Cell>{sample.ModSummary}</Table.Cell>
            <Table.Cell>{sample.Location?.name}</Table.Cell>
            <Table.Cell>{sample.Status}</Table.Cell>
            <Table.Cell>
              <Button onclick={() => EditSample(sample.DisplayId)}>Edit</Button>
            </Table.Cell>
          </Table.Row>
        {/each}
      </Table.Body>
    </Table.Root>
  </div>
  <!-- Pagination Controls -->
  <div class="flex justify-center mt-6">
    <Pagination.Root
      count={filteredSamples().length}
      perPage={pageSize}
      bind:page
    >
      {#snippet children({ pages, currentPage })}
        <Pagination.Content>
          <Pagination.Item>
            <Pagination.PrevButton disabled={page === 1} />
          </Pagination.Item>
          {#each pages as pageObj (pageObj.key)}
            {#if pageObj.type === "ellipsis"}
              <Pagination.Item>
                <Pagination.Ellipsis />
              </Pagination.Item>
            {:else}
              <Pagination.Item>
                <Pagination.Link
                  page={pageObj}
                  isActive={currentPage === pageObj.value}
                >
                  {pageObj.value}
                </Pagination.Link>
              </Pagination.Item>
            {/if}
          {/each}
          <Pagination.Item>
            <Pagination.NextButton disabled={page === totalPages()} />
          </Pagination.Item>
        </Pagination.Content>
      {/snippet}
    </Pagination.Root>
  </div>
</div>
