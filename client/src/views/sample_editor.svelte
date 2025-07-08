<script lang="ts">
  import { onMount } from "svelte";
  import { Button } from "$lib/components/ui/button";
  import ProductSelect from "$lib/components/selects/product-select.svelte";
  import LocationSelect from "$lib/components/selects/location-select.svelte";
  import StateSelect from "$lib/components/selects/state-select.svelte";
  import * as Card from "$lib/components/ui/card";
  import { IDStringToBlob } from "$lib/components/id_helper";
  import { AppStore } from "$lib/components/app_store";
  import { Sample } from "$lib/components/sample";
  import { SampleMod } from "$lib/components/sample_mod";
  import * as Table from "$lib/components/ui/table";
  let sample: Sample = $state(new Sample({}));
  let showFittedOnly = $state(false);
  let sample_id = "";
  let product_id = $state("");
  let location_id = $state("");
  let sample_state = $state("");
  let new_mod_name = $state("");
  let adding_mod = $state(false);
  let add_mod_error = $state("");

  async function fetchSample() {
    const params = new URLSearchParams(window.location.search);
    const new_sample_id = params.get("sample_id") || "";
    if (!new_sample_id) {
      return;
    }
    // Fetch sample data
    const res = await fetch(`/api/sample/${new_sample_id}`);
    if (res.ok) {
      const json_data = await res.json();
      sample = new Sample(
        { ...json_data.sample, mods: json_data.mods },
        AppStore
      );
      sample_id = sample.DisplayId;
      product_id = sample.Product?.id || "";
      location_id = sample.Location?.id || "";
      sample_state = sample.state || "";
    } else {
      sample = new Sample(
        {
          ID: IDStringToBlob(new_sample_id),
          State: "available",
        },
        AppStore
      );
      sample_id = new_sample_id;
    }
  }

  async function addMod() {
    if (!new_mod_name.trim()) return;
    adding_mod = true;
    add_mod_error = "";
    try {
      const res = await fetch(`/api/sample/${sample_id}/mods`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ name: new_mod_name }),
      });
      if (!res.ok) {
        add_mod_error = await res.text();
      } else {
        new_mod_name = "";
        await fetchSample();
      }
    } catch (e) {
      add_mod_error = "Failed to add mod.";
    }
    adding_mod = false;
  }

  async function removeMod(mod_id: string) {
    if (!mod_id) return;
    await fetch(`/api/sample/${sample_id}/mods/${mod_id}`, {
      method: "DELETE",
    });
    await fetchSample();
  }

  async function saveSample(e: SubmitEvent) {
    e.preventDefault();
    const form = new FormData();
    form.append("location_id", location_id);
    form.append("product_id", product_id);
    form.append("state", sample_state);
    const res = await fetch(`/api/sample/${sample_id}`, {
      method: "POST",
      body: form,
    });
    if (!res.ok) throw new Error(await res.text());
    // Optionally, show a success message or redirect
    await fetchSample();
  }

  onMount(fetchSample);
</script>

<Card.Header>
  <Card.Title>Edit Sample</Card.Title>
</Card.Header>
<Card.Content>
  {#if sample}
    <form onsubmit={saveSample} class="space-y-4">
      <div>
        <label class="block font-bold" for="sample-id">Sample ID</label>
        <div id="sample-id">{sample.DisplayId}</div>
      </div>
      <div>
        <ProductSelect
          bind:bindValue={product_id}
          placeholder="Select a product variant"
          id="variant-select"
        />
      </div>
      <div>
        <LocationSelect bind:bindValue={location_id} />
      </div>
      <div>
        <StateSelect bind:bindValue={sample_state} />
      </div>
      <div>
        <Button type="submit">Save</Button>
      </div>
    </form>
    <div class="mt-6">
      <h3 class="font-bold mb-2">Mods</h3>
      <div class="flex flex-wrap gap-2 mb-2 items-center">
        <label class="flex items-center gap-1 ml-4">
          <input type="checkbox" bind:checked={showFittedOnly} />
          <span class="text-sm">Show fitted only</span>
        </label>
      </div>
      {#if add_mod_error}
        <div class="text-red-500 text-sm mb-2">{add_mod_error}</div>
      {/if}
      <Table.Root>
        <Table.Header>
          <Table.Row>
            <Table.Head>Name</Table.Head>
            <Table.Head>Date Added</Table.Head>
            <Table.Head>Date Removed</Table.Head>
            <Table.Head>Status</Table.Head>
            <Table.Head>Action</Table.Head>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {#if sample.mods.length === 0}
            <Table.Row>
              <Table.Cell colspan={5} class="text-muted-foreground"
                >No mods.</Table.Cell
              >
            </Table.Row>
          {:else}
            {#each [...sample.mods]
              .filter((mod) => !showFittedOnly || !mod.time_removed)
              .sort((a, b) => {
                // Sort by time_added descending, then by time_removed descending if both removed
                const aTime = a.time_removed ? a.time_removed.getTime() : (a.time_added?.getTime() ?? 0);
                const bTime = b.time_removed ? b.time_removed.getTime() : (b.time_added?.getTime() ?? 0);
                return bTime - aTime;
              }) as mod (mod.id)}
              <Table.Row
                class={!mod.time_removed
                  ? "font-bold bg-green-50 dark:bg-green-900/20"
                  : ""}
              >
                <Table.Cell>{mod.name}</Table.Cell>
                <Table.Cell
                  >{mod.time_added?.toLocaleString?.() ?? ""}</Table.Cell
                >
                <Table.Cell
                  >{mod.time_removed
                    ? mod.time_removed.toLocaleString()
                    : "-"}</Table.Cell
                >
                <Table.Cell>
                  {#if !mod.time_removed}
                    <span class="text-green-700 dark:text-green-400"
                      >Fitted</span
                    >
                  {:else}
                    <span class="text-gray-500">Removed</span>
                  {/if}
                </Table.Cell>
                <Table.Cell>
                  {#if !mod.time_removed}
                    <Button
                      size="sm"
                      variant="destructive"
                      onclick={() => removeMod(mod.id)}>Remove</Button
                    >
                  {/if}
                </Table.Cell>
              </Table.Row>
            {/each}
          {/if}
        </Table.Body>
      </Table.Root>
      <div class="mt-4">
        <input
          type="text"
          class="border rounded px-2 py-1"
          placeholder="Add mod..."
          bind:value={new_mod_name}
          onkeydown={(e) => {
            if (e.key === "Enter") {
              e.preventDefault();
              addMod();
            }
          }}
          disabled={adding_mod}
        />
        <Button onclick={addMod} disabled={adding_mod || !new_mod_name.trim()}
          >Add</Button
        >
      </div>
    </div>
  {/if}
</Card.Content>
