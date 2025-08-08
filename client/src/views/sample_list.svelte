<script lang="ts">
  import * as Table from "$lib/components/ui/table";
  import { Button } from "$lib/components/ui/button";
  import ProductSelect from "$lib/components/selects/product-select.svelte";
  import LocationSelect from "$lib/components/selects/location-select.svelte";
  import StateSelect from "$lib/components/selects/state-select.svelte";
  import { AppStore } from "$lib/components/app_store";
  import { Checkbox } from "$lib/components/ui/checkbox";
  import * as Pagination from "$lib/components/ui/pagination";
  import * as Card from "$lib/components/ui/card";
  import { Separator } from "$lib/components/ui/separator";
  let selectedProduct = $state("");
  let selectedLocation = $state("");
  let page = $state(1);
  import { onMount, onDestroy } from "svelte";
  import { Input } from "$lib/components/ui/input";
  import UserSelect from "$lib/components/selects/user-select.svelte";
  import { Label } from "$lib/components/ui/label";
  let pageSize = $state(1);
  let tableContainer: HTMLDivElement | null = $state(null);
  let firstRow: HTMLElement | null = $state(null);
  let paginationContainer: HTMLDivElement | null = $state(null);

  let showUnassigned = $state(false);
  let modQuery = $state("");
  let selectedState = $state("");
  let ownerQuery = $state(""); // Will hold user id
  let issueQuery = $state("");

  function calculatePageSize() {
    if (!tableContainer) {
      console.warn("tableContainer not available, cannot calculate page size");
      setTimeout(() => {
        // Fallback in case firstRow is not immediately available
        calculatePageSize();
      }, 1000);
      return;
    }
    tableContainer.style.height = "100%";

    if (!firstRow) {
      // Fallback if firstRow is not available
      pageSize = 1;
      console.warn("firstRow not available, using fallback pageSize of 1");
      setTimeout(() => {
        // Fallback in case firstRow is not immediately available
        calculatePageSize();
      }, 1000);
      return;
    }

    let ROW_HEIGHT = 9999;
    const measured = firstRow.clientHeight;
    if (measured && measured > 0) {
      ROW_HEIGHT = measured;
    }
    const totalHeight =
      tableContainer.clientHeight - (paginationContainer?.clientHeight || 0) ||
      0;
    const rows = Math.max(1, Math.floor(totalHeight / ROW_HEIGHT));
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
    const directChildren = $AppStore.products.filter(
      (p) => p.parentProductID === parentId
    );
    let allIds: string[] = directChildren.map((p) => p.id);
    for (const child of directChildren) {
      allIds = allIds.concat(getDescendantProductIds(child.id));
    }
    return allIds;
  }
  function getDescendantLocationIds(parentId: string): string[] {
    const directChildren = $AppStore.locations.filter(
      (l) => l.parentLocationID === parentId
    );
    let allIds: string[] = directChildren.map((l) => l.id);
    for (const child of directChildren) {
      allIds = allIds.concat(getDescendantLocationIds(child.id));
    }
    return allIds;
  }

  const filteredSamples = $derived(() => {
    let result = $AppStore.samples;
    if (!showUnassigned && selectedState !== "unassigned") {
      result = result.filter((sample) => {
        return sample.state !== "unassigned";
      });
    }
    if (selectedState && selectedState !== "") {
      result = result.filter((sample) => sample.state === selectedState);
    }
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
    if (modQuery != "") {
      const queries = modQuery
        .toLowerCase()
        .split(",")
        .map((q) => q.trim());
      result = result.filter((sample) => {
        if (sample.ModSummary === "No mods") return false;
        if (sample.ModSummary === "No active mods") return false;
        return queries.some((query) =>
          sample.ModSummary.toLowerCase().includes(query)
        );
      });
    }
    if (ownerQuery && ownerQuery !== "") {
      result = result.filter((sample) => sample.Owner?.id === ownerQuery);
    }
    if (issueQuery.trim() !== "") {
      const iq = issueQuery.toLowerCase();
      result = result.filter((sample) =>
        sample.productIssue?.toLowerCase().includes(iq)
      );
    }
    return result;
  });

  const totalPages = $derived(() =>
    Math.max(1, Math.ceil(filteredSamples().length / pageSize))
  );
  const paginatedSamples = $derived(() =>
    filteredSamples().slice((page - 1) * pageSize, page * pageSize)
  );

  $effect(() => {
    selectedLocation ||
      selectedProduct ||
      modQuery ||
      selectedState ||
      showUnassigned;

    page = 1; // Reset to first page when filters change
  });

  function EditSample(sampleId: string) {
    window.location.href = `/app?sample_id=${sampleId}`;
  }
