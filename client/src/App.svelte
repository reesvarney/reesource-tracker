<script lang="ts">
  import SampleList from "$views/sample_list.svelte";
  import SampleEditor from "$views/sample_editor.svelte";
  import ProductEditor from "$views/product_editor.svelte";
  import LocationEditor from "$views/location_editor.svelte";
  import FindSample from "$views/find_sample.svelte";
  import BulkApply from "$views/bulk_apply.svelte";
  import { Button } from "$lib/components/ui/button";
  import SampleCodeGenerator from "$views/sample_code_generator.svelte";
  import * as Tabs from "$lib/components/ui/tabs";

  import { onMount } from "svelte";
  import * as Card from "$lib/components/ui/card";
  import { UpdateAppStore } from "$lib/components/app_store";
  import { Toaster } from "$lib/components/ui/sonner/index.js";

  let currentPage = $state("quick_actions");
  onMount(() => {
    if (window.location.search.includes("sample")) {
      currentPage = "sample_edit";
    } else if (window.location.search.includes("product")) {
      currentPage = "product_edit";
    } else {
      currentPage = "quick_actions";
    }

    UpdateAppStore();
  });

  $effect(() => {
    if (currentPage !== "sample_edit") {
      const url = new URL(window.location.href);
      url.searchParams.delete("sample");
      window.history.replaceState({}, "", url.toString());
    }
  });

  let bulk_apply_active = $derived(currentPage === "bulk_apply");
  let find_sample_active = $derived(currentPage === "find_sample");
  $inspect(bulk_apply_active, find_sample_active);
</script>

<Toaster />

<main class="w-full overflow-hidden p-6 flex flex-col justify-stretch">
  <Card.Root class="grow max-h-full overflow-y-auto pb-4">
    <Card.Content class="grow max-h-full flex flex-col justify-stretch h-full">
      {#if currentPage === "quick_actions"}
        <div
          class="flex flex-col gap-4 items-stretch w-full h-full justify-center p-2"
        >
          <Button
            size="lg"
            class="text-lg py-6"
            onclick={() => (currentPage = "find_sample")}>Find Sample</Button
          >
          <Button
            size="lg"
            class="text-lg py-6"
            onclick={() => (currentPage = "bulk_apply")}>Bulk Apply</Button
          >
          <Button
            size="lg"
            class="text-lg py-6"
            onclick={() => (currentPage = "sample_list")}>Sample List</Button
          >
          <Button
            size="lg"
            class="text-lg py-6"
            onclick={() => (currentPage = "sample_code_generator")}
            >Provision Sample Codes</Button
          >
          <Button
            size="lg"
            class="text-lg py-6"
            onclick={() => (currentPage = "product_edit")}>Products</Button
          >
          <Button
            size="lg"
            class="text-lg py-6"
            onclick={() => (currentPage = "location_edit")}>Locations</Button
          >
        </div>
      {:else}
        <Tabs.Root bind:value={currentPage} class="w-full grow flex flex-col">
          <Tabs.List>
            <Tabs.Trigger value="find_sample" class="w-full"
              >Find Sample</Tabs.Trigger
            >
            <Tabs.Trigger value="bulk_apply" class="w-full"
              >Bulk Apply</Tabs.Trigger
            >
            <Tabs.Trigger value="sample_list">Sample List</Tabs.Trigger>
            <Tabs.Trigger value="sample_code_generator"
              >Provision Sample Codes</Tabs.Trigger
            >
            <Tabs.Trigger value="product_edit">Products</Tabs.Trigger>
            <Tabs.Trigger value="location_edit">Locations</Tabs.Trigger>
            {#if currentPage === "sample_edit"}
              <Tabs.Trigger value="sample_edit"
                >Sample {window.location.search.split(
                  "sample_id="
                )[1]}</Tabs.Trigger
              >
            {/if}
          </Tabs.List>
          <Tabs.Content value="find_sample">
            <FindSample bind:active={find_sample_active} />
          </Tabs.Content>
          <Tabs.Content value="bulk_apply">
            <BulkApply bind:active={bulk_apply_active} />
          </Tabs.Content>
          <Tabs.Content value="sample_list">
            <SampleList />
          </Tabs.Content>
          <Tabs.Content value="product_edit">
            <ProductEditor />
          </Tabs.Content>
          <Tabs.Content value="sample_edit">
            <SampleEditor />
          </Tabs.Content>
          <Tabs.Content value="location_edit">
            <LocationEditor />
          </Tabs.Content>
          <Tabs.Content value="sample_code_generator">
            <SampleCodeGenerator />
          </Tabs.Content>
        </Tabs.Root>
      {/if}
    </Card.Content>
  </Card.Root>
</main>

<style>
  main {
    @media print {
      max-height: none;
      overflow: visible;
      height: auto;
    }

    @media screen {
      max-height: 100vh;
      overflow: hidden;
      height: 100vh;
    }
  }
</style>
