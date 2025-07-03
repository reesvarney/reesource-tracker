<script lang="ts">
  import SampleList from "$views/sample_list.svelte";
  import SampleEdit from "$views/sample_editor.svelte";
  import ProductEditor from "$views/product_editor.svelte";
  import * as Tabs from "$lib/components/ui/tabs";

  import { onMount } from "svelte";
  import * as Card from "$lib/components/ui/card";
  let currentPage = $state("sample_list");
  onMount(() => {
    // This is where you can initialize any global state or perform actions
    // when the app is mounted, if needed.
    if (window.location.search.includes("sample")) {
      currentPage = "sample_edit";
    } else {
      currentPage = "sample_list";
    }
  });

  $effect(() => {
    // Remove the sample query parameter from the URL if the sample_editor tab is not active
    if (currentPage !== "sample_edit") {
      const url = new URL(window.location.href);
      url.searchParams.delete("sample");
      window.history.replaceState({}, "", url.toString());
    }
  });
</script>

<main>
  <div class="w-full h-full p-2 flex flex-col gap-2">
    <Card.Root>
      <Card.Content>
        <Tabs.Root bind:value={currentPage} class="w-full">
          <Tabs.List>
            <Tabs.Trigger value="sample_list">Samples</Tabs.Trigger>
            <Tabs.Trigger value="product_edit">Products</Tabs.Trigger>
            {#if currentPage === "sample_edit"}
              <Tabs.Trigger value="sample_edit"
                >Sample {window.location.search.split(
                  "sample="
                )[1]}</Tabs.Trigger
              >
            {/if}
          </Tabs.List>
          <Tabs.Content value="sample_list">
            <SampleList />
          </Tabs.Content>
          <Tabs.Content value="product_edit">
            <ProductEditor />
          </Tabs.Content>
          <Tabs.Content value="sample_edit">
            <SampleEdit />
          </Tabs.Content>
        </Tabs.Root>
      </Card.Content>
    </Card.Root>
  </div>
</main>

<style>
</style>