</script>

<div class="max-h-full grow h-full overflow-auto">
  <div class="flex flex-col h-full max-h-full">
    <Card.Root class="grow">
      <Card.Header>
        <Card.Title>Sample List</Card.Title>
        <Card.Description>
          Find and manage samples in the system.
        </Card.Description>
        <div class=" flex flex-wrap items-start gap-12 flex-row w-full">
          <div class="flex flex-row gap-8">
            <div>
              <Label class=" mb-2" for="product-select">Product</Label>
              <ProductSelect
                bind:bindValue={selectedProduct}
                filterMode={true}
              />
            </div>

            <div>
              <Label class=" mb-2" for="issue-query">Issue</Label>
              <Input
                type="text"
                placeholder="Filter by issue"
                class="input input-bordered w-full max-w-xs"
                bind:value={issueQuery}
              />
            </div>
          </div>
          <div>
            <Label class=" mb-2" for="mod-query">Mods</Label>
            <Input
              type="text"
              placeholder="Filter by mods"
              class="input input-bordered w-full max-w-xs"
              bind:value={modQuery}
            />
          </div>

          <div>
            <Label class=" mb-2" for="location-select">Location</Label>
            <LocationSelect
              bind:bindValue={selectedLocation}
              filterMode={true}
            />
          </div>
          <div>
            <Label class=" mb-2" for="owner-query">Owner</Label>
            <UserSelect bind:bindValue={ownerQuery} filterMode={true} />
          </div>

          <div>
            <Label class=" mb-2" for="state-select">State</Label>
            <div class="flex items-center self-end gap-2">
              <StateSelect
                bind:bindValue={selectedState}
                filterMode={true}
                showUnassigned={true}
              />
              {#if selectedState === ""}
                <Checkbox id="show-unassigned" bind:checked={showUnassigned} />
                <Label for="show-unassigned" class="no-wrap min-w-max">
                  Include Unassigned Samples
                </Label>
              {/if}
            </div>
          </div>

          <div class="self-end align-self-end grow flex flex-row justify-end">
            <Button
              variant="outline"
              onclick={() => {
                $AppStore.currentPage = "find_sample";
                $AppStore = $AppStore;
              }}>Find Sample by ID or QR Code</Button
            >
          </div>
        </div>
      </Card.Header>
      <Separator />
      <Card.Content class="h-full flex flex-col">
        <div
          bind:this={tableContainer}
          class="grow basis-0 max-h-full h-full overflow-hidden"
        >
          <Table.Root>
            <Table.Header>
              <Table.Row>
                <Table.Head>ID</Table.Head>
                <Table.Head>Product</Table.Head>
                <Table.Head>Issue</Table.Head>
                <Table.Head>Mods</Table.Head>
                <Table.Head>Location</Table.Head>
                <Table.Head>Owner</Table.Head>
                <Table.Head>Status</Table.Head>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {#each paginatedSamples() as sample, i}
                <Table.Row
                  bind:ref={
                    () => firstRow, (el) => (i === 0 ? (firstRow = el) : null)
                  }
                >
                  <Table.Cell>{sample.DisplayId}</Table.Cell>
                  <Table.Cell
                    >{sample.Product?.CombinedName}
                    {#if sample.Product?.partNumber}{`(${sample.Product?.partNumber})`}{/if}</Table.Cell
                  >
                  <Table.Cell>{sample.productIssue}</Table.Cell>
                  <Table.Cell class="max-w-32 overflow-hidden text-ellipsis">
                    {sample.ModSummary}
                  </Table.Cell>
                  <Table.Cell>{sample.Location?.name}</Table.Cell>
                  <Table.Cell>{sample.Owner?.name}</Table.Cell>
                  <Table.Cell>{sample.Status}</Table.Cell>
                  <Table.Cell>
                    <Button onclick={() => EditSample(sample.DisplayId)}>
                      Edit
                    </Button>
                  </Table.Cell>
                </Table.Row>
              {/each}
            </Table.Body>
          </Table.Root>
        </div>
        <!-- Pagination Controls -->
        <div class="flex justify-center mt-6" bind:this={paginationContainer}>
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
      </Card.Content>
    </Card.Root>
  </div>
</div>
